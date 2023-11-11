package helpers

import (
	"errors"
	"fmt"
)

/*
Contains a series of errors used throughout the application
*/
var (
	ErrFileAlreadyProcessed  = errors.New("file has already been processed")
	ErrMissingRequiredFields = errors.New("missing required fields")
	ErrNoMatchesFound        = errors.New("no matches found")
	ErrNoDirPath             = errors.New("no directory path found")
	ErrNoFileName            = errors.New("no file name found")
	ErrNoFileExtension       = errors.New("no file extension found")
	ErrClosestDirUnknown     = errors.New("something went very wrong getting the closest dir") // Should never happen
	ErrPleaseWaitForProcess  = errors.New("please wait for the current process to finish")
	ErrIndexOutOfBounds      = errors.New("index out of bounds")
	ErrUserStoppedProcess    = errors.New("user stopped process")
	ErrConvertedFileExists   = errors.New("converted file already exists")
	ErrConvertTrackEmpty     = errors.New("convert track is empty")
	ErrBuildingConvertTrack  = errors.New("error building convert track")
	ErrConvertingTrack       = errors.New("error converting track")
	ErrInFilePathRequired    = errors.New("InFilePath required")
	ErrOperationFinished     = errors.New("operation finished")
	ErrInDirPathRequired     = errors.New("InDirPath required")
	ErrGettingClosestDir     = errors.New("error getting closest dir")
	ErrGettingListableURI    = errors.New("error getting listable URI")
	ErrOperationNotFound     = errors.New("operation not found")
	ErrConvertTrack          = errors.New("error converting track")
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
)

func HandleFatalError(err error) {
	fmt.Println(err)
	panic(err)
}
