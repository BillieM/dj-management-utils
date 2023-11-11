package operations

import (
	"context"
	"fmt"

	"github.com/billiem/seren-management/src/helpers"
)

/*
This file serves as an entrypoint for all operations
*/

/*
OperationProcess is used to provide callbacks to the operations package
*/
type OperationProcess interface {
	StepCallback(StepInfo)
	ExitCallback()
}

/*
StepInfo is returned to the StepCallback after each step

It provides information about the step that just finished
*/
type StepInfo struct {
	Progress   float64
	Message    string
	Error      error
	Importance helpers.Importance
}

/*
ConvertSingleMp3Params is used as a way to pass arguments to ConvertSingleMp3
*/
type ConvertSingleMp3Params struct {
	InFilePath string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
}

func (p ConvertSingleMp3Params) check() error {
	if p.InFilePath == "" {
		return helpers.ErrInFilePathRequired
	}

	return nil
}

func ConvertSingleMp3(ctx context.Context, cfg helpers.Config, o OperationProcess, params ConvertSingleMp3Params) {

	err := params.check()

	if err != nil {
		panic(err)
	}

	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray([]string{params.InFilePath}, params.OutDirPath)

	_, _, _ = convertTrackArray, alreadyExistsCnt, errs
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
		return helpers.ErrInDirPathRequired
	}

	return nil
}

func ConvertFolderMp3(ctx context.Context, cfg helpers.Config, o OperationProcess, params ConvertFolderMp3Params) {

	ctx, cancelCauseFunc := context.WithCancelCause(ctx)
	defer cancelCauseFunc(helpers.ErrOperationFinished)

	err := params.check()

	if err != nil {
		cancelCauseFunc(err)
	}

	o.StepCallback(stageStepInfo("Finding files to convert"))
	convertFilePaths, err := getConvertPaths(cfg, params.InDirPath, params.Recursion)
	o.StepCallback(stageStepInfo(fmt.Sprintf("Found %v potential files to convert", len(convertFilePaths))))

	if err != nil {
		// TODO: unsure if this function is going to return errors we still want to process? look into further
		cancelCauseFunc(err)
	}

	o.StepCallback(stageStepInfo("Checking found files"))
	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray(convertFilePaths, params.OutDirPath)
	o.StepCallback(stageStepInfo(fmt.Sprintf("%v files already exist, %v left to convert", alreadyExistsCnt, len(convertTrackArray))))

	for _, err := range errs {
		o.StepCallback(warningStepInfo(err))
	}

	o.StepCallback(stageStepInfo("Converting files to mp3"))
	parallelProcessConvertTrackArray(ctx, o, convertTrackArray)

	o.StepCallback(processFinishedStepInfo("Finished converting files to mp3"))

	o.ExitCallback()
}
