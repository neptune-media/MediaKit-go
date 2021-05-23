package mkvmerge

import (
	"os/exec"
	"strings"
)

type Runner struct {
	Args           []string
	InputFilename  string
	OutputFilename string

	output []byte
}

func (r *Runner) Do() error {
	args := r.buildArgs()
	c := exec.Command("mkvmerge", args...)
	o, err := c.CombinedOutput()

	r.output = make([]byte, len(o))
	copy(r.output, o)
	return err
}

func (r *Runner) GetCommandString() string {
	cmd := []string{"mkvmerge"}
	cmd = append(cmd, r.buildArgs()...)
	return strings.Join(cmd, " ")
}

func (r *Runner) GetOutput() []byte {
	o := make([]byte, len(r.output))
	copy(o, r.output)
	return o
}

func (r *Runner) buildArgs() []string {
	var args []string
	if r.OutputFilename != "" {
		args = append(args, "-o", r.OutputFilename)
	}

	if len(r.Args) > 0 {
		args = append(args, r.Args...)
	}

	if r.InputFilename != "" {
		args = append(args, r.InputFilename)
	}

	return args
}
