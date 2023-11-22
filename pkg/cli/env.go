package cli

import (
	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

type cliEnv struct {
	helpers.Config
	*database.SerenDB
}

func (e cliEnv) opEnv() operations.OpEnv {
	return operations.OpEnv{
		Config:  e.Config,
		SerenDB: e.SerenDB,
	}
}

func buildCliEnv(configPath string) (*cliEnv, error) {

	cfg, err := helpers.LoadCLIConfig(configPath)

	if err != nil {
		return nil, err
	}

	db, err := database.Connect()

	if err != nil {
		return nil, err
	}

	e := &cliEnv{cfg, db}

	return e, nil
}
