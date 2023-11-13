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
	SkipLog    bool
	Progress   float64
	Message    string
	Error      error
	Importance helpers.Importance
}

func SeparateSingleStem(ctx context.Context, cfg helpers.Config, o OperationProcess, opts SeparateSingleStemOpts) {

	err := opts.check()

	if err != nil {
		panic(err)
	}

}

func SeparateFolderStem(ctx context.Context, cfg helpers.Config, o OperationProcess, opts SeparateFolderStemOpts) {

	err := opts.check()

	if err != nil {
		panic(err)
	}

}

func ConvertSingleMp3(ctx context.Context, cfg helpers.Config, o OperationProcess, opts ConvertSingleMp3Opts) {

	err := opts.check()

	if err != nil {
		panic(err)
	}

	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray([]string{opts.InFilePath}, opts.OutDirPath)

	_, _, _ = convertTrackArray, alreadyExistsCnt, errs
}

/*
ConvertFolderMp3 converts all files in a folder to mp3
*/
func ConvertFolderMp3(ctx context.Context, cfg helpers.Config, o OperationProcess, opts ConvertFolderMp3Opts) {

	err := opts.check()

	if err != nil {
		o.StepCallback(dangerStepInfo(err))
		return
	}

	o.StepCallback(stageStepInfo("Finding files to convert"))
	convertFilePaths, err := getConvertPaths(cfg, opts.InDirPath, opts.Recursion)
	o.StepCallback(stageStepInfo(fmt.Sprintf("Found %v potential files to convert", len(convertFilePaths))))

	if err != nil {
		o.StepCallback(dangerStepInfo(err))
		return
	}

	o.StepCallback(stageStepInfo("Checking found files"))
	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray(convertFilePaths, opts.OutDirPath)
	o.StepCallback(stageStepInfo(fmt.Sprintf("%v files already exist, %v left to convert", alreadyExistsCnt, len(convertTrackArray))))

	for _, err := range errs {
		o.StepCallback(warningStepInfo(err))
	}

	o.StepCallback(stageStepInfo("Converting files to mp3"))
	parallelProcessConvertTrackArray(ctx, o, convertTrackArray)

	o.StepCallback(processFinishedStepInfo("Finished converting files to mp3"))

	o.ExitCallback()
}
