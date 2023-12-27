package helpers

import (
	"errors"
	"fmt"
	"strings"
)

/*
ErrorContains is a helper function to check errors in tests
*/
func ErrorContains(out error, want error) bool {
	if out == nil {
		return want == nil
	}
	if want == nil {
		return false
	}
	return strings.Contains(out.Error(), want.Error())
}

/*
Contains a series of errors used throughout the application
*/
var (
	ErrBuildingConvertTrack      = errors.New("error building convert track")
	ErrClosestDirUnknown         = errors.New("something went very wrong getting the closest dir") // Should never happen
	ErrConvertedFileExists       = errors.New("converted file already exists")
	ErrConvertingTrack           = errors.New("error converting track")
	ErrConvertTrack              = errors.New("error converting track")
	ErrConvertTrackEmpty         = errors.New("convert track is empty")
	ErrFileAlreadyProcessed      = errors.New("file has already been processed")
	ErrGettingClosestDir         = errors.New("error getting closest dir")
	ErrGettingListableURI        = errors.New("error getting listable URI")
	ErrIndexOutOfBounds          = errors.New("index out of bounds")
	ErrInDirPathRequired         = errors.New("InDirPath required")
	ErrInFilePathRequired        = errors.New("InFilePath required")
	ErrMissingRequiredFields     = errors.New("missing required fields")
	ErrNoDirPath                 = errors.New("no directory path found")
	ErrNoFileExtension           = errors.New("no file extension found")
	ErrNoFileName                = errors.New("no file name found")
	ErrNoMatchesFound            = errors.New("no matches found")
	ErrOperationFinished         = errors.New("operation finished")
	ErrOperationNotFound         = errors.New("operation not found")
	ErrBusyPleaseFinishFirst     = errors.New("please finish what you're doing first")
	ErrUserStoppedProcess        = errors.New("user stopped process")
	ErrBuildingStemTrack         = errors.New("error building stem track")
	ErrStemOutputExists          = errors.New("stem extraction output already exists")
	ErrStemTrackEmpty            = errors.New("stem track is empty")
	ErrDemucsSepStep             = errors.New("error running demucs seperation step")
	ErrMergeM4AStep              = errors.New("error running merge m4a step")
	ErrAddMetadataStep           = errors.New("error running add metadata step")
	ErrCleanupStep               = errors.New("error running cleanup step")
	ErrInvalidStemSeparationType = errors.New("invalid stem separation type")
	ErrInvalidPlatform           = errors.New("platform is invalid")
	ErrConfigDoesNotExist        = errors.New("config does not exist")
	ErrMissingPlaylistURL        = errors.New("missing playlist URL")
	ErrExtractingHydrationString = errors.New("error extracting hydration string")
	ErrHydratableKeyNotFound     = errors.New("hydratable key not found")
	ErrTrackMissingID            = errors.New("track missing ID")
	ErrPlaylistAlreadyExists     = errors.New("playlist already exists")
	ErrMissingRedirectURI        = errors.New("missing redirect URI")
	ErrRequestingPlaylist        = errors.New("error requesting playlist")
	ErrUserCancelled             = errors.New("user cancelled operation")
	ErrSoundCloudClientIDNotSet  = errors.New("soundcloud client ID not set")
	ErrPleaseWaitForDownload     = errors.New("please wait for the track to finish downloading")
)

var (
	GenErrBuildingConvertTrack = func(path string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrBuildingConvertTrack, path, err)
	}
	GenErrConvertingTrack = func(name string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrConvertingTrack, name, err)
	}
	GenErrClosestDirUnknown = func(path string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrClosestDirUnknown, path, err)
	}
	GenErrGettingClosestDir = func(err error) error {
		return fmt.Errorf("%s: %w", ErrGettingClosestDir, err)
	}
	GenErrGettingListableURI = func(err error) error {
		return fmt.Errorf("%s: %w", ErrGettingListableURI, err)
	}
	GenErrConvertTrack = func(name string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrConvertTrack, name, err)
	}
	GenErrBuildingStemTrack = func(name string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrBuildingStemTrack, name, err)
	}
	GenErrDemucsSepStep = func(name string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrDemucsSepStep, name, err)
	}
	GenErrMergeM4AStep = func(name string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrMergeM4AStep, name, err)
	}
	GenErrAddMetadataStep = func(name string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrAddMetadataStep, name, err)
	}
	GenErrCleanupStep = func(name string, err error) error {
		return fmt.Errorf("%s %s: %w", ErrCleanupStep, name, err)
	}
)
