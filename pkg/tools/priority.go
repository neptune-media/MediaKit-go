package tools

import (
	"os/exec"
	"strings"

	"github.com/go-logr/logr"
)

type PriorityCmd struct {
	*exec.Cmd
	Logger      logr.Logger
	LowPriority bool
}

func (c *PriorityCmd) Start() error {
	c.Logger.V(1).Info("executing command", "command", strings.Join(c.Args, " "))
	err := c.Cmd.Start()
	if err != nil {
		return err
	}

	if c.LowPriority {
		err = ReduceProcessPriority(c.Cmd.Process)
		if err != nil {
			return err
		}
	}

	return nil
}
