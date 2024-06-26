package cli

import (
	"fmt"
	"os"

	"github.com/billiem/seren-management/pkg/collection"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/streaming"
	"github.com/urfave/cli/v2"
)

func flattenDir(c *cli.Context) error {

	e, err := buildCliEnv(c.String("config"))

	if err != nil {
		return err
	}

	absPath, err := helpers.GetAbsOrWdPath(c.Args().First())

	if err != nil {
		return err
	}

	opEnv := e.opEnv()

	opEnv.FlattenDirectory(absPath)

	return nil
}

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
	opEnv.BuildOperationHandler(func(f float64) {
		fmt.Println(f)
	}, func(_ map[string]any) {

	}, func(err error) {
		fmt.Println(err)
	})

	opEnv.ReadCollection(c.Context, traktorCollectionOpts)

	return nil
}

func getSoundcloudPlaylist(c *cli.Context) error {

	// e, err := buildCliEnv(c.String("config"))

	// if err != nil {
	// 	return err
	// }

	// soundcloudOpts := operations.GetSoundCloudPlaylistOpts{
	// 	PlaylistURL: c.String("url"),
	// }

	// opEnv := e.opEnv()
	// opEnv.RegisterStepHandler(stepHandler{})

	// opEnv.GetSoundCloudPlaylist(c.Context, soundcloudOpts, func(p streaming.SoundCloudPlaylist) {})

	return nil
}

func getSpotifyPlaylist(c *cli.Context) error {

	e, err := buildCliEnv(c.String("config"))

	if err != nil {
		return err
	}

	_ = e

	return nil
}

func generateClientID(c *cli.Context) error {

	clientID, err := streaming.GenerateSoundCloudClientID()

	if err != nil {
		return err
	}

	fmt.Println(clientID)

	return nil
}
