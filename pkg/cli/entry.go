package cli

import (
	"os"
	"time"

	"github.com/urfave/cli/v2"
)

func Entry() {
	cmd := &cli.App{
		Version:  "v0.0.1",
		Compiled: time.Now(),
		Authors: []*cli.Author{
			{
				Name:  "Billie M",
				Email: "billiemerz@gmail.com",
			},
		},
		Usage: "A collection of useful utilities to help manage your DJ library",
		Commands: []*cli.Command{
			{
				Name:    "convertmp3",
				Aliases: []string{"cmp3"},
				Usage:   "Converts a single mp3 file to wav",
				Action:  convertMp3,
				Flags:   []cli.Flag{},
			},
		},
	}

	cmd.Run(os.Args)
}
