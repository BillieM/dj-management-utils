package cli

import (
	"fmt"
	"os"

	"github.com/billiem/seren-management/pkg/collection"
	"github.com/billiem/seren-management/pkg/helpers"
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

	collectionInPath, err := helpers.GetAbsOrWdPath(c.String("in"))
	if err != nil {
		return err
	}

	collectionOutPath, err := helpers.GetAbsOrWdPath(c.String("out"))
	if err != nil {
		return err
	}

	traktorCollectionOpts := collection.ReadTraktorOpts{
		CollectionInPath:  collectionInPath,
		CollectionOutPath: collectionOutPath,
	}

	opEnv := e.opEnv()
	opEnv.RegisterStepHandler(stepHandler{})

	opEnv.ReadCollection(c.Context, traktorCollectionOpts)

	return nil
}

func getSoundcloudPlaylists(c *cli.Context) error {

	e, err := buildCliEnv(c.String("config"))

	if err != nil {
		return err
	}

	_ = e

	return nil
}

func getSpotifyPlaylists(c *cli.Context) error {

	e, err := buildCliEnv(c.String("config"))

	if err != nil {
		return err
	}

	_ = e

	return nil
}
