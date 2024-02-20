package ffmpeg

import "fmt"

// GenericAudioOptions represents a generic audio codec
type GenericAudioOptions struct {
	// -b:a flag
	Bitrate int

	Codec string
}

func (o *GenericAudioOptions) GetOptions() []string {
	args := []string{
		o.Codec,
	}
	args = append(args, o.getBitrate()...)

	return filter(args, func(s string) bool { return len(s) > 0 })
}

func (o *GenericAudioOptions) getBitrate() []string {
	if o.Bitrate == 0 {
		return []string{}
	}

	return []string{"-b:a", fmt.Sprintf("%dk", o.Bitrate)}
}
