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
StepHandler is used to provide callbacks to the operations package
*/
type StepHandler interface {
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
func (e *OpEnv) SeparateSingleStem(ctx context.Context, opts SeparateSingleStemOpts) {

	defer func() {
		e.step(progressOnlyStepInfo(1))
		e.exit()
	}()

	_, err := opts.Check()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Checking file to separate"))
	stemTrackArray, alreadyExistsCnt, errs := buildStemTrackArray([]string{opts.InFilePath}, opts.OutDirPath, opts.Type)

	if len(errs) > 0 {
		e.step(warningStepInfo(errs[0]))
		return
	}

	if alreadyExistsCnt > 0 {
		e.step(warningStepInfo(helpers.ErrConvertedFileExists))
		return
	}

	e.step(stageStepInfo("Converting file to stems"))
	e.parallelProcessStemTrackArray(ctx, stemTrackArray)
	e.step(processFinishedStepInfo("Finished"))
}

/*
SeparateFolderStem separates stems from all files in a folder
*/
func (e *OpEnv) SeparateFolderStem(ctx context.Context, opts SeparateFolderStemOpts) {

	defer func() {
		e.step(progressOnlyStepInfo(1))
		e.exit()
	}()

	_, err := opts.Check()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Finding files to convert"))
	stemFilePaths, err := e.getStemPaths(opts.InDirPath, opts.Recursion)
	e.step(stageStepInfo(fmt.Sprintf("Found %v potential files to convert", len(stemFilePaths))))

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Checking found files"))
	stemTrackArray, alreadyExistsCnt, errs := buildStemTrackArray(stemFilePaths, opts.OutDirPath, opts.Type)
	e.step(stageStepInfo(fmt.Sprintf("%v files already exist, %v left to convert", alreadyExistsCnt, len(stemTrackArray))))

	for _, err := range errs {
		e.step(warningStepInfo(err))
	}

	e.step(stageStepInfo("Converting files to stems"))
	e.parallelProcessStemTrackArray(ctx, stemTrackArray)
	e.step(processFinishedStepInfo("Finished"))
}

/*
ConvertSingleMp3 converts a single file to mp3
*/
func (e *OpEnv) ConvertSingleMp3(ctx context.Context, opts ConvertSingleMp3Opts) {

	defer func() {
		e.step(progressOnlyStepInfo(1))
		e.exit()
	}()

	_, err := opts.Check()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Checking file to convert"))
	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray([]string{opts.InFilePath}, opts.OutDirPath)

	if len(errs) > 0 {
		e.step(warningStepInfo(errs[0]))
		return
	}

	if alreadyExistsCnt > 0 {
		e.step(warningStepInfo(helpers.ErrConvertedFileExists))
		return
	}

	e.step(stageStepInfo("Converting file to mp3"))
	e.parallelProcessConvertTrackArray(ctx, convertTrackArray)
	e.step(processFinishedStepInfo("Finished"))
}

/*
ConvertFolderMp3 converts all files in a folder to mp3
*/
func (e *OpEnv) ConvertFolderMp3(ctx context.Context, opts ConvertFolderMp3Opts) {

	defer func() {
		e.step(progressOnlyStepInfo(1))
		e.exit()
	}()

	_, err := opts.Check()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Finding files to convert"))
	convertFilePaths, err := e.getConvertPaths(opts.InDirPath, opts.Recursion)
	e.step(stageStepInfo(fmt.Sprintf("Found %v potential files to convert", len(convertFilePaths))))

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}

	e.step(stageStepInfo("Checking found files"))
	convertTrackArray, alreadyExistsCnt, errs := buildConvertTrackArray(convertFilePaths, opts.OutDirPath)
	e.step(stageStepInfo(fmt.Sprintf("%v files already exist, %v left to convert", alreadyExistsCnt, len(convertTrackArray))))

	for _, err := range errs {
		e.step(warningStepInfo(err))
	}

	e.step(stageStepInfo("Converting files to mp3"))
	e.parallelProcessConvertTrackArray(ctx, convertTrackArray)
	e.step(processFinishedStepInfo("Finished"))
}

/*
ReadCollection reads a collection for a given platform and stores it in the database
*/
func (e *OpEnv) ReadCollection(ctx context.Context, opts ReadCollectionOpts) {

	collection := opts.Build(e.Config)

	err := collection.ReadCollection()

	if err != nil {
		e.step(dangerStepInfo(err))
		return
	}
}
