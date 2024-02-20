package ffmpeg

// CopyOptions represents the copy codec
type CopyOptions struct{}

func (o *CopyOptions) GetOptions() []string {
	return []string{"copy"}
}
