package operations

import (
	"context"
	"errors"

	"github.com/billiem/seren-management/src/helpers"
)

/*
This file serves as an entrypoint for all operations
*/

type OperationProcess interface {
	StepCallback(StepInfo)
	ExitCallback()
}

type StepInfo struct {
	Progress float64
	Message  string
}

/*
Converts a single (non-mp3) file to mp3

Files which can be converted are found in config.ExtensionsToConvertToMp3
*/

type ConvertSingleMp3Params struct {
	InFilePath string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
}

func (p ConvertSingleMp3Params) check() error {
	if p.InFilePath == "" {
		return errors.New("inFilePath is required")
	}

	return nil
}

func ConvertSingleMp3(ctx context.Context, cfg helpers.Config, o OperationProcess, params ConvertSingleMp3Params) {

	err := params.check()

	if err != nil {
		panic(err)
	}

	convertTrackArray, errors := buildConvertTrackArray([]string{params.InFilePath}, params.OutDirPath)

	_ = errors
	_ = convertTrackArray
}

/*
Converts a folder of (non-mp3) files to mp3

Files which can be converted are found in config.ExtensionsToConvertToMp3
*/
type ConvertFolderMp3Params struct {
	InDirPath  string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
	Recursion  bool   // Optional
}

func (p ConvertFolderMp3Params) check() error {
	if p.InDirPath == "" {
		return errors.New("InDirPath is required")
	}

	return nil
}

func ConvertFolderMp3(ctx context.Context, cfg helpers.Config, o OperationProcess, params ConvertFolderMp3Params) {

	ctx, cancelCauseFunc := context.WithCancelCause(ctx)
	defer cancelCauseFunc(errors.New("operation finished"))

	err := params.check()

	if err != nil {
		cancelCauseFunc(err)
	}

	convertFilePaths, err := getConvertPaths(cfg, params.InDirPath, params.Recursion)

	if err != nil {
		// TODO: unsure if this function is going to return errors we still want to process? look into further
		cancelCauseFunc(err)
	}

	convertTrackArray, errs := buildConvertTrackArray(convertFilePaths, params.OutDirPath)

	for _, err := range errs {
		o.StepCallback(StepInfo{
			Progress: 0,
			Message:  err.Error(),
		})
	}

	parallelProcessConvertTrackArray(ctx, o, convertTrackArray)

	o.ExitCallback()
}
