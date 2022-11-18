//go:build !windows

package tools

import (
	"golang.org/x/sys/unix"
	"os"
)

// ReduceProcessPriority is a multi-os helper for reducing
// the run priority of a process.
func ReduceProcessPriority(p *os.Process) error {
	// Priority 19 is the lowest priority
	return unix.Setpriority(unix.PRIO_PROCESS, p.Pid, 19)
}
