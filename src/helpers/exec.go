package helpers

import (
	"errors"
	"fmt"
	"os/exec"
)

// CmdExec Execute a command
func CmdExec(args ...string) (string, error) {

	baseCmd := args[0]
	cmdArgs := args[1:]

	cmd := exec.Command(baseCmd, cmdArgs...)
	out, err := cmd.CombinedOutput()
	outStr := string(out)

	if err != nil {
		var errDetail error
		if execExitError, ok := err.(*exec.ExitError); ok {
			errDetail = errors.New(string(execExitError.Stderr))
		} else if execError, ok := err.(*exec.Error); ok {
			errDetail = execError.Err
		} else {
			errDetail = errors.New("unknown execution error")
		}
		return outStr, fmt.Errorf("%s: %w", err.Error(), errDetail)
	}

	return outStr, nil
}
