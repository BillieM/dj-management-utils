package operations

import "github.com/billiem/seren-management/src/helpers"

// Gets all of the files in the given dirpath
func (o ConvertFolderMp3Params) getConvertPaths() ([]string, error) {
	convertPaths, err := helpers.GetFilesInDir(o.InDirPath, o.Recursion)
	if err != nil {
		return nil, err
	}
	var validConvertPaths []string
	for _, path := range convertPaths {
		if helpers.IsExtensionInArray(path, o.Config.ExtensionsToConvertToMp3) {
			validConvertPaths = append(validConvertPaths, path)
		}
	}
	return validConvertPaths, nil
}

// func (d *Data) convertMp3(track Track) Track {

// 	if track == (Track{}) {
// 		return track
// 	}

// 	if helpers.doesFileExist(track.NewPath) {
// 		fmt.Printf("file already exists: %s, skipping\n", track.NewPath)
// 		return Track{}
// 	}

// 	fmt.Printf("converting %s to %s\n", track.OriginalPath, track.NewPath)

// 	err := cmdExec(
// 		"ffmpeg",
// 		"-i", track.OriginalPath,
// 		"-b:a", "320k",
// 		track.NewPath,
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("converted %s to %s\n", track.OriginalPath, track.NewPath)

// 	// create dir for old file if it doesn't exist
// 	os.MkdirAll(track.StorageDir, os.ModePerm)

// 	// move the original file to the newpathforold

// 	err = cmdExec(
// 		"mv",
// 		track.OriginalPath,
// 		track.NewPathForOld,
// 	)

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Printf("moved %s to %s\n", track.OriginalPath, track.NewPathForOld)

// 	return track
// }
