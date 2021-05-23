package ffprobe

import (
	"encoding/json"
	"io"
)

func ReadFFProbeOutput(r io.Reader) (*Output, error) {
	o := &Output{
		Frames: make([]FrameInfo, 0),
	}

	if err := json.NewDecoder(r).Decode(o); err != nil {
		return nil, err
	}

	return o, nil
}
