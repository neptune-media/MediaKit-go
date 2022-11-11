package tools

import (
	"context"
	"io"
)

type Check interface {
	Validate() error
}

type Tool interface {
	Do() error
	DoWithContext(context.Context) error

	GetCommand() string
	GetCommandArgs() []string
	GetStdout() []byte
	GetStderr() []byte

	// GetOutputBuffers returns buffers for writing stdout and stderr to.
	GetOutputBuffers() (io.Writer, io.Writer)

	// IsLowPriority returns if the process should be run at low priority
	IsLowPriority() bool
}
