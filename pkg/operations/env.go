package operations

import (
	"context"

	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations/internal"

	mp3 "github.com/billiem/seren-management/pkg/operations/mp3"
	stems "github.com/billiem/seren-management/pkg/operations/stems"
)

/*
OpEnv is the base environment used for all operations, it contains all the
dependencies required to perform operations.

Individual operations may generate their own environments that are subsets of
this struct, for example, the StemEnv and Mp3Env structs.
*/
type OpEnv struct {
	helpers.Config
	Logger helpers.SerenLogger
	internal.OperationHandler
	*data.SerenDB

	Mp3EnvBuilder  func() Mp3Env
	StemEnvBuilder func() StemEnv
}

func (e *OpEnv) BuildOperationHandler(
	progressFunc func(float64),
	completeFunc func(map[string]any),
	errorFunc func(error),
) {
	e.OperationHandler = *internal.BuildOperationHandler(
		progressFunc,
		completeFunc,
		errorFunc,
	)
}

/*
Mp3Env defines the interface for the mp3 conversion operations

This allows us to mock out a custom environment for testing, checking if
the underlying methods are being called correctly
*/
type Mp3Env interface {
	GetMp3Paths(string, bool) ([]string, error)
	GetMp3Tracks([]string, string) ([]mp3.ConvertTrack, int, []error)
	ConvertMp3Tracks(context.Context, []mp3.ConvertTrack)
}

/*
AttachDefaultMp3EnvBuilder attaches the default mp3 environment builder to the OpEnv

This method is called by
*/
func (e *OpEnv) AttachDefaultMp3EnvBuilder() {
	e.Mp3EnvBuilder = func() Mp3Env {
		return &mp3.Mp3Env{
			OperationHandler: &e.OperationHandler,
			Config:           e.Config,
			Logger:           e.Logger,
		}
	}
}

/*
StemEnv defines the interface for the stem conversion operations

This allows us to mock out a custom environment for testing, checking if
the underlying methods are being called correctly
*/
type StemEnv interface {
	GetStemPaths(string, bool) ([]string, error)
	GetStemTracks([]string, string, stems.StemSeparationType) ([]stems.StemTrack, int, []error)
	ConvertStemTracks(context.Context, []stems.StemTrack)
}

/*
AttachDefaultStemEnvBuilder attaches the default stem environment builder to the OpEnv

This allows us to provide a custom StemEnvBuilder func for testing
*/
func (e *OpEnv) AttachDefaultStemEnvBuilder() {
	e.StemEnvBuilder = func() StemEnv {
		return &stems.StemEnv{
			OperationHandler: &e.OperationHandler,
			Config:           e.Config,
			Logger:           e.Logger,
		}
	}
}
