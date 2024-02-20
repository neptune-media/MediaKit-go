package ffmpeg

// EncodingOptions provide a method for specifying various video and audio
// codec options, and formatting them to ffmpeg commandline arguments
type EncodingOptions interface {
	GetOptions() []string
}
