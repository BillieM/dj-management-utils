package operations

import (
	"github.com/billiem/seren-management/src/helpers"
)

type Track struct {
	name string
}

type AudioFile struct {
	originalFileInfo helpers.FileInfo // File path before any operations are performed (if "", this is a new file)
	newFileInfo      helpers.FileInfo // File path after all operations are performed (if "", this file should not be moved)
	clearOnExit      bool             // If true, this file should be deleted after all operations are performed
	extension        string           // File extension (e.g. ".mp3")
}

/*
ConvertTrack is used as part of the process for converting non-mp3 audio files to mp3
*/

type ConvertTrack struct {
	Track
	originalFile AudioFile
	newFile      AudioFile
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
	if outDirPath != "" {
		newFileInfo.DirPath = outDirPath
	}

	return ConvertTrack{
		Track: Track{
			name: origFileInfo.FileName,
		},
		originalFile: AudioFile{
			originalFileInfo: origFileInfo,
		},
		newFile: AudioFile{
			newFileInfo: newFileInfo,
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
