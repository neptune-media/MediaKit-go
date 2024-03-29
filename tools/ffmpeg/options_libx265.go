package ffmpeg

import (
	"fmt"
)

// Libx265Options represent the libx265 codec
type Libx265Options struct {
	// Corresponds to the -crf flag
	CRF int

	// -preset flag
	Preset string

	// -tune flag
	Tune string

	// Add -crf flag
	UseCRF bool
}

func (o *Libx265Options) GetOptions() []string {
	args := []string{
		"libx265",
	}
	args = append(args, o.getConstantRateFactor()...)
	args = append(args, o.getPreset()...)
	args = append(args, o.getTune()...)

	return filter(args, func(s string) bool { return len(s) > 0 })
}

func (o *Libx265Options) getConstantRateFactor() []string {
	if !o.UseCRF {
		return []string{}
	}

	return []string{"-crf", fmt.Sprintf("%d", o.CRF)}
}

func (o *Libx265Options) getPreset() []string {
	if o.Preset == "" {
		return []string{}
	}

	return []string{"-preset", o.Preset}
}

func (o *Libx265Options) getTune() []string {
	if o.Tune == "" {
		return []string{}
	}

	return []string{"-tune", o.Tune}
}
