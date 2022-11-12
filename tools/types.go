package tools

import (
	"context"
	"io"
)

// Check is used to validate if conditions are met for executing
// a particular tool.
type Check interface {
	Validate() error
}

// Tool is a wrapper for an executable program, responsible for generating
// all the command-line arguments needed to execute successfully.
type Tool interface {
	// Do should execute the program
	Do() error

	// DoWithContext should execute the program under the provided context
	DoWithContext(context.Context) error

	// GetCommand should return the name of the binary to execute (ffmpeg, mkvmerge, etc)
	GetCommand() string

	// GetCommandArgs should return a list of arguments to pass to the command
	GetCommandArgs() []string

	// GetStdout should return the stdout output after the program is executed
	GetStdout() []byte

	// GetStderr should return the stderr output after the program is executed
	GetStderr() []byte

	// GetOutputBuffers returns buffers for writing stdout and stderr to.
	GetOutputBuffers() (io.Writer, io.Writer)

	// IsLowPriority returns if the process should be run at low priority
	IsLowPriority() bool
}
