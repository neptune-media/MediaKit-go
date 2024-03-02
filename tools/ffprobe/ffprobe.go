package ffprobe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/neptune-media/MediaKit-go/tools"
	"io"
)

// FFProbe is a wrapper for the ffprobe binary
type FFProbe struct {
	// Filename is the path of the file to analyze
	Filename string

	// GetFrameCount - when true, counts the number of frames in Filename
	GetFrameCount bool
	// GetFrame - when true, returns details of every frame in Filename
	GetFrames bool
	// LowPriority - when true, runs the ffprobe at a low process priority
	LowPriority bool
	// Threads - determines the number of CPU threads to use (0 for auto)
	Threads int
	// UseThreads - when true, sets the -threads flag with Threads value
	UseThreads bool

	stdout       []byte
	stderr       []byte
	stdoutBuffer bytes.Buffer
	stderrBuffer bytes.Buffer
}

func (f *FFProbe) Do() error {
	return f.DoWithContext(context.Background())
}

func (f *FFProbe) DoWithContext(ctx context.Context) error {
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

func (f *FFProbe) GetCommand() string {
	return "ffprobe"
}

func (f *FFProbe) GetCommandArgs() []string {
	args := []string{
		"-of",
		"json=compact=1",
	}

	if f.GetFrameCount {
		args = append(args, "-show_streams", "-count_frames")
	}

	if f.GetFrames {
		args = append(args, "-show_frames")
	}

	if f.UseThreads {
		args = append(args, "-threads", fmt.Sprintf("%d", f.Threads))
	}

	args = append(args, f.Filename)
	return args
}

func (f *FFProbe) GetOutput() (*Output, error) {
	output := &Output{}
	err := json.Unmarshal(f.stdout, output)
	return output, err
}

func (f *FFProbe) GetStdout() []byte {
	return f.stdout
}

func (f *FFProbe) GetStderr() []byte {
	return f.stderr
}

func (f *FFProbe) GetOutputBuffers() (io.Writer, io.Writer) {
	return &f.stdoutBuffer, &f.stderrBuffer
}

func (f *FFProbe) IsLowPriority() bool {
	return f.LowPriority
}
