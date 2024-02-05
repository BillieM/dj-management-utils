package operations

import (
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations/internal"
)

/*
Mp3Env is the environment used for all mp3 conversion operations, it
is a subset of the OpEnv struct
*/
type Mp3Env struct {
	*internal.OperationHandler
	Config helpers.Config
	Logger helpers.SerenLogger
}
