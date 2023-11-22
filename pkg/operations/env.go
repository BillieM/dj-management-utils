package operations

import (
	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
)

type OpEnv struct {
	helpers.Config
	*database.SerenDB
}
