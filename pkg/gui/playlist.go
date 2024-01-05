package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/gui/uihelpers"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
openPlaylistPopup opens a popup window for a given playlist

This is called when a user clicks the 'open playlist' button on a playlist
*/
func (e *guiEnv) openPlaylistPopup(playlist streaming.SoundCloudPlaylist) {

	loading := newViewLoading(fmt.Sprintf("Loading tracks for %s...", playlist.Name))

	var trackListBinding iwidget.TrackListBinding
	selectedTrack := &iwidget.SelectedTrackBinding{}

	// Build the track list widget (displays list of tracks to select)
	trackListSection := iwidget.NewTrackListSection(
		e.getWidgetBase(),
		&trackListBinding,
		selectedTrack,
		iwidget.TrackListFuncs{
			RefreshSoundCloudPlaylist: e.getRefreshSoundCloudPlaylistFunc(playlist, &trackListBinding),
		},
	)

	// Build the track section widget (displays info about selected track)
	trackSection := iwidget.NewTrackSection(
		e.getWidgetBase(),
		iwidget.TrackFuncs{
			DownloadSoundCloudTrack: e.getDownloadSoundCloudTrackFunc(selectedTrack, playlist.Name),
			SaveSoundCloudTrackToDB: e.getSaveSoundCloudTrackFunc(selectedTrack),
			OnError:                 e.displayErrorDialog,
		},
	)
	trackSection.Bind(selectedTrack)

	splitContainer := container.NewHSplit(
		trackListSection,
		trackSection,
	)
	splitContainer.SetOffset(0)

	content := container.NewStack(
		splitContainer,
		loading,
	)

	playlistPopup := uihelpers.NewPercentagePopup(
		playlist.Name,
		content,
		e.mainWindow,
		e.resizeEvents,
		0.9, 0.9,
		fyne.NewSize(1280, 0),
		func(close func()) {
			e.guiState.busy = false
		},
	)

	// Load the tracks in the background, hide the loading screen when done
	// This should be quick as it only requires a database query
	go func(tlb *iwidget.TrackListBinding) {
		err := e.loadSoundCloudPlaylistTracks(playlist.ExternalID, tlb)
		if err != nil {
			playlistPopup.Hide()
			e.displayErrorDialog(err)
			return
		}
		loading.Hide()
	}(&trackListBinding)

	e.guiState.busy = true

	playlistPopup.Show()
}
