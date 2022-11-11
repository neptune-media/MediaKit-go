package tools

import (
	"context"
	"os/exec"
)

func ExecTool(ctx context.Context, t Tool) error {
	cmd := exec.CommandContext(ctx, t.GetCommand(), t.GetCommandArgs()...)

	// Set the buffers used to hold stdout and stderr
	cmd.Stdout, cmd.Stderr = t.GetOutputBuffers()

	// Start the process
	if err := cmd.Start(); err != nil {
		return err
	}

	if t.IsLowPriority() {
		if err := ReduceProcessPriority(cmd.Process); err != nil {
			return err
		}
	}

	// Wait for the process to exit
	return cmd.Wait()
}
