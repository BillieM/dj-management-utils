package operations

import (
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations/internal"
)

/*
StemEnv is the environment used for all stem conversion operations, it
is a subset of the OpEnv struct
*/
type StemEnv struct {
	*internal.OperationHandler
	helpers.Config
	Logger helpers.SerenLogger
}
