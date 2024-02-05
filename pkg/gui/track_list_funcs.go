package gui

import (
	"context"
	"fmt"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/operations"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
Provides functions used to generate the functions used to interface with the track_list of a given playlist

These functions are called to generate functions that are passed through to widget constructors
*/

/*
getDownloadSoundCloudTrackFunc returns a function that can be used to download a SoundCloud track, we call
this function when opening a playlist popup.

Generating the function here saves us passing lots of data down the track widgets (i.e. env/ playlist name)
May consider changing this and attaching these to the widgets in the future
*/
func (e *guiEnv) getDownloadSoundCloudTrackFunc(selectedTrack *iwidget.SelectedTrackBinding, playlistName string) func() {
	return func() {

		selectedTrack.LockSelected()
		defer selectedTrack.UnlockSelected()

		track := selectedTrack.TrackBinding.Track

		opEnv := e.opEnv()
		opEnv.BuildOperationHandler(
			func(i float64) {},
			func(data map[string]any) {
				filePath, ok := data["filepath"].(string)
				if !ok {
					e.showErrorDialog(fault.Wrap(
						fault.New("error casting filepath to string"),
						fmsg.WithDesc(
							"error parsing filepath from operation data",
							"Error parsing track download results",
						),
					), true)
					return
				}
				track.LocalPath = filePath
				e.showInfoDialog(
					"Download Successful",
					fmt.Sprintf("Downloaded %s to %s", track.Name, track.LocalPath),
				)
			},
			func(err error) {
				e.showErrorDialog(fault.Wrap(
					err,
					fmsg.WithDesc(
						"error downloading soundcloud track",
						"Error downloading SoundCloud track",
					),
				), true)
			},
		)

		opEnv.DownloadSoundCloudFile(*track, playlistName)
	}
}

func (e *guiEnv) getSaveSoundCloudTrackFunc(selectedTrack *iwidget.SelectedTrackBinding) func() {

	return func() {

		track := selectedTrack.TrackBinding.Track

		ctx := context.Background()
		ctx = fctx.WithMeta(ctx,
			"track_name", track.Name,
			"track_permalink", track.Permalink,
			"track_external_id", fmt.Sprintf("%d", track.ExternalID),
			"track_local_path", track.LocalPath,
		)

		err := e.SerenDB.TxUpsertSoundCloudTracks([]data.SoundcloudTrack{track.ToDB()})
		if err != nil {
			e.showErrorDialog(
				fault.Wrap(
					err,
					fctx.With(ctx),
					fmsg.With("error saving track to database"),
				),
				true,
			)
			return
		}
	}
}

/*
getRefreshSoundCloudPlaylistFunc returns a function that can be used to refresh a SoundCloud playlist

Generating the function here saves us passing lots of data down the track widgets (i.e. env/ playlist name)
*/
func (e *guiEnv) getRefreshSoundCloudPlaylistFunc(playlist streaming.SoundCloudPlaylist, trackListBinding *iwidget.TrackListBinding) func() {

	ctx := context.Background()

	ctx = fctx.WithMeta(ctx,
		"playlist_name", playlist.Name,
		"playlist_external_id", fmt.Sprintf("%d", playlist.ExternalID),
	)

	processResultsFunc := func(p streaming.SoundCloudPlaylist, err error) {

		if err != nil {
			e.showErrorDialog(fault.Wrap(
				err,
				fctx.With(ctx),
				fmsg.WithDesc(
					"err refreshing SoundCloud playlist",
					"Error refreshing SoundCloud playlist",
				),
			), true)
			return
		}

		currentTracksMap := make(map[int64]streaming.SoundCloudTrack)
		existingTracksMap := make(map[int64]streaming.SoundCloudTrack)

		for _, t := range p.Tracks {
			currentTracksMap[t.ExternalID] = t
		}

		var tracksToSave []streaming.SoundCloudTrack

		/*
			store all tracks prior to updating inside the existingTracksMap,
				we can use this to figure out which tracks are new

			we also set the RemovedFromPlaylist flag to true if a track is no longer in the playlist,
				and then add to the tracksToSave slice to be updated in the database
		*/
		for _, t := range trackListBinding.Tracks {
			existingTracksMap[t.ExternalID] = *t
			v, ok := currentTracksMap[t.ExternalID]
			if !ok {
				t.RemovedFromPlaylist = true
				tracksToSave = append(tracksToSave, *t)
			}
			*t = v
		}

		/*
			iterate through the tracks in the playlist,
				if the track does not exist in the existingTracksMap,
				it is new, so we add it to the trackListBinding list (to be displayed in the UI)
				we do not save it to the database, as this was already done inside operations.go
					(though i may want to change this, but it will require changing the refresh playlist func too)
		*/
		for _, t := range p.Tracks {
			_, ok := existingTracksMap[t.ExternalID]
			if !ok {
				// track does not exist in current list, so we add it & save it
				t.Playlists = append(t.Playlists, playlist)
				newT := t
				trackListBinding.Tracks = append(trackListBinding.Tracks, &newT)
				tracksToSave = append(tracksToSave, newT)
			}
		}

		if len(tracksToSave) > 0 {

			dataT := make([]data.SoundcloudTrack, len(tracksToSave))
			for i, t := range tracksToSave {
				dataT[i] = t.ToDB()
			}

			err = e.SerenDB.TxUpsertSoundCloudTracks(dataT)

			if err != nil {
				e.showErrorDialog(fault.Wrap(
					err,
					fctx.With(ctx),
					fmsg.WithDesc(
						"error saving tracks to database",
						"Error saving tracks to database",
					),
				), true)
				return
			}
		}

		trackListBinding.ApplyFilterSort()
	}

	return func() {

		opEnv := e.opEnv()
		opEnv.BuildOperationHandler(
			func(i float64) {},
			func(_ map[string]any) {},
			func(err error) {},
		)

		opts := operations.GetSoundCloudPlaylistOpts{
			PlaylistURL: playlist.Permalink,
			Refresh:     true,
		}

		opEnv.GetSoundCloudPlaylist(ctx, opts, processResultsFunc)
	}
}
