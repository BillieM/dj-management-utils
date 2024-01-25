package internal

import (
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/helpers"
)

type AudioFile struct {
	FileInfo       FileInfo
	DeleteOnFinish bool // If true, this file should be deleted after all operations are performed
}

type FileInfo struct {
	FullPath      string
	DirPath       string
	FileName      string
	FileExtension string
}

func (f FileInfo) BuildFullPath() string {
	return helpers.JoinFilepathToSlash(f.DirPath, f.FileName+f.FileExtension)
}

/*
Returns FileInfo struct from a given file path

FileInfo contains information about the file path, including the directory path, file name, and file extension
*/
func SplitFilePath(s string) (FileInfo, error) {

	// errors don't matter in this case
	dirPath, _ := helpers.GetDirPathFromFilePath(s)
	fileName, _ := helpers.GetFileNameFromFilePath(s)
	fileExtension, _ := helpers.GetFileExtensionFromFilePath(s)

	if dirPath == "" && fileName == "" && fileExtension == "" {
		return FileInfo{}, fault.New("no matches found")
	}

	return FileInfo{
		FullPath:      s,
		DirPath:       dirPath,
		FileName:      fileName,
		FileExtension: fileExtension,
	}, nil
}

/*
Calls SplitFilePath but checks all required fields are present
*/
func SplitFilePathRequired(s string) (FileInfo, error) {
	fileInfo, err := SplitFilePath(s)
	if err != nil {
		return fileInfo, fault.Wrap(
			err,
			fmsg.With("error splitting file path"),
		)
	}
	if fileInfo.DirPath == "" || fileInfo.FileName == "" || fileInfo.FileExtension == "" {
		return fileInfo, fault.New("missing required fields")
	}
	return fileInfo, nil
}
