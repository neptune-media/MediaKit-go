package mediakit

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"time"

	"github.com/neptune-media/MediaKit-go/ffprobe"
)

type FrameSeeker struct {
	Frames []time.Duration
	pos    int
}

type FrameArray []time.Duration

func ReadVideoIFrames(filename string) ([]time.Duration, error) {
	frames := make([]time.Duration, 0)

	c := exec.Command("ffprobe",
		"-select_streams",
		"v",
		"-show_frames",
		"-of",
		"json=compact=1",
		filename,
	)

	stdout, err := c.StdoutPipe()
	if err != nil {
		return nil, err
	}

	fmt.Printf("Analyzing file, this may take a few minutes...\n")
	if err := c.Start(); err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	go func() {
		for {
			if _, err := io.Copy(&buf, stdout); err != nil {
				return
			}
		}
	}()

	if err := c.Wait(); err != nil {
		return nil, err
	}

	parsed, err := ffprobe.ReadFFProbeOutput(&buf)
	if err != nil {
		return nil, err
	}

	for _, frame := range parsed.Frames {
		if frame.PictType == "I" {
			frames = append(frames, time.Duration(frame.PktPTS)*time.Millisecond)
		}
	}

	return frames, nil
}

func ReadFrameArray(r io.Reader) ([]time.Duration, error) {
	frames := make([]time.Duration, 0)
	for {
		var t time.Duration
		if _, err := fmt.Fscanf(r, "%d\n", &t); err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		frames = append(frames, t)
	}

	return frames, nil
}

func (a FrameArray) Write(w io.Writer) error {
	for _, frame := range []time.Duration(a) {
		if _, err := fmt.Fprintf(w, "%d\n", frame); err != nil {
			return err
		}
	}

	return nil
}

func (f *FrameSeeker) Current() time.Duration {
	return f.Frames[f.pos]
}

func (f *FrameSeeker) AtEnd() bool {
	return !(f.pos < len(f.Frames))
}

func (f *FrameSeeker) Next() {
	if !f.AtEnd() {
		f.pos += 1
	}
}

func (f *FrameSeeker) Peek() time.Duration {
	if f.pos+1 < len(f.Frames) {
		return f.Frames[f.pos+1]
	}

	return f.Frames[f.pos]
}

func (f *FrameSeeker) Reset() {
	f.pos = 0
}
