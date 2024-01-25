package operations

import (
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations/internal"
)

/*
ConvertTrack is used as part of the process for converting non-mp3 audio files to mp3

Must be exported to be used in github.com/deliveryhero/pipeline/v2
*/

type ConvertTrack struct {
	ID int

	Name         string
	OriginalFile internal.AudioFile
	NewFile      internal.AudioFile
}

/*
GetConvertPaths gets all of the files in the provided directory which should be converted to mp3 based on the config

if recursion is true, will also get files in subdirectories
*/
func (e *Mp3Env) GetConvertPaths(inDirPath string, recursion bool) ([]string, error) {
	convertPaths, err := helpers.GetFilesInDir(inDirPath, recursion)
	if err != nil {
		return nil, err
	}
	var validConvertPaths []string
	for _, path := range convertPaths {
		if helpers.IsExtensionInArray(path, e.Config.ExtensionsToConvertToMp3) {
			validConvertPaths = append(validConvertPaths, path)
		}
	}
	return validConvertPaths, nil
}

/*
BuildConvertTracks builds an array of ConvertTrack structs from an array of file paths

File paths have been pre-validated to ensure they are valid files which can be converted
by the GetConvertPaths function
*/
func BuildConvertTracks(paths []string, outDirPath string) ([]ConvertTrack, int, []error) {
	var tracks []ConvertTrack
	var errs []error
	var alreadyExistsCnt int

	for i, path := range paths {
		track, err := buildConvertTrack(i, path, outDirPath)

		if err != nil {
			if err == helpers.ErrConvertedFileExists {
				alreadyExistsCnt++
				continue
			}

			errs = append(errs, helpers.GenErrBuildingConvertTrack(path, err))
			continue
		}

		tracks = append(tracks, track)
	}

	return tracks, alreadyExistsCnt, errs

}

/*
buildConvertTrack builds a ConvertTrack struct from a file path

File path has been pre-validated to ensure it is a valid file which can be converted
*/
func buildConvertTrack(id int, path string, outDirPath string) (ConvertTrack, error) {

	origFileInfo, err := internal.SplitFilePathRequired(path)

	if err != nil {
		return ConvertTrack{}, err
	}

	newFileInfo := origFileInfo

	// Populate info for the new file
	newFileInfo.FileExtension = ".mp3"
	if outDirPath != "" {
		newFileInfo.DirPath = outDirPath
	}
	newFileInfo.FullPath = newFileInfo.BuildFullPath()

	if helpers.DoesFileExist(newFileInfo.FullPath) {
		return ConvertTrack{}, helpers.ErrConvertedFileExists
	}

	return ConvertTrack{
		ID:   id,
		Name: origFileInfo.FileName,
		OriginalFile: internal.AudioFile{
			FileInfo: origFileInfo,
		},
		NewFile: internal.AudioFile{
			FileInfo: newFileInfo,
		},
	}, nil
}
