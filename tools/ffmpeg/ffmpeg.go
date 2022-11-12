package ffmpeg

import (
	"bytes"
	"context"
	"fmt"
	"github.com/neptune-media/MediaKit-go/tools"
	"io"
)

func filter(ss []string, test func(string) bool) []string {
	r := make([]string, 0)
	for _, s := range ss {
		if test(s) {
			r = append(r, s)
		}
	}
	return r
}

// FFmpeg is a wrapper for the ffmpeg binary
type FFmpeg struct {
	// A list of ISO 639-1 language codes to select for output
	AudioLanguages []string

	// Audio output encoding options
	AudioOptions EncodingOptions

	// Discard audio from output
	DiscardAudio bool

	// Discard subtitles from output
	DiscardSubtitles bool

	// Discard video from output
	DiscardVideo bool

	// Modifies output to enable "Fast Start" for web streaming
	EnableFastStart bool

	// A list of raw args to set for input
	InputArgs []string

	// Path to read input from
	InputFilename string

	// A list of raw args to set for output
	OutputArgs []string

	// Path to store output in
	OutputFilename string

	// A list of ISO 639-1 language codes to select for output
	SubtitleLanguages []string

	// Subtitle output encoding options
	SubtitleOptions EncodingOptions

	// Uses a lower process priority for ffmpeg
	UseLowerPriority bool

	// Video output encoding options
	VideoOptions EncodingOptions

	stdout       []byte
	stderr       []byte
	stdoutBuffer bytes.Buffer
	stderrBuffer bytes.Buffer
}

func (f *FFmpeg) Do() error {
	return f.DoWithContext(context.Background())
}

func (f *FFmpeg) DoWithContext(ctx context.Context) error {
	// Reset buffers
	f.stdout = make([]byte, 0)
	f.stderr = make([]byte, 0)
	f.stdoutBuffer.Reset()
	f.stderrBuffer.Reset()

	// Execute
	err := tools.ExecTool(ctx, f)

	// Copy output to buffer for later
	f.stdout = make([]byte, f.stdoutBuffer.Len())
	f.stderr = make([]byte, f.stderrBuffer.Len())
	copy(f.stdout, f.stdoutBuffer.Bytes())
	copy(f.stderr, f.stderrBuffer.Bytes())

	return err
}

func (f *FFmpeg) GetCommand() string {
	return "ffmpeg"
}

func (f *FFmpeg) GetCommandArgs() []string {
	args := make([]string, 0)
	args = append(args, f.InputArgs...)
	args = append(args,
		"-i",
		f.InputFilename,
	)

	if f.DiscardAudio {
		args = append(args, "-an")
	} else if f.AudioOptions != nil {
		args = append(args, "-c:a")
		args = append(args, f.AudioOptions.GetCodecOptions()...)
	}

	if f.DiscardSubtitles {
		args = append(args, "-sn")
	} else if f.SubtitleOptions != nil {
		args = append(args, "-c:s")
		args = append(args, f.SubtitleOptions.GetCodecOptions()...)
	}

	if f.DiscardVideo {
		args = append(args, "-vn")
	} else if f.VideoOptions != nil {
		args = append(args, "-c:v")
		args = append(args, f.VideoOptions.GetCodecOptions()...)
	}

	if f.EnableFastStart {
		args = append(args, "-movflags", "+faststart")
	}

	args = append(args, "-map", "0:v:0")
	if len(f.AudioLanguages) > 0 {
		for _, lang := range f.AudioLanguages {
			args = append(args, "-map", fmt.Sprintf("0:a:m:language:%s", lang))
		}
	}

	if len(f.SubtitleLanguages) > 0 {
		for _, lang := range f.SubtitleLanguages {
			args = append(args, "-map", fmt.Sprintf("0:s:m:language:%s", lang))
		}
	}

	args = append(args, f.OutputArgs...)
	args = append(args,
		"-y",
		f.OutputFilename,
	)

	return filter(args, func(s string) bool { return len(s) > 0 })
}

func (f *FFmpeg) GetStdout() []byte {
	return f.stdout
}

func (f *FFmpeg) GetStderr() []byte {
	return f.stderr
}

func (f *FFmpeg) GetOutputBuffers() (io.Writer, io.Writer) {
	return &f.stdoutBuffer, &f.stderrBuffer
}

func (f *FFmpeg) IsLowPriority() bool {
	return f.UseLowerPriority
}
