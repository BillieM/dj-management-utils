package operations

import (
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
	operations "github.com/billiem/seren-management/pkg/operations/mp3"
)

type OpEnv struct {
	helpers.Config
	Logger helpers.SerenLogger
	operationHandler
	*data.SerenDB

	FinishOperation
	ProgressOperation

	Mp3EnvBuilder
}

/*
Mp3Env defines the interface for the mp3 conversion operations
*/
type Mp3Env interface {
	GetConvertPaths() ([]string, error)
	GetConvertTracks([]string) ([]operations.ConvertTrack, int, []error)
	ConvertTracks([]string) (int, int, []error)
}

/*
Mp3EnvBuilder defines the interface for building the mp3 conversion environment

This allows us to mock out a custom environment for testing whilst
minimising the amount of boilerplate required to call the operations
package from the UI packages
*/
type Mp3EnvBuilder interface {
	mp3Env() *operations.Mp3Env
}

func (e *OpEnv) Mp3Env() *operations.Mp3Env {
	return &operations.Mp3Env{
		Config: e.Config,
		Logger: e.Logger,
	}
}

type FinishOperation interface {
	FinishError(error)
	FinishSuccess(map[string]any)
}

type ProgressOperation interface {
	Progress(float64)
}
