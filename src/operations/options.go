package operations

import "github.com/billiem/seren-management/src/helpers"

/*
SeperateSingleStemOptions is used as a way to pass arguments to SeperateSingleStem
*/
type SeparateSingleStemOpts struct {
	InFilePath string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
}

/*
check checks the options for the SeperateSingleStem operation
*/
func (p SeparateSingleStemOpts) check() error {
	if p.InFilePath == "" {
		return helpers.ErrInFilePathRequired
	}

	return nil
}

/*
SeperateFolderStemOptions contains the options for SeperateFolderStem
*/
type SeparateFolderStemOpts struct {
	InDirPath  string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
	Recursion  bool   // Optional
}

/*
check checks the options for the SeperateFolderStem operation
*/
func (p SeparateFolderStemOpts) check() error {
	if p.InDirPath == "" {
		return helpers.ErrInDirPathRequired
	}

	return nil
}

/*
ConvertSingleMp3Options is used as a way to pass arguments to ConvertSingleMp3
*/
type ConvertSingleMp3Opts struct {
	InFilePath string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
}

/*
check checks the options for the ConvertSingleMp3 operation
*/
func (p ConvertSingleMp3Opts) check() error {
	if p.InFilePath == "" {
		return helpers.ErrInFilePathRequired
	}

	return nil
}

/*
ConvertFolderMp3Options contains the options for ConvertFolderMp3
*/
type ConvertFolderMp3Opts struct {
	InDirPath  string // Mandatory
	OutDirPath string // Optional - if not provided, will use the same dir as the input file
	Recursion  bool   // Optional
}

/*
check checks the options for the ConvertFolderMp3 operation
*/
func (p ConvertFolderMp3Opts) check() error {
	if p.InDirPath == "" {
		return helpers.ErrInDirPathRequired
	}

	return nil
}
