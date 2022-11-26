package shellout

import (
	"bytes"
	"errors"
	"os/exec"
)

// A Result represents the result of running a shell command.
type Result struct {
	Stdout   bytes.Buffer
	Stderr   bytes.Buffer
	ExitCode int
}

// Run executes a shell command and returns the Result, or an error indicating why the command
// could not be completed.
func Run(name string, arg ...string) (*Result, error) {
	result := &Result{}
	cmd := exec.Command("git", arg...)
	cmd.Stdout = &result.Stdout
	cmd.Stderr = &result.Stderr
	err := cmd.Run()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		result.ExitCode = exitErr.ExitCode()
	} else if err != nil {
		return nil, err
	}
	return result, nil
}
