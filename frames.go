package mediakit

import (
	"fmt"
	"io"
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

	fmt.Printf("Analyzing file, this may take a few minutes...\n")
	runner := &ffprobe.Runner{Filename: filename}

	parsed, err := runner.GetFrames()
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

func (a FrameArray) WriteTo(w io.Writer) (int64, error) {
	var n int64 = 0
	for _, frame := range []time.Duration(a) {
		count, err := fmt.Fprintf(w, "%d\n", frame)
		n += int64(count)
		if err != nil {
			return n, err
		}
	}

	return n, nil
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
