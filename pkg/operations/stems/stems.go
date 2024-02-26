package operations

import (
	"context"
	"fmt"
	"os"
	"strings"

	b64 "encoding/base64"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
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

func (e *StemEnv) ConvertStemTracks(ctx context.Context, tracks []StemTrack) {

	numSteps := 3

	if tracks[0].StemsOnly {
		numSteps = 1
	}

	e.BuildProgressTracker(len(tracks), numSteps)

	tracksChan := pipeline.Emit(tracks...)

	demucsOut := pipeline.ProcessConcurrently(ctx, 1, pipeline.NewProcessor(func(ctx context.Context, t StemTrack) (StemTrack, error) {
		if t == (StemTrack{}) {
			return t, helpers.ErrStemTrackEmpty
		}

		if t.SkipDemucs {
			e.Logger.Info(fmt.Sprintf("Skipping demucs separation for: %s", t.Name))
			return t, nil
		}

		e.Logger.Info(fmt.Sprintf("Performing demucs separation for: %s", t.Name))
		t, err := e.demucsSeparate(t)

		if err != nil {
			return t, err
		}

		e.Logger.Info(fmt.Sprintf("Finished demucs separation for: %s", t.Name))
		e.ProcessStep(t.ID)

		return t, nil
	}, func(t StemTrack, err error) {

		if !strings.Contains(err.Error(), "context canceled") {
			e.Logger.NonFatalError(fault.Wrap(
				err,
				fctx.With(fctx.WithMeta(
					ctx,
					"name", t.Name,
				)),
				fmsg.WithDesc(
					"demucs separation error",
					"There was an error calling demucs to separate the stems",
				),
			))
		}

		e.ProcessComplete(t.ID)
	}), tracksChan)

	mergeM4aOut := pipeline.ProcessConcurrently(ctx, 2, pipeline.NewProcessor(func(ctx context.Context, t StemTrack) (StemTrack, error) {
		if t == (StemTrack{}) {
			return t, helpers.ErrStemTrackEmpty
		}

		if t.StemsOnly {
			return t, nil
		}

		e.Logger.Info(fmt.Sprintf("Merging files to Traktor stem file for: %s", t.Name))

		t, err := e.mergeToM4a(t)

		if err != nil {
			return t, err
		}

		e.Logger.Info(fmt.Sprintf("Finished merging files for: %s", t.Name))
		e.ProcessStep(t.ID)

		return t, nil
	}, func(t StemTrack, err error) {

		if !strings.Contains(err.Error(), "context canceled") {
			e.Logger.NonFatalError(fault.Wrap(
				err,
				fctx.With(fctx.WithMeta(
					ctx,
					"name", t.Name,
				)),
				fmsg.WithDesc(
					"error merging files",
					"There was an error merging the stems into a single file",
				),
			))
		}

		e.ProcessComplete(t.ID)
	}), demucsOut)

	cleanupOut := pipeline.ProcessConcurrently(ctx, 4, pipeline.NewProcessor(func(ctx context.Context, t StemTrack) (StemTrack, error) {
		if t == (StemTrack{}) {
			return t, helpers.ErrStemTrackEmpty
		}

		t, err := e.cleanUp(t)

		if err != nil {
			return t, err
		}

		return t, nil

	}, func(t StemTrack, err error) {
		if !strings.Contains(err.Error(), "context canceled") {
			e.Logger.NonFatalError(fault.Wrap(
				err,
				fctx.With(fctx.WithMeta(
					ctx,
					"name", t.Name,
				)),
				fmsg.WithDesc(
					"error cleaning up",
					"There was an error cleaning up the stem files",
				),
			))
		}
	}), mergeM4aOut)

	for range cleanupOut {
		t := <-cleanupOut
		e.Logger.Info(fmt.Sprintf("Finished processing: %s", t.Name))
	}
}

/*
demucsSeparate calls demucs to split a file into stem tracks
*/
func (e *StemEnv) demucsSeparate(track StemTrack) (StemTrack, error) {

	// create stem dir if it doesn't exist
	os.MkdirAll(track.StemDir, os.ModePerm)

	demucsArgs := []string{
		"demucs",
		"--out", track.StemDir,
		"--filename", fmt.Sprintf("%s{stem}.{ext}", track.StemDir),
		"--jobs", "4",
		"--name", "htdemucs",
	}

	if e.Config.CudaEnabled {
		demucsArgs = append(demucsArgs, "-d", "cuda")
	}

	if track.OriginalFile.FileInfo.FileExtension == ".mp3" {
		demucsArgs = append(demucsArgs, "--mp3")
	}

	demucsArgs = append(demucsArgs, track.OriginalFile.FileInfo.FullPath)

	// run demucs
	out, err := helpers.CmdExec(
		demucsArgs...,
	)

	if err != nil {
		e.Logger.Debug(out)
		return track, err
	}

	return track, nil
}

func (e *StemEnv) mergeToM4a(track StemTrack) (StemTrack, error) {

	// create output file dir if it doesn't exist
	os.MkdirAll(track.OutFile.FileInfo.DirPath, os.ModePerm)

	// convert stems to m4a
	out, err := helpers.CmdExec(
		"ffmpeg",
		"-i", track.OriginalFile.FileInfo.FullPath,
		"-i", track.DrumsFile.AudioFile.FileInfo.FullPath,
		"-i", track.BassFile.AudioFile.FileInfo.FullPath,
		"-i", track.OtherFile.AudioFile.FileInfo.FullPath,
		"-i", track.VocalsFile.AudioFile.FileInfo.FullPath,
		"-map", "0", "-map", "1", "-map", "2", "-map", "3", "-map", "4",
		"-metadata", fmt.Sprintf("udta:0:type=stem:src=base64,%s", e.getTraktorMetadata()),
		"-vn",
		track.OutFile.FileInfo.FullPath,
	)

	if err != nil {
		e.Logger.Debug(out)
		return track, err
	}

	return track, nil
}

func (e *StemEnv) cleanUp(track StemTrack) (StemTrack, error) {
	// deletes stem files/ dirs

	if !track.StemsOnly {
		os.RemoveAll(track.StemDir)
	}

	return track, nil
}

func (e *StemEnv) getTraktorMetadata() string {
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
