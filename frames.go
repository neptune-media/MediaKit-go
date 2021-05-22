package mediakit

import (
	"fmt"
	"time"

	"github.com/neptune-media/MediaKit-go/ffprobe"
)

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
