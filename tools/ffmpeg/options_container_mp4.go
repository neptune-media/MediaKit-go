package ffmpeg

// Mp4ContainerOptions represents options for outputting to a mp4 container
type Mp4ContainerOptions struct {
	// Modifies output to enable "Fast Start" for web streaming
	EnableFastStart bool
}

func (o *Mp4ContainerOptions) GetOptions() []string {
	args := []string{"-f", "mp4"}
	args = append(args, o.getMovFlags()...)
	return args
}

func (o *Mp4ContainerOptions) getMovFlags() []string {
	flags := make([]string, 0)

	if o.EnableFastStart {
		flags = append(flags, "+faststart")
	}

	if len(flags) == 0 {
		return flags
	}

	return append([]string{"-movflags"}, flags...)
}
