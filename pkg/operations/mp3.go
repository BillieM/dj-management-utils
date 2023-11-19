package operations

import (
	"context"
	"fmt"
	"strings"

	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/deliveryhero/pipeline/v2"
)

/*
getConvertPath gets all of the files in the provided directory which should be converted to mp3 based on the config

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

	if len(tracks) == 0 {
		return
	}

	var completedTracks int
	var totalTracks = len(tracks)

	tracksChan := pipeline.Emit(tracks...)

	convertOut := pipeline.ProcessConcurrently(ctx, 4, pipeline.NewProcessor(func(ctx context.Context, t ConvertTrack) (ConvertTrack, error) {

		if t == (ConvertTrack{}) {
			return t, helpers.ErrConvertTrackEmpty
		}

		o.StepCallback(stepStartedStepInfo(fmt.Sprintf("Converting: %s", t.Name)))

		t, err := convertTrack(t)
		if err != nil {
			return t, err
		}
		completedTracks++

		o.StepCallback(stepFinishedStepInfo(
			fmt.Sprintf("Converted: %s", t.Name),
			float64(completedTracks)/float64(totalTracks),
		))

		return t, nil
	}, func(t ConvertTrack, err error) {
		completedTracks++
		if strings.Contains(err.Error(), "context canceled") {
			o.StepCallback(
				progressOnlyStepInfo(float64(completedTracks) / float64(totalTracks)),
			)
		} else {
			o.StepCallback(stepWarningStepInfo(
				helpers.GenErrConvertTrack(t.Name, err),
				float64(completedTracks)/float64(totalTracks),
			))
		}
	}), tracksChan)

	for range convertOut {
		t := <-convertOut
		_ = t
	}
}

func convertTrack(track ConvertTrack) (ConvertTrack, error) {

	// create dir for new file if it doesn't exist
	err := helpers.CreateDirIfNotExists(track.NewFile.FileInfo.DirPath)

	if err != nil {
		return track, err
	}

	_, err = helpers.CmdExec(
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
		_, err = helpers.CmdExec(
			"rm",
			track.OriginalFile.FileInfo.FullPath,
		)

		if err != nil {
			return track, err
		}
	}

	return track, nil

}
