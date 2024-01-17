package operations

import (
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
)

type OpEnv struct {
	helpers.Config
	Logger helpers.SerenLogger
	operationHandler
	*data.SerenDB
}
