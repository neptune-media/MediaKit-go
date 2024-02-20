package ffmpeg

import "fmt"

// MkvContainerOptions represents options for outputting to a mp4 container
type MkvContainerOptions struct {
	// Defines amount of space (in kilobytes) at beginning of file to reserve for writing video cues.
	// Recommended values are 50 per hour of video, see:
	// https://www.ffmpeg.org/ffmpeg-formats.html#matroska
	ReserveIndexSpace int
}

func (o *MkvContainerOptions) GetOptions() []string {
	args := []string{"-f", "matroska"}
	args = append(args, o.getMuxerOptions()...)
	return args
}

func (o *MkvContainerOptions) getMuxerOptions() []string {
	flags := make([]string, 0)

	if o.ReserveIndexSpace > 0 {
		flags = append(flags, "-reserve_index_space", fmt.Sprintf("%dk", o.ReserveIndexSpace))
	}

	return flags
}
