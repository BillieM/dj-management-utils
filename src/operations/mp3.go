package operations

import (
	"context"
	"fmt"
	"time"

	"github.com/billiem/seren-management/src/helpers"
	"github.com/deliveryhero/pipeline/v2"
)

// Gets all of the files in the given dirpath
func getConvertPaths(cfg helpers.Config, inDirPath string, recursion bool) ([]string, error) {
	convertPaths, err := helpers.GetFilesInDir(inDirPath, recursion)
	if err != nil {
		return nil, err
	}
	var validConvertPaths []string
	for _, path := range convertPaths {
		if helpers.IsExtensionInArray(path, cfg.ExtensionsToConvertToMp3) {
			validConvertPaths = append(validConvertPaths, path)
		}
	}
	return validConvertPaths, nil
}

func parallelProcessConvertTrackArray(ctx context.Context, o OperationProcess, tracks []ConvertTrack) {

	var completedTracks int
	var totalTracks = len(tracks)

	tracksChan := pipeline.Delay(ctx, time.Millisecond*500, pipeline.Emit(tracks...))

	convertOut := pipeline.ProcessConcurrently(ctx, 4, pipeline.NewProcessor(func(ctx context.Context, t ConvertTrack) (ConvertTrack, error) {
		completedTracks++
		o.StepCallback(float64(completedTracks) / float64(totalTracks))
		fmt.Println(completedTracks, totalTracks, float64(completedTracks)/float64(totalTracks))
		return t, nil
	}, func(t ConvertTrack, err error) {
		fmt.Printf("t.Name: %s failed because: %s\n", t.Name, err.Error())
	}), tracksChan)

	for range convertOut {
		select {
		case <-ctx.Done():
			return
		case t := <-convertOut:
			fmt.Println(t)
			_ = t
		}
	}
	// for t := range convertOut {
	// 	fmt.Println(t)
	// 	_ = t
	// }
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
