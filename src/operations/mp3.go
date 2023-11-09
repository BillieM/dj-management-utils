package operations

import (
	"context"
	"errors"
	"fmt"

	"github.com/billiem/seren-management/src/helpers"
	"github.com/deliveryhero/pipeline/v2"
)

/*
Gets all of the files in the provided directory which can be converted to mp3

if recursion is true, will also get files in subdirectories
*/
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

	tracksChan := pipeline.Emit(tracks...)

	convertOut := pipeline.ProcessConcurrently(ctx, 1, pipeline.NewProcessor(func(ctx context.Context, t ConvertTrack) (ConvertTrack, error) {
		t, err := convertTrack(t)
		if err != nil {
			return t, err
		}
		completedTracks++
		o.StepCallback(float64(completedTracks) / float64(totalTracks))
		fmt.Println(completedTracks, totalTracks, float64(completedTracks)/float64(totalTracks))
		return t, nil
	}, func(t ConvertTrack, err error) {
		completedTracks++
		fmt.Printf("t.Name: %s failed because: %s\n", t.Name, err.Error())
	}), tracksChan)

	for range convertOut {
		select {
		case <-ctx.Done():
			return
		case t := <-convertOut:
			fmt.Println(t)
			_ = t
		default:
			continue
		}
	}
}

func convertTrack(track ConvertTrack) (ConvertTrack, error) {

	if track == (ConvertTrack{}) {
		return track, errors.New("convert track is empty")
	}

	if helpers.DoesFileExist(track.NewFile.FileInfo.FullPath) {
		return track, errors.New("file already exists: " + track.NewFile.FileInfo.FullPath)
	}

	// create dir for new file if it doesn't exist
	err := helpers.CreateDirIfNotExists(track.NewFile.FileInfo.DirPath)

	if err != nil {
		return track, err
	}

	err = helpers.CmdExec(
		"ffmpeg",
		"-i", track.OriginalFile.FileInfo.FullPath,
		"-b:a", "320k",
		track.NewFile.FileInfo.FullPath,
	)

	if err != nil {
		return track, err
	}

	// create dir for old file if it doesn't exist
	err = helpers.CreateDirIfNotExists(track.OriginalFile.FileInfo.DirPath)

	if err != nil {
		return track, err
	}

	// delete the original file if DeleteOnFinish is true
	if track.OriginalFile.DeleteOnFinish {
		err = helpers.CmdExec(
			"rm",
			track.OriginalFile.FileInfo.FullPath,
		)

		if err != nil {
			return track, err
		}
	}

	return track, nil

}
