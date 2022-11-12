package mkvmerge

import (
	"bytes"
	"context"
	"github.com/neptune-media/MediaKit-go/tools"
	"io"
)

// MKVMerge is a wrapper for the mkvmerge binary
type MKVMerge struct {
	// Args is a list of raw arguments to pass to mkvmerge
	Args []string
	// InputFilename is the path of the file to read
	InputFilename string
	// LowPriority - when true, runs the mkvmerge at a low process priority
	LowPriority bool
	// OutputFilename is the path of the file to write
	OutputFilename string

	stdout       []byte
	stderr       []byte
	stdoutBuffer bytes.Buffer
	stderrBuffer bytes.Buffer
}

func (m *MKVMerge) Do() error {
	return m.DoWithContext(context.Background())
}

func (m *MKVMerge) DoWithContext(ctx context.Context) error {
	// Reset buffers
	m.stdout = make([]byte, 0)
	m.stderr = make([]byte, 0)
	m.stdoutBuffer.Reset()
	m.stderrBuffer.Reset()

	// Execute
	err := tools.ExecTool(ctx, m)

	// Copy output to buffer for later
	m.stdout = make([]byte, m.stdoutBuffer.Len())
	m.stderr = make([]byte, m.stderrBuffer.Len())
	copy(m.stdout, m.stdoutBuffer.Bytes())
	copy(m.stderr, m.stderrBuffer.Bytes())

	return err
}

func (m *MKVMerge) GetCommand() string {
	return "mkvmerge"
}

func (m *MKVMerge) GetCommandArgs() []string {
	args := make([]string, 0)

	if m.OutputFilename != "" {
		args = append(args, "-o", m.OutputFilename)
	}

	if len(m.Args) > 0 {
		args = append(args, m.Args...)
	}

	if m.InputFilename != "" {
		args = append(args, m.InputFilename)
	}

	return args
}

func (m *MKVMerge) GetStdout() []byte {
	return m.stdout
}

func (m *MKVMerge) GetStderr() []byte {
	return m.stderr
}

func (m *MKVMerge) GetOutputBuffers() (io.Writer, io.Writer) {
	return &m.stdoutBuffer, &m.stderrBuffer
}

func (m *MKVMerge) IsLowPriority() bool {
	return m.LowPriority
}
