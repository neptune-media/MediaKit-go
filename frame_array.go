package mediakit

import (
	"fmt"
	"io"
	"time"
)

// FrameArray is a list of timecodes representing frame positions in a video file
type FrameArray []time.Duration

func ReadFrameArray(r io.Reader) (FrameArray, error) {
	frames := make(FrameArray, 0)
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
