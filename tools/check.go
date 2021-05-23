package tools

import (
	"fmt"
	"os/exec"
)

type Checks struct {
	checks []Check
}

type ExecutableCheck struct {
	Executable string
}

func NewChecks(checks ...Check) *Checks {
	return &Checks{
		checks: checks,
	}
}

func (c *Checks) AddCheck(check Check) {
	c.checks = append(c.checks, check)
}

func (c *Checks) Run() error {
	for _, check := range c.checks {
		if err := check.Validate(); err != nil {
			return err
		}
	}

	return nil
}

func (e *ExecutableCheck) Validate() error {
	p, err := exec.LookPath(e.Executable)
	if err == nil {
		fmt.Printf("Found executable %s: %s\n", e.Executable, p)
	}

	return err
}
