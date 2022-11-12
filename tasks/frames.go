package tasks

import (
	"context"
	"fmt"
	"github.com/neptune-media/MediaKit-go/tools/ffprobe"
	"time"
)

// ReadVideoIFramesFromFile reads the specified file and returns a list of I-frames
func ReadVideoIFramesFromFile(filename string) ([]time.Duration, error) {
	fmt.Printf("Analyzing file, this may take a few minutes...\n")
	runner := &ffprobe.FFProbe{Filename: filename, GetFrames: true, LowPriority: true}
	return ReadVideoIFrames(runner)
}

// ReadVideoIFrames executes the provided FFProbe to return a list of I-frames
func ReadVideoIFrames(runner *ffprobe.FFProbe) ([]time.Duration, error) {
	frames := make([]time.Duration, 0)
	runner.GetFrames = true

	ctx, cancel := context.WithCancel(context.TODO())
	defer cancel()

	err := runner.DoWithContext(ctx)
	if err != nil {
		fmt.Printf("===== begin stdout =====\n%s\n===== end stdout =====\n", runner.GetStdout())
		fmt.Printf("===== begin stderr =====\n%s\n===== end stderr =====\n", runner.GetStderr())
		return nil, err
	}

	parsed, err := runner.GetOutput()
	if err != nil {
		return nil, err
	}

	for _, frame := range parsed.Frames {
		if frame.PictType == "I" {
			frames = append(frames, time.Duration(frame.PTS)*time.Millisecond)
		}
	}

	return frames, nil
}
