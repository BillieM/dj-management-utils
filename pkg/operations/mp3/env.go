package operations

import "github.com/billiem/seren-management/pkg/helpers"

/*
Mp3Env is the environment used for all mp3 conversion operations, it
is a simple wrapper around the OpEnv struct
*/
type Mp3Env struct {
	helpers.Config
	Logger helpers.SerenLogger
}
