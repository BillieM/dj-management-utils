package helpers

import (
	"fmt"
	"os/exec"
)

// CmdExec Execute a command
func CmdExec(args ...string) error {

	baseCmd := args[0]
	cmdArgs := args[1:]

	cmd := exec.Command(baseCmd, cmdArgs...)
	_, err := cmd.Output()
	if err != nil {
		fmt.Println(string(err.(*exec.ExitError).Stderr))
		return err
	}

	return nil
}
