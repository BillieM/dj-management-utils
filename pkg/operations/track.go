package operations

import (
	"github.com/billiem/seren-management/pkg/helpers"
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
	ID int

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

	origFileInfo, err := helpers.SplitFilePathRequired(path)

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
		ID: id,
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

	ID int

	OriginalFile AudioFile

	OutFile AudioFile // this is the .stem.m4a used by Traktor

	StemDir    string // The directory where the stem files will be created
	SkipDemucs bool   // If true, skip the demucs step (i.e. stem files exist on Traktor type)
	StemsOnly  bool   // If true, skip the merge/ metadata steps (i.e. only stems are required)

	BassFile   StemFile
	DrumsFile  StemFile
	OtherFile  StemFile
	VocalsFile StemFile
}

type StemFile struct {
	AudioFile
}

/*
buildStemTrackArray builds an array of StemTrack structs from an array of file paths
*/

func buildStemTrackArray(paths []string, outDirPath string, stemType StemSeparationType) ([]StemTrack, int, []error) {
	var tracks []StemTrack
	var errs []error
	var alreadyExistsCnt int

	for i, path := range paths {
		track, err := buildStemTrack(i, path, outDirPath, stemType)

		if err != nil {
			if err == helpers.ErrStemOutputExists {
				alreadyExistsCnt++
				continue
			}

			errs = append(errs, helpers.GenErrBuildingStemTrack(path, err))
			continue
		}

		tracks = append(tracks, track)
	}

	return tracks, alreadyExistsCnt, errs
}

/*
buildStemTrack builds a StemTrack struct from a file path
*/
func buildStemTrack(id int, path string, outDirPath string, stemType StemSeparationType) (StemTrack, error) {

	origFileInfo, err := helpers.SplitFilePathRequired(path)

	if err != nil {
		return StemTrack{}, err
	}

	var newFileInfo helpers.FileInfo

	baseStemDirPath := helpers.JoinFilepathToSlash(origFileInfo.DirPath, origFileInfo.FileName) + "/"
	if outDirPath != "" {
		baseStemDirPath = helpers.JoinFilepathToSlash(outDirPath, origFileInfo.FileName) + "/"
	}

	deleteOnFinish := stemType == Traktor
	var skipDemucs bool
	var stemsOnly bool

	bassFile := buildStemFile(baseStemDirPath, "bass", origFileInfo.FileExtension, deleteOnFinish)
	drumsFile := buildStemFile(baseStemDirPath, "drums", origFileInfo.FileExtension, deleteOnFinish)
	otherFile := buildStemFile(baseStemDirPath, "other", origFileInfo.FileExtension, deleteOnFinish)
	vocalsFile := buildStemFile(baseStemDirPath, "vocals", origFileInfo.FileExtension, deleteOnFinish)

	// Check if the demucs output already exists
	stemsExist := helpers.DoesFileExist(bassFile.FileInfo.FullPath) &&
		helpers.DoesFileExist(drumsFile.FileInfo.FullPath) &&
		helpers.DoesFileExist(otherFile.FileInfo.FullPath) &&
		helpers.DoesFileExist(vocalsFile.FileInfo.FullPath)

	// Build the out file only if generating a Traktor stem file (out file is the .stem.m4a used by Traktor)
	if stemType == Traktor {
		newFileInfo = origFileInfo
		newFileInfo.FileExtension = ".stem.m4a"
		if outDirPath != "" {
			newFileInfo.DirPath = outDirPath
		}
		newFileInfo.FullPath = newFileInfo.BuildFullPath()

		if stemsExist {
			skipDemucs = true
		}

		if helpers.DoesFileExist(newFileInfo.FullPath) {
			return StemTrack{}, helpers.ErrStemOutputExists
		}
	} else if stemType == FourTrack {

		stemsOnly = true

		if stemsExist {
			return StemTrack{}, helpers.ErrStemOutputExists
		}
	}

	return StemTrack{
		ID: id,
		Track: Track{
			Name: origFileInfo.FileName,
		},
		OriginalFile: AudioFile{
			FileInfo: origFileInfo,
		},
		OutFile: AudioFile{
			FileInfo: newFileInfo,
		},
		StemDir:    baseStemDirPath,
		SkipDemucs: skipDemucs,
		StemsOnly:  stemsOnly,
		BassFile:   bassFile,
		DrumsFile:  drumsFile,
		OtherFile:  otherFile,
		VocalsFile: vocalsFile,
	}, nil
}

func buildStemFile(baseStemDirPath string, fileName string, extension string, deleteOnFinish bool) StemFile {

	stemFileInfo := helpers.FileInfo{
		DirPath:       baseStemDirPath,
		FileName:      fileName,
		FileExtension: extension,
	}

	stemFileInfo.FullPath = stemFileInfo.BuildFullPath()

	return StemFile{
		AudioFile: AudioFile{
			FileInfo:       stemFileInfo,
			DeleteOnFinish: deleteOnFinish,
		},
	}
}
