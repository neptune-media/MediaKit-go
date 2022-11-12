package tools

import (
	"fmt"
	"os/exec"
)

// Checks is a helper to run multiple Check in one go
type Checks struct {
	checks []Check
}

// ExecutableCheck is used to determine if a specific executable
// can be found in the PATH env variable
type ExecutableCheck struct {
	// Executable is the name of the executable to find
	Executable string
}

func NewChecks(checks ...Check) *Checks {
	return &Checks{
		checks: checks,
	}
}

// AddCheck adds a Check to be validated later
func (c *Checks) AddCheck(check Check) {
	c.checks = append(c.checks, check)
}

// Run will run all added Check, and stop on the first error
func (c *Checks) Run() error {
	for _, check := range c.checks {
		if err := check.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// Validate calls exec.LookPath to find the specified executable
func (e *ExecutableCheck) Validate() error {
	p, err := exec.LookPath(e.Executable)
	if err == nil {
		fmt.Printf("Found executable %s: %s\n", e.Executable, p)
	}

	return err
}
