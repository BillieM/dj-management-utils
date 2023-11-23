package cli

import (
	"errors"
	"os"
	"time"

	"github.com/billiem/seren-management/pkg/helpers"
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
			{
				Name:    "read-collection",
				Aliases: []string{"rc"},
				Usage:   "Read a platforms collection into the applications database",
				Subcommands: []*cli.Command{
					{
						Name:    "traktor",
						Aliases: []string{"t"},
						Usage:   "Reads a Traktor collection into the applications database",
						Action:  readTraktorCollection,
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "in",
								Aliases:  []string{"i"},
								Usage:    "Path to the Traktor collection file, if not given we default to the path stored in application config",
								Required: false,
							},
							&cli.StringFlag{
								Name:     "out",
								Aliases:  []string{"o"},
								Usage:    "Path to store the new traktor collection file, if false we default to {in}_new.nml",
								Required: false,
							},
						},
					},
				},
			},
		},
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     "config",
				Aliases:  []string{"c"},
				Usage:    "Optional path to the config file",
				Value:    "config.json",
				Required: false,
				Action: func(c *cli.Context, s string) error {
					if _, err := os.Stat(s); errors.Is(err, os.ErrNotExist) {
						return helpers.ErrConfigDoesNotExist
					}
					return nil
				},
			},
		},
	}

	cmd.Run(os.Args)
}
