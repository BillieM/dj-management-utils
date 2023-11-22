package cli

import (
	"fmt"
	"os"

	"github.com/billiem/seren-management/pkg/collection"
	"github.com/billiem/seren-management/pkg/operations"
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

func readTraktorCollection(c *cli.Context) error {

	e, err := buildCliEnv(c.String("config"))

	if err != nil {
		return err
	}

	traktorCollectionOpts := collection.ReadTraktorOpts{
		CollectionPath: c.String("path"),
	}

	operations.ReadCollection(c.Context, e.opEnv(), operationProcess{}, traktorCollectionOpts)

	return nil
}
