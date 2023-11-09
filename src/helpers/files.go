package helpers

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
)

type FileInfo struct {
	FullPath      string
	DirPath       string
	FileName      string
	FileExtension string
}

func (f FileInfo) BuildFullPath() string {
	return fmt.Sprintf("%s%s%s", f.DirPath, f.FileName, f.FileExtension)
}

func DoesFileExist(path string) bool {
	if _, err := os.Stat(path); !errors.Is(err, os.ErrNotExist) {
		return true
	}
	return false
}

/*
Uses regex to check if a string contains any of the given extensions

TODO: May want to adjust this to use GetExtensionFromFilePath for maintainabilities sake
*/
func IsExtensionInArray(s string, a []string) bool {
	for _, v := range a {
		if regexContains(s, `(?i)\.`+v+`$`) {
			return true
		}
	}
	return false
}

// TODO: do i actually need this anymore, we're using FileInfo instead?
func ReplaceTrackExtension(s string, r string, a []string) string {
	for _, v := range a {
		s = regexReplace(s, `(?i)\.`+v+`$`, r)
	}
	return s
}

func GetDirPathFromFilePath(s string) (string, error) {
	fileInfo, err := SplitFilePath(s)
	if err != nil {
		return "", err
	}
	if fileInfo.DirPath == "" {
		return "", errors.New("no directory path found")
	}
	return fileInfo.DirPath, nil
}

func GetFileNameFromFilePath(s string) (string, error) {
	fileInfo, err := SplitFilePath(s)
	if err != nil {
		return "", err
	}
	if fileInfo.FileName == "" {
		return "", errors.New("no file name found")
	}
	return fileInfo.FileName, nil
}

/*
Returns the file extension from a given file path (including the dot)
*/
func GetFileExtensionFromFilePath(s string) (string, error) {
	fileInfo, err := SplitFilePath(s)
	if err != nil {
		return "", err
	}
	if fileInfo.FileExtension == "" {
		return "", errors.New("no file extension found")
	}
	return fileInfo.FileExtension, nil
}

/*
Returns FileInfo struct from a given file path

FileInfo contains information about the file path, including the directory path, file name, and file extension
*/
func SplitFilePath(s string) (FileInfo, error) {
	re := regexp.MustCompile(`(.*/)?(.*?)(\..*)?$`)
	matches := re.FindStringSubmatch(s)
	if !ContainsNonEmptyString(matches) {
		return FileInfo{}, errors.New("no matches found")
	}

	return FileInfo{
		FullPath:      s,
		DirPath:       matches[1],
		FileName:      matches[2],
		FileExtension: matches[3],
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
		return fileInfo, errors.New("missing required fields")
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
	// fmt.Println(*rCnt, path)
	if err != nil {
		if *rCnt <= 4 {
			return GetClosestDir(filepath.Join(path, ".."), baseDirPath, rCnt)
		} else if *rCnt == 5 {
			return GetClosestDir(baseDirPath, baseDirPath, rCnt)
		} else if *rCnt == 6 {
			return GetClosestDir(filepath.Join("/"), baseDirPath, rCnt)
		} else {
			return "", errors.New("Something went very wrong getting the cloest dir, err: " + err.Error())
		}
	}
	if fi.IsDir() {
		return path, nil
	} else {
		return GetClosestDir(filepath.Join(path, ".."), baseDirPath, rCnt)
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
