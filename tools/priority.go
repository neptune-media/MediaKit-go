//go:build !windows

package tools

import "os"

func ReduceProcessPriority(p *os.Process) error {
	// TODO: Fill this out when we have a linux machine to test on
	// For now, just be a no-op
	return nil
}
