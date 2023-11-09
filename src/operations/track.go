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

type ProcessTrack interface {
}

/*
ConvertTrack is used as part of the process for converting non-mp3 audio files to mp3
*/

type ConvertTrack struct {
	Track
	OriginalFile AudioFile
	NewFile      AudioFile
}

func buildConvertTrackArray(paths []string, outDirPath string) ([]ConvertTrack, []error) {
	var tracks []ConvertTrack
	var errors []error

	for _, path := range paths {
		track, err := buildConvertTrack(path, outDirPath)

		if err != nil {
			errors = append(errors, err)
			continue
		}

		tracks = append(tracks, track)
	}

	return tracks, errors

}

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
