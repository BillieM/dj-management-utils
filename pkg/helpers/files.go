package helpers

import (
	"errors"
	"fmt"
	"io/fs"
	"mime"
	"os"
	"path/filepath"

	"github.com/billiem/seren-management/pkg/projectpath"
)

type FileInfo struct {
	FullPath      string
	DirPath       string
	FileName      string
	FileExtension string
}

func (f FileInfo) BuildFullPath() string {
	return JoinFilepathToSlash(f.DirPath, f.FileName+f.FileExtension)
}

func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		return true
	}
	return false
}

/*
Uses regex to check if a string contains any of the given extensions
*/
func IsExtensionInArray(s string, a []string) bool {
	for _, v := range a {
		if RegexContains(s, `(?i)\.`+v+`$`) {
			return true
		}
	}
	return false
}

func ReplaceTrackExtension(s string, r string, a []string) string {
	for _, v := range a {
		s = RegexReplace(s, `(?i)\.`+v+`$`, r)
	}
	return s
}

func GetDirPathFromFilePath(s string) (string, error) {
	dir, _ := filepath.Split(s)
	dir = filepath.ToSlash(dir)
	if dir == "" {
		return "", ErrNoDirPath
	}
	return dir, nil
}

/*
GetFileNameFromFilePath returns the file name from a given file path (without the extension)

TODO: add an argument to include the extension
*/
func GetFileNameFromFilePath(s string) (string, error) {
	_, file := filepath.Split(s)
	if file == "" {
		return "", ErrNoFileName
	}
	fileName := file[:len(file)-len(filepath.Ext(file))]

	if fileName == "" {
		return "", ErrNoFileName
	}

	return fileName, nil
}

/*
Returns the file extension from a given file path (including the dot)
*/
func GetFileExtensionFromFilePath(s string) (string, error) {
	ext := filepath.Ext(s)
	if ext == "" {
		return "", ErrNoFileExtension
	}
	return ext, nil
}

/*
Returns FileInfo struct from a given file path

FileInfo contains information about the file path, including the directory path, file name, and file extension
*/
func SplitFilePath(s string) (FileInfo, error) {

	// errors don't matter in this case
	dirPath, _ := GetDirPathFromFilePath(s)
	fileName, _ := GetFileNameFromFilePath(s)
	fileExtension, _ := GetFileExtensionFromFilePath(s)

	if dirPath == "" && fileName == "" && fileExtension == "" {
		return FileInfo{}, ErrNoMatchesFound
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
		return fileInfo, err
	}
	if fileInfo.DirPath == "" || fileInfo.FileName == "" || fileInfo.FileExtension == "" {
		return fileInfo, ErrMissingRequiredFields
	}
	return fileInfo, nil
}

func GetFilesInDir(dirPath string, recursion bool) ([]string, error) {

	var filePaths []string

	if recursion {
		err := filepath.WalkDir(dirPath, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				filePaths = append(filePaths, path)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
	} else {
		files, err := os.ReadDir(dirPath)
		if err != nil {
			return nil, err
		}
		for _, file := range files {
			filePaths = append(filePaths, fmt.Sprintf("%s/%s", dirPath, file.Name()))
		}
	}

	return filePaths, nil
}

/*
Provides a recursive way of finding the closest directory to a given path,
or the base directory if no 'close directory' is found within 4 recursive calls
If no base directory is found, it will return the root directory (i.e. /)
*/
func GetClosestDir(path string, baseDirPath string, rCnt *int) (string, error) {
	*rCnt++
	fi, err := os.Stat(path)

	if path == "" {
		// if the path is empty, we may aswell just check with baseDirPath/ projectpath.Root
		*rCnt = 5
	}

	if err != nil {
		if *rCnt <= 4 {
			return GetClosestDir(JoinFilepathToSlash(path, ".."), baseDirPath, rCnt)
		} else if *rCnt == 5 {
			return GetClosestDir(baseDirPath, baseDirPath, rCnt)
		} else if *rCnt == 6 {
			return GetClosestDir(projectpath.Root, baseDirPath, rCnt)
		} else {
			return "", GenErrClosestDirUnknown(path, err)
		}
	}
	if fi.IsDir() {
		return path, nil
	} else {
		return GetClosestDir(
			JoinFilepathToSlash(path, ".."),
			baseDirPath, rCnt,
		)
	}
}

/*
 */
func CreateDirIfNotExists(path string) error {
	err := os.MkdirAll(path, os.ModePerm)

	if err != nil {
		return err
	}

	return nil
}

func JoinFilepathToSlash(a ...string) string {
	return filepath.ToSlash(filepath.Join(a...))
}

func RemoveFileExtension(s string) string {
	return s[:len(s)-len(filepath.Ext(s))]
}

/*
GetAbsOrWdPath returns the same path if it is absolute, otherwise it will return the path joined to the current working directory

Returns an empty string if the path is empty
*/
func GetAbsOrWdPath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if filepath.IsAbs(path) {
		return path, nil
	}
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	return JoinFilepathToSlash(cwd, path), nil
}

/*
GetAbsOrProjPath returns the same path if it is absolute, otherwise it will return the path joined to the project root directory

Returns an empty string if the path is empty
*/
func GetAbsOrProjPath(path string) (string, error) {
	if path == "" {
		return "", nil
	}

	if filepath.IsAbs(path) {
		return path, nil
	}
	return JoinFilepathToSlash(projectpath.Root, path), nil
}

/*
GetFileExtensionFromContentType returns the file extension from a given content type
*/
func GetFileExtensionFromContentType(contentType string) (string, error) {
	exts, err := mime.ExtensionsByType(contentType)

	if err != nil {
		return "", err
	}

	if len(exts) == 0 {
		return "", ErrNoFileExtension
	}

	return exts[0], nil
}

/*
GetFileExtensionFromContentDisposition returns the filename (with extension) from a given content disposition
*/
func GetFileNameFromContentDisposition(contentDisposition string) (string, error) {
	_, params, err := mime.ParseMediaType(contentDisposition)

	if err != nil {
		return "", err
	}

	if params["filename"] == "" {
		return "", ErrNoFileExtension
	}

	return params["filename"], nil
}

/*
GetAudioExtensions returns a list of possible audio file extensions.

This is primarily used for generic file filtering by extension
*/
func GetAudioExtensions() []string {
	return []string{
		"mp3",
		"ogg",
		"aac",
		"alac",
		"flac",
		"aif",
		"aiff",
		"wav",
		"aifc",
		"mp4",
		"m4a",
	}
}
