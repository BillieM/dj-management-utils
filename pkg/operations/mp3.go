package operations

import (
	"context"
	"fmt"
	"strings"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/deliveryhero/pipeline/v2"
)

/*
getConvertPath gets all of the files in the provided directory which should be converted to mp3 based on the config

if recursion is true, will also get files in subdirectories
*/
func (e *OpEnv) getConvertPaths(inDirPath string, recursion bool) ([]string, error) {
	convertPaths, err := helpers.GetFilesInDir(inDirPath, recursion)
	if err != nil {
		return nil, err
	}
	var validConvertPaths []string
	for _, path := range convertPaths {
		if helpers.IsExtensionInArray(path, e.Config.ExtensionsToConvertToMp3) {
			validConvertPaths = append(validConvertPaths, path)
		}
	}
	return validConvertPaths, nil
}

func (e *OpEnv) parallelProcessConvertTrackArray(ctx context.Context, tracks []ConvertTrack) {

	if len(tracks) == 0 {
		return
	}

	p := BuildProgress(len(tracks), 1)

	tracksChan := pipeline.Emit(tracks...)

	convertOut := pipeline.ProcessConcurrently(ctx, 4, pipeline.NewProcessor(func(ctx context.Context, t ConvertTrack) (ConvertTrack, error) {

		if t == (ConvertTrack{}) {
			return t, helpers.ErrConvertTrackEmpty
		}

		e.Logger.Info(fmt.Sprintf("Converting: %s", t.Name))

		t, err := e.convertTrack(t)
		if err != nil {
			return t, err
		}

		e.Logger.Info(fmt.Sprintf("Finished converting: %s", t.Name))

		e.progress(p.Complete(t.ID))

		return t, nil
	}, func(t ConvertTrack, err error) {

		if !strings.Contains(err.Error(), "context canceled") {
			e.Logger.NonFatalError(fault.Wrap(
				err,
				fmsg.With("error processing convert track"),
			))
		}
		e.progress(p.Complete(t.ID))
	}), tracksChan)

	for range convertOut {
		t := <-convertOut
		_ = t
	}
}

func (e *OpEnv) convertTrack(track ConvertTrack) (ConvertTrack, error) {

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
