package helpers

import (
	"os/exec"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
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
			errDetail = fault.New(string(execExitError.Stderr))
		} else if execError, ok := err.(*exec.Error); ok {
			errDetail = execError.Err
		} else {
			errDetail = fault.New("unknown execution error")
		}
		return outStr, fault.Wrap(
			err,
			fmsg.With(errDetail.Error()),
		)
	}

	return outStr, nil
}
