package operations

import (
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations/internal"
)

type OperationOptions interface {
	Check() (bool, error)
}

/*
SeperateSingleStemOptions is used as a way to pass arguments to SeperateSingleStem
*/
type SeparateSingleStemOpts struct {
	InFilePath string                      // Mandatory
	OutDirPath string                      // Optional - if not provided, will use the same dir as the input file
	Type       internal.StemSeparationType // Mandatory
}

/*
check checks the options for the SeperateSingleStem operation
*/
func (p SeparateSingleStemOpts) Check() (bool, error) {
	if p.InFilePath == "" {
		return false, helpers.ErrInFilePathRequired
	}
	if !p.Type.Check() {
		return false, helpers.ErrInvalidStemSeparationType
	}

	return true, nil
}

/*
SeperateFolderStemOptions contains the options for SeperateFolderStem
*/
type SeparateFolderStemOpts struct {
	InDirPath  string                      // Mandatory
	OutDirPath string                      // Optional - if not provided, will use the same dir as the input file
	Recursion  bool                        // Optional
	Type       internal.StemSeparationType // Mandatory
}

/*
check checks the options for the SeperateFolderStem operation
*/
func (p SeparateFolderStemOpts) Check() (bool, error) {
	if p.InDirPath == "" {
		return false, helpers.ErrInDirPathRequired
	}

	return true, nil
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
func (p ConvertSingleMp3Opts) Check() (bool, error) {
	if p.InFilePath == "" {
		return false, helpers.ErrInFilePathRequired
	}

	return true, nil
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
func (p ConvertFolderMp3Opts) Check() (bool, error) {
	if p.InDirPath == "" {
		return false, helpers.ErrInDirPathRequired
	}

	return true, nil
}

/*
GetSoundCloudPlaylistOpts contains the options for GetSoundCloudPlaylist
*/
type GetSoundCloudPlaylistOpts struct {
	PlaylistURL string // Mandatory
	Refresh     bool   // Optional
}

/*
check checks the options for the GetSoundCloudPlaylist operation
*/
func (p GetSoundCloudPlaylistOpts) Check() (bool, error) {
	if p.PlaylistURL == "" {
		return false, helpers.ErrMissingPlaylistURL
	}

	// perform some regex mapping to check the url is correct

	return true, nil
}
