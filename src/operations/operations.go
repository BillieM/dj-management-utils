package operations

import (
	"context"

	"github.com/billiem/seren-management/src/helpers"
	"github.com/k0kubun/pp"
)

/*
This file serves as an entrypoint for all operations
*/

type BaseOperationParams struct {
	Config             *helpers.Config
	OperationShortName string
	Context            context.Context
	StepCallback       func(float64)

	Steps []func(Track) (Track, error)
}

/*
Converts a single (non-mp3) file to mp3

Files which can be converted are found in config.ExtensionsToConvertToMp3
*/

type ConvertSingleMp3Params struct {
	BaseOperationParams

	InFilePath string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
}

func (o ConvertSingleMp3Params) ExecuteOperation() {
	o.OperationShortName = "convert-single-mp3"

	convertTrackArray, errors := buildConvertTrackArray([]string{o.InFilePath}, o.OutDirPath)

	_ = errors
	_ = convertTrackArray
}

/*
Converts a folder of (non-mp3) files to mp3

Files which can be converted are found in config.ExtensionsToConvertToMp3
*/
type ConvertFolderMp3Params struct {
	BaseOperationParams

	InDirPath  string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
	Recursion  bool   // Optional
}

func (o ConvertFolderMp3Params) ExecuteOperation() {
	o.OperationShortName = "convert-folder-mp3"

	convertFilePaths, err := o.getConvertPaths()

	if err != nil {
		panic(err)
	}

	convertTrackArray, errors := buildConvertTrackArray(convertFilePaths, o.OutDirPath)

	for _, track := range convertTrackArray {
		pp.Println(track)
	}

	for _, err := range errors {
		pp.Println(err)
	}

	_ = errors
	_ = convertTrackArray
}

// func ConvertFolderMp3(o OperationParams) {

// 	// pipeConfig := parapipe.Config{
// 	// 	ProcessErrors: false,
// 	// }
// 	// pipe := parapipe.NewPipeline(pipeConfig).
// 	// 	Pipe(1, func(v interface{}) interface{} {
// 	// 		inputVal := v.(string)
// 	// 	}).
// 	// 	Pipe(1, func(v interface{}) interface{} {
// 	// 		inputVal := v.(string)
// 	// 	})

// 	// for result := range pipe.Out() {
// 	// 	fmt.Println(result)
// 	// }
// }
