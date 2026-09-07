package tools

import (
	"context"
	"os/exec"
)

type Executable struct {
	// Name is the name of the executable to find
	Name string

	// Path is the full detected path to the executable
	Path string

	// VersionArg is the arguments to pass to the executable to get
	// version information.  If left empty, no check is made.
	VersionArg string

	// Version is the detected version of the executable
	Version string
}

func (e *Executable) Validate(ctx context.Context) error {
	path, err := exec.LookPath(e.Name)
	if err != nil {
		return err
	}

	e.Path = path

	if e.VersionArg != "" {
		cmd := exec.CommandContext(ctx, e.Name, e.VersionArg)
		out, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}

		e.Version = string(out)
	}

	return nil
}
