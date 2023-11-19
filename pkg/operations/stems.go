package operations

import (
	"context"
	"fmt"
	"os"
	"strings"

	b64 "encoding/base64"

	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/deliveryhero/pipeline/v2"
)

type DemucsModels int

const (
	Demucs DemucsModels = iota
	DemucsFT
	Demucs6
	DemucsMMI
	MDX
	MDXExtra
	MDXQ
	MDXQExtra
	SIG
)

func (d DemucsModels) String() string {
	switch d {
	case Demucs:
		return "htdemucs"
	case DemucsFT:
		return "htdemucs_ft"
	case Demucs6:
		return "htdemucs_6s"
	case DemucsMMI:
		return "hdemucs_mmi"
	case MDX:
		return "mdx"
	case MDXExtra:
		return "mdx_extra"
	case MDXQ:
		return "mdx_q"
	case MDXQExtra:
		return "mdx_extra_q"
	case SIG:
		return "SIG"
	default:
		return "htdemucs"
	}
}

/*
getStemPaths gets all of the files in the provided directory which should be converted to stems based on the config

if recursion is true, will also get files in subdirectories
*/
func getStemPaths(cfg helpers.Config, inDirPath string, recursion bool) ([]string, error) {
	stemPaths, err := helpers.GetFilesInDir(inDirPath, recursion)
	if err != nil {
		return nil, err
	}
	var validStemPaths []string
	for _, path := range stemPaths {
		if helpers.IsExtensionInArray(path, cfg.ExtensionsToConvertToStems) {
			validStemPaths = append(validStemPaths, path)
		}
	}
	return validStemPaths, nil
}

func parallelProcessStemTrackArray(ctx context.Context, o OperationProcess, tracks []StemTrack) {

	if len(tracks) == 0 {
		return
	}

	numSteps := 3

	if tracks[0].StemsOnly {
		numSteps = 1
	}

	p := buildProgress(len(tracks), numSteps)

	tracksChan := pipeline.Emit(tracks...)

	demucsOut := pipeline.ProcessConcurrently(ctx, 1, pipeline.NewProcessor(func(ctx context.Context, t StemTrack) (StemTrack, error) {
		if t == (StemTrack{}) {
			return t, helpers.ErrStemTrackEmpty
		}

		if t.SkipDemucs {
			o.StepCallback(stepFinishedStepInfo(fmt.Sprintf("Skipping demucs seperation for: %s", t.Name), p.step(t.ID)))
			return t, nil
		}

		o.StepCallback(stepStartedStepInfo(fmt.Sprintf("Performing demucs separation for: %s", t.Name)))
		t, err := demucsSeparate(t)

		if err != nil {
			return t, err
		}

		o.StepCallback(stepFinishedStepInfo("Finished demucs separation for "+t.Name, p.step(t.ID)))

		return t, nil
	}, func(t StemTrack, err error) {
		if !strings.Contains(err.Error(), "context canceled") {
			o.StepCallback(
				warningStepInfo(helpers.GenErrDemucsSepStep(t.Name, err)),
			)
		}
	}), tracksChan)

	mergeM4aOut := pipeline.ProcessConcurrently(ctx, 2, pipeline.NewProcessor(func(ctx context.Context, t StemTrack) (StemTrack, error) {
		if t == (StemTrack{}) {
			return t, helpers.ErrStemTrackEmpty
		}

		if t.StemsOnly {
			return t, nil
		}

		o.StepCallback(stepStartedStepInfo(fmt.Sprintf("Merging files to Traktor stem file for: %s", t.Name)))

		t, err := mergeToM4a(t)

		if err != nil {
			return t, err
		}

		o.StepCallback(stepFinishedStepInfo("Finished merging files for: "+t.Name, p.step(t.ID)))

		return t, nil
	}, func(t StemTrack, err error) {
		if !strings.Contains(err.Error(), "context canceled") {
			o.StepCallback(
				warningStepInfo(helpers.GenErrMergeM4AStep(t.Name, err)),
			)
		}
	}), demucsOut)

	addMetadataOut := pipeline.ProcessConcurrently(ctx, 4, pipeline.NewProcessor(func(ctx context.Context, t StemTrack) (StemTrack, error) {
		if t == (StemTrack{}) {
			return t, helpers.ErrStemTrackEmpty
		}

		if t.StemsOnly {
			return t, nil
		}

		o.StepCallback(stepStartedStepInfo(fmt.Sprintf("Adding metadata to Traktor stem file for: %s", t.Name)))
		t, err := addMetadata(t)

		if err != nil {
			return t, err
		}

		o.StepCallback(stepFinishedStepInfo("Finished adding metadata for: "+t.Name, p.step(t.ID)))

		return t, nil
	}, func(t StemTrack, err error) {
		if !strings.Contains(err.Error(), "context canceled") {
			o.StepCallback(
				warningStepInfo(helpers.GenErrAddMetadataStep(t.Name, err)),
			)
		}
	}), mergeM4aOut)

	cleanupOut := pipeline.ProcessConcurrently(ctx, 4, pipeline.NewProcessor(func(ctx context.Context, t StemTrack) (StemTrack, error) {
		if t == (StemTrack{}) {
			return t, helpers.ErrStemTrackEmpty
		}

		t, err := cleanUp(t)

		if err != nil {
			return t, err
		}

		return t, nil

	}, func(t StemTrack, err error) {
		if !strings.Contains(err.Error(), "context canceled") {
			o.StepCallback(
				warningStepInfo(helpers.GenErrCleanupStep(t.Name, err)),
			)
		}
	}), addMetadataOut)

	for range cleanupOut {
		t := <-cleanupOut
		o.StepCallback(trackFinishedStepInfo(fmt.Sprintf("Finished processing: %s", t.Name), p.complete(t.ID)))
	}
}

/*
demucsSeparate calls demucs to split a file into stem tracks
*/
func demucsSeparate(track StemTrack) (StemTrack, error) {

	// create stem dir if it doesn't exist
	os.MkdirAll(track.StemDir, os.ModePerm)

	demucsArgs := []string{
		"demucs",
		"--out", track.StemDir,
		"--filename", fmt.Sprintf("%s{stem}.{ext}", track.StemDir),
		"--jobs", "4",
		"--name", "htdemucs",
		"-d", "cuda",
		track.OriginalFile.FileInfo.FullPath,
	}

	if track.OriginalFile.FileInfo.FileExtension == ".mp3" {
		demucsArgs = append(demucsArgs, "--mp3")
	}

	// run demucs
	_, err := helpers.CmdExec(
		demucsArgs...,
	)

	if err != nil {
		return track, err
	}

	return track, nil
}

func mergeToM4a(track StemTrack) (StemTrack, error) {

	// create output file dir if it doesn't exist
	os.MkdirAll(track.OutFile.FileInfo.DirPath, os.ModePerm)

	// convert stems to m4a
	_, err := helpers.CmdExec(
		"ffmpeg",
		"-i", track.OriginalFile.FileInfo.FullPath,
		"-i", track.DrumsFile.AudioFile.FileInfo.FullPath,
		"-i", track.BassFile.AudioFile.FileInfo.FullPath,
		"-i", track.OtherFile.AudioFile.FileInfo.FullPath,
		"-i", track.VocalsFile.AudioFile.FileInfo.FullPath,
		"-map", "0", "-map", "1", "-map", "2", "-map", "3", "-map", "4",
		track.OutFile.FileInfo.FullPath,
	)

	if err != nil {
		return track, err
	}

	return track, nil
}

func addMetadata(track StemTrack) (StemTrack, error) {

	// add metadata to m4a
	_, err := helpers.CmdExec(
		"MP4Box",
		track.OutFile.FileInfo.FullPath,
		"-udta", fmt.Sprintf("0:type=stem:src=base64,%s", getTraktorMetadata()),
	)
	if err != nil {
		return track, err
	}

	return track, nil
}

func cleanUp(track StemTrack) (StemTrack, error) {
	// deletes stem files/ dirs

	if !track.StemsOnly {
		os.RemoveAll(track.StemDir)
	}

	return track, nil
}

func getTraktorMetadata() string {
	drumColour := "#009E73"
	bassColour := "#D55E00"
	otherColour := "#CC79A7"
	vocalColour := "#56B4E9"

	dataString := fmt.Sprintf(
		`{
			"mastering_dsp": {
				"compressor": {
					"ratio": 3, 
					"output_gain": 0.5, 
					"enabled": false, 
					"release": 0.300000011920929, 
					"attack": 0.003000000026077032, 
					"input_gain": 0.5, 
					"threshold": 0, 
					"hp_cutoff": 300, 
					"dry_wet": 50
				}, 
				"limiter": {
					"release": 0.05000000074505806, 
					"threshold": 0, 
					"ceiling": -0.3499999940395355, 
					"enabled": false
				}
			}, 
			"version": 1, 
			"stems": [
				{"color": "%s", "name": "Drums"}, 
				{"color": "%s", "name": "Bass"}, 
				{"color": "%s", "name": "Other"}, 
				{"color": "%s", "name": "Vocals"}
			]
		}`,
		drumColour,
		bassColour,
		otherColour,
		vocalColour,
	)
	sEnc := b64.StdEncoding.EncodeToString([]byte(dataString))
	return sEnc
}
