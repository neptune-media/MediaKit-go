package ffmpeg

// CopyOptions represents the copy codec
type CopyOptions struct{}

func (o *CopyOptions) GetCodecOptions() []string {
	return []string{"copy"}
}
