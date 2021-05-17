package mediakit

import (
	"bytes"
	"fmt"
	"github.com/neptune-media/MediaKit-go/ffprobe"
	"io"
	"os/exec"
)

type FrameNumber int

func ReadVideoIFrames(filename string) ([]FrameNumber, error) {
	frames := make([]FrameNumber, 0)

	c := exec.Command("ffprobe",
		"-select_streams",
		"v",
		"-show_frames",
		"-of",
		"json",
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
			frames = append(frames, FrameNumber(frame.PktPTS))
		}
	}

	return frames, nil
}
