package cli

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
)

func convertMp3(c *cli.Context) error {
	workingDir, err := os.Getwd()

	if err != nil {
		return err
	}

	fmt.Println(workingDir)

	return nil
}
