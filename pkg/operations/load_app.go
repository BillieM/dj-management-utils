package operations

import (
	"context"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
TODO: handle all these with step handlers
*/

func (e *OpEnv) CheckLocalPaths() {

	var brokenPathChanged []streaming.SoundCloudTrack

	tracks, err := e.SerenDB.ListSoundCloudTracksHasLocalPath(context.Background())
	if err != nil {
		e.Logger.Error(fault.Flatten(fault.Wrap(
			err,
			fmsg.With("error listing tracks with local path in db"),
		)))
		return
	}

	for _, track := range tracks {

		t := streaming.SoundCloudTrack{}
		t.LoadFromDB(track)

		fileExists := helpers.DoesFileExist(t.LocalPath)

		if fileExists == t.LocalPathBroken {

			// if no local path and not broken, skip it
			// as it's not a broken path
			if t.LocalPath == "" && !t.LocalPathBroken {
				continue
			}

			t.LocalPathBroken = !fileExists

			brokenPathChanged = append(brokenPathChanged, t)
		}
	}

	if len(brokenPathChanged) == 0 {
		e.Logger.Debug("no tracks with changed broken status")
		return
	}

	e.Logger.Debugf(fmt.Sprintf(
		"found %d tracks with changed broken status\n",
		len(brokenPathChanged),
	))

	dataT := make([]data.SoundcloudTrack, len(brokenPathChanged))

	for i, t := range brokenPathChanged {
		dataT[i] = t.ToDB()
	}

	err = e.SerenDB.TxUpsertSoundCloudTracks(dataT)

	if err != nil {
		e.Logger.NonFatalError(fault.Wrap(
			err,
			fmsg.With(
				"error updating tracks in db",
			),
		))
		return
	}
}

func (e *OpEnv) IndexCollections() {

}

func (e *OpEnv) IndexLocalFolders() {

}
