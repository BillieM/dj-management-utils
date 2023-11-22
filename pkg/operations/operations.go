package operations

import (
	"context"
	"fmt"

	"github.com/billiem/seren-management/pkg/helpers"
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

/*
SeperateSingleStem separates stems from a single file
*/
func SeparateSingleStem(ctx context.Context, env OpEnv, o OperationProcess, opts SeparateSingleStemOpts) {

	defer func() {
		o.StepCallback(progressOnlyStepInfo(1))
		o.ExitCallback()
	}()

	_, err := opts.Check()

	if err != nil {
		o.StepCallback(dangerStepInfo(err))
		return
	}

	o.StepCallback(stageStepInfo("Checking file to separate"))
	stemTrackArray, alreadyExistsCnt, errs := buildStemTrackArray([]string{opts.InFilePath}, opts.OutDirPath, opts.Type)

	if len(errs) > 0 {
		o.StepCallback(warningStepInfo(errs[0]))
		return
	}

	if alreadyExistsCnt > 0 {
		o.StepCallback(warningStepInfo(helpers.ErrConvertedFileExists))
		return
	}

	o.StepCallback(stageStepInfo("Converting file to stems"))
	parallelProcessStemTrackArray(ctx, o, stemTrackArray)
	o.StepCallback(processFinishedStepInfo("Finished"))
}

/*
SeparateFolderStem separates stems from all files in a folder
*/
func SeparateFolderStem(ctx context.Context, env OpEnv, o OperationProcess, opts SeparateFolderStemOpts) {

	defer func() {
		o.StepCallback(progressOnlyStepInfo(1))
		o.ExitCallback()
	}()

	_, err := opts.Check()

	if err != nil {
		o.StepCallback(dangerStepInfo(err))
		return
	}

	o.StepCallback(stageStepInfo("Finding files to convert"))
	stemFilePaths, err := getStemPaths(opts.InDirPath, opts.Recursion, env.Config.ExtensionsToSeparateToStems)
	o.StepCallback(stageStepInfo(fmt.Sprintf("Found %v potential files to convert", len(stemFilePaths))))

	if err != nil {
		o.StepCallback(dangerStepInfo(err))
		return
	}

	o.StepCallback(stageStepInfo("Checking found files"))
	stemTrackArray, alreadyExistsCnt, errs := buildStemTrackArray(stemFilePaths, opts.OutDirPath, opts.Type)
	o.StepCallback(stageStepInfo(fmt.Sprintf("%v files already exist, %v left to convert", alreadyExistsCnt, len(stemTrackArray))))

	for _, err := range errs {
		o.StepCallback(warningStepInfo(err))
	}

	o.StepCallback(stageStepInfo("Converting files to stems"))
	parallelProcessStemTrackArray(ctx, o, stemTrackArray)
	o.StepCallback(processFinishedStepInfo("Finished"))
}

/*
ConvertSingleMp3 converts a single file to mp3
*/
func ConvertSingleMp3(ctx context.Context, env OpEnv, o OperationProcess, opts ConvertSingleMp3Opts) {

	defer func() {
		o.StepCallback(progressOnlyStepInfo(1))
		o.ExitCallback()
	}()

	_, err := opts.Check()

	if err != nil {
		o.StepCallback(dangerStepInfo(err))
		return
	}

	o.StepCallback(stageStepInfo("Checking file to convert"))
	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray([]string{opts.InFilePath}, opts.OutDirPath)

	if len(errs) > 0 {
		o.StepCallback(warningStepInfo(errs[0]))
		return
	}

	if alreadyExistsCnt > 0 {
		o.StepCallback(warningStepInfo(helpers.ErrConvertedFileExists))
		return
	}

	o.StepCallback(stageStepInfo("Converting file to mp3"))
	parallelProcessConvertTrackArray(ctx, o, convertTrackArray)
	o.StepCallback(processFinishedStepInfo("Finished"))
}

/*
ConvertFolderMp3 converts all files in a folder to mp3
*/
func ConvertFolderMp3(ctx context.Context, env OpEnv, o OperationProcess, opts ConvertFolderMp3Opts) {

	defer func() {
		o.StepCallback(progressOnlyStepInfo(1))
		o.ExitCallback()
	}()

	_, err := opts.Check()

	if err != nil {
		o.StepCallback(dangerStepInfo(err))
		return
	}

	o.StepCallback(stageStepInfo("Finding files to convert"))
	convertFilePaths, err := getConvertPaths(opts.InDirPath, opts.Recursion, env.Config.ExtensionsToConvertToMp3)
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
	o.StepCallback(processFinishedStepInfo("Finished"))
}

/*
ReadCollection reads a collection for a given platform and stores it in the database
*/
func ReadCollection(ctx context.Context, env OpEnv, o OperationProcess, opts ReadCollectionOpts) {

	collection := opts.Build(env.Config)

	collection.ReadCollection()
}
