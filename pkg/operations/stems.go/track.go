package operations

import (
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations/internal"
)

/*
StemTrack is used as part of the process for converting audio files into stems

Must be exported to be used in github.com/deliveryhero/pipeline/v2
*/
type StemTrack struct {
	ID   int
	Name string

	OriginalFile internal.AudioFile

	OutFile internal.AudioFile // this is the .stem.m4a used by Traktor

	StemDir    string // The directory where the stem files will be created
	SkipDemucs bool   // If true, skip the demucs step (i.e. stem files exist on Traktor type)
	StemsOnly  bool   // If true, skip the merge/ metadata steps (i.e. only stems are required)

	BassFile   StemFile
	DrumsFile  StemFile
	OtherFile  StemFile
	VocalsFile StemFile
}

type StemFile struct {
	internal.AudioFile
}

/*
buildStemTrackArray builds an array of StemTrack structs from an array of file paths
*/

func buildStemTrackArray(paths []string, outDirPath string, stemType internal.StemSeparationType) ([]StemTrack, int, []error) {
	var tracks []StemTrack
	var errs []error
	var alreadyExistsCnt int

	for i, path := range paths {
		track, err := BuildStemTrack(i, path, outDirPath, stemType)

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
BuildStemTrack builds a StemTrack struct from a file path
*/
func BuildStemTrack(id int, path string, outDirPath string, stemType internal.StemSeparationType) (StemTrack, error) {

	origFileInfo, err := internal.SplitFilePathRequired(path)

	if err != nil {
		return StemTrack{}, err
	}

	var newFileInfo internal.FileInfo

	baseStemDirPath := helpers.JoinFilepathToSlash(origFileInfo.DirPath, origFileInfo.FileName) + "/"
	if outDirPath != "" {
		baseStemDirPath = helpers.JoinFilepathToSlash(outDirPath, origFileInfo.FileName) + "/"
	}

	deleteOnFinish := stemType == internal.Traktor
	var skipDemucs bool
	var stemsOnly bool

	bassFile := BuildStemFile(baseStemDirPath, "bass", origFileInfo.FileExtension, deleteOnFinish)
	drumsFile := BuildStemFile(baseStemDirPath, "drums", origFileInfo.FileExtension, deleteOnFinish)
	otherFile := BuildStemFile(baseStemDirPath, "other", origFileInfo.FileExtension, deleteOnFinish)
	vocalsFile := BuildStemFile(baseStemDirPath, "vocals", origFileInfo.FileExtension, deleteOnFinish)

	// Check if the demucs output already exists
	stemsExist := helpers.DoesFileExist(bassFile.FileInfo.FullPath) &&
		helpers.DoesFileExist(drumsFile.FileInfo.FullPath) &&
		helpers.DoesFileExist(otherFile.FileInfo.FullPath) &&
		helpers.DoesFileExist(vocalsFile.FileInfo.FullPath)

	// Build the out file only if generating a Traktor stem file (out file is the .stem.m4a used by Traktor)
	if stemType == internal.Traktor {
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
	} else if stemType == internal.FourTrack {

		stemsOnly = true

		if stemsExist {
			return StemTrack{}, helpers.ErrStemOutputExists
		}
	}

	return StemTrack{
		ID:   id,
		Name: origFileInfo.FileName,
		OriginalFile: internal.AudioFile{
			FileInfo: origFileInfo,
		},
		OutFile: internal.AudioFile{
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

func BuildStemFile(baseStemDirPath string, fileName string, extension string, deleteOnFinish bool) StemFile {

	stemFileInfo := internal.FileInfo{
		DirPath:       baseStemDirPath,
		FileName:      fileName,
		FileExtension: extension,
	}

	stemFileInfo.FullPath = stemFileInfo.BuildFullPath()

	return StemFile{
		AudioFile: internal.AudioFile{
			FileInfo:       stemFileInfo,
			DeleteOnFinish: deleteOnFinish,
		},
	}
}
