package gui

import (
	"context"
	"database/sql"
	"net/url"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/operations"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
Provides functions used to interface with a playlist, these functions are called to generate functions that are passed
through to widget constructors
*/

/*
getAddPlaylistCallback returns a function that can be used to add a playlist.

This function calls SoundCloud, and adds the playlist to the database, it is attached to the
'add playlist' button
*/
func (e *guiEnv) getAddPlaylistCallback(playlistBindVals *iwidget.PlaylistBindingList, onAdd func()) func(string) {
	ctx := context.Background()
	ctx, ctxClose := context.WithCancel(ctx)
	opEnv := e.opEnv()
	opEnv.RegisterOperationHandler(
		func(i operations.OperationProgressInfo) {},
		func(i operations.OperationFinishedInfo) {
			ctxClose()
		},
	)

	return func(urlRaw string) {

		ctx = fctx.WithMeta(
			ctx,
			"url_raw", urlRaw,
		)

		pbi := iwidget.PlaylistBindingItem{
			Base: e.getWidgetBase(),
		}

		netUrl, err := url.Parse(urlRaw)

		if err != nil {
			err = fault.Wrap(
				err,
				fctx.With(ctx),
				fmsg.WithDesc(
					"error parsing url",
					"Could not parse the URL, please check it is correct",
				),
			)

			e.logger.NonFatalError(err)
			pbi.SetFailed(err)

			playlistBindVals.Append(&pbi)
			onAdd()
			return
		}

		netUrl.RawQuery = ""

		ctx = fctx.WithMeta(
			ctx,
			"parsed_url", netUrl.String(),
		)

		pbi.SetFinding(
			streaming.SoundCloudPlaylist{
				SearchUrl: netUrl.String(),
			},
		)
		playlistBindVals.Append(&pbi)
		onAdd()

		go opEnv.GetSoundCloudPlaylist(ctx, operations.GetSoundCloudPlaylistOpts{
			PlaylistURL: netUrl.String(),
		}, func(p streaming.SoundCloudPlaylist, err error) {
			if err != nil {
				err = fault.Wrap(
					err,
					fctx.With(ctx),
					fmsg.WithDesc(
						"error getting playlist",
						"There was an error getting the playlist from SoundCloud",
					),
				)

				pbi.SetFailed(err)
				e.logger.NonFatalError(err)

				return
			}

			pbi.SetFound(p)
			onAdd()
		})
	}
}

func (e *guiEnv) loadSoundCloudPlaylists(playlistBindingList *iwidget.PlaylistBindingList) error {
	playlists, err := e.SerenDB.ListSoundCloudPlaylists(context.Background())

	if err != nil {
		return fault.Wrap(
			err,
			fmsg.WithDesc(
				"error getting soundcloud playlists from db",
				"There was an error getting the SoundCloud playlists from the database",
			),
		)
	}

	e.logger.Debugf("successfully got %v playlists from db, loading into gui", len(playlists))

	playlistBindingList.Load(playlists)

	return nil
}

func (e *guiEnv) loadSoundCloudPlaylistTracks(playlistExtID int64, trackListBinding *iwidget.TrackListBinding) error {
	tracks, err := e.SerenDB.ListSoundCloudTracksByPlaylistExternalID(
		context.Background(),
		sql.NullInt64{Valid: true, Int64: playlistExtID},
	)
	if err != nil {
		return fault.Wrap(
			err,
			fmsg.WithDesc(
				"error getting soundcloud tracks from db",
				"There was an error getting the SoundCloud tracks from the database for this playlist",
			),
		)
	}
	streamTracks := make([]*streaming.SoundCloudTrack, len(tracks))
	for i, track := range tracks {
		streamTrack := streaming.SoundCloudTrack{}
		streamTrack.LoadFromDB(track)
		streamTracks[i] = &streamTrack
	}
	trackListBinding.Set(streamTracks)
	trackListBinding.ApplyFilterSort()

	return nil
}
