package tools

import "os/exec"

type PriorityCmd struct {
	*exec.Cmd
	LowPriority bool
}

func (c *PriorityCmd) Start() error {
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
