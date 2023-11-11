package operations

import (
	"github.com/billiem/seren-management/src/helpers"
)

/*
Shared between operation track types
*/

type Track struct {
	Name string
}

type AudioFile struct {
	FileInfo       helpers.FileInfo
	DeleteOnFinish bool // If true, this file should be deleted after all operations are performed
}

/*
ConvertTrack is used as part of the process for converting non-mp3 audio files to mp3
*/

type ConvertTrack struct {
	Track
	OriginalFile AudioFile
	NewFile      AudioFile
}

/*
buildConvertTrackArray builds an array of ConvertTrack structs from an array of file paths

File paths have been pre-validated to ensure they are valid files which can be converted
*/
func buildConvertTrackArray(paths []string, outDirPath string) ([]ConvertTrack, int, []error) {
	var tracks []ConvertTrack
	var errs []error
	var alreadyExistsCnt int

	for _, path := range paths {
		track, err := buildConvertTrack(path, outDirPath)

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
func buildConvertTrack(path string, outDirPath string) (ConvertTrack, error) {

	origFileInfo, err := helpers.SplitFilePathRequired(path)

	if err != nil {
		return ConvertTrack{}, err
	}

	newFileInfo := origFileInfo

	// Populate info for the new file
	newFileInfo.FileExtension = ".mp3"
	if outDirPath != "" {
		// TODO: option to preserve folder structure if call was recursive
		// currently, recursive calls will only preserve the folder structure
		// if outDirPath is not provided
		newFileInfo.DirPath = outDirPath
	}
	newFileInfo.FullPath = newFileInfo.BuildFullPath()

	if helpers.DoesFileExist(newFileInfo.FullPath) {
		return ConvertTrack{}, helpers.ErrConvertedFileExists
	}

	return ConvertTrack{
		Track: Track{
			Name: origFileInfo.FileName,
		},
		OriginalFile: AudioFile{
			FileInfo: origFileInfo,
		},
		NewFile: AudioFile{
			FileInfo: newFileInfo,
		},
	}, nil
}

/*
StemTrack is used as part of the process for converting audio files into stems
*/
type StemTrack struct {
	Track

	stemFiles []StemFile
}

type StemFile struct {
	AudioFile
}
