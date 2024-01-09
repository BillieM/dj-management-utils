package cli

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

type cliEnv struct {
	helpers.Config
	*data.SerenDB
	logger helpers.SerenLogger
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
		return nil, fault.Wrap(
			err,
			fmsg.With("Error loading CLI config"),
		)
	}

	loggers, err := helpers.BuildAppLoggers()

	if err != nil {
		return nil, fault.Wrap(
			err,
			fmsg.With("Error building loggers"),
		)
	}

	sDB, err := data.Connect(cfg, loggers.DBLogger)

	if err != nil {
		return nil, fault.Wrap(
			err,
			fmsg.With("Error connecting to database"),
		)
	}

	e := &cliEnv{cfg, sDB, loggers.AppLogger}

	return e, nil
}
