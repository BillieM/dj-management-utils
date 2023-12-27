package gui

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/gui/uihelpers"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*

 */

type playlistState int

const (
	NotSet playlistState = iota
	Found
	Finding
	Failed
)

func (p playlistState) String() string {
	switch p {
	case NotSet:
		return "NotSet"
	case Found:
		return "Found"
	case Finding:
		return "Downloading"
	case Failed:
		return "Failed"
	default:
		return "Unknown"
	}
}

/*
playlistBindingList stores a list of playlistBindingItem structs

It is used to display a list of playlists as playlistWidgets in the UI
*/
type playlistBindingList struct {
	bindBase

	Items []*playlistBindingItem
}

func (i *playlistBindingList) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *playlistBindingList) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (i *playlistBindingList) GetItem(index int) (binding.DataItem, error) {
	i.Lock()
	defer i.Unlock()
	if index < 0 || index >= len(i.Items) {
		return nil, helpers.ErrIndexOutOfBounds
	}

	return i.Items[index], nil
}

func (i *playlistBindingList) Length() int {
	i.Lock()
	defer i.Unlock()
	return len(i.Items)
}

func (i *playlistBindingList) Append(p *playlistBindingItem) {
	i.Lock()
	defer i.Unlock()
	i.Items = append(i.Items, p)
}

/*
load loads all playlists from the database into the playlistBindingList
*/
func (i *playlistBindingList) load(s *data.SerenDB) {

	// TODO err handling...
	playlists, _ := s.ListSoundCloudPlaylists(context.Background())

	for _, playlist := range playlists {

		p := streaming.SoundCloudPlaylist{}
		p.LoadFromDB(playlist, nil)

		i.Append(&playlistBindingItem{
			playlist: p,
			state:    Found,
		})
	}
}

/*
playlistBindingItem is a struct that contains the data for a playlist

It is used to display a playlist as a playlistWidget in the UI
*/
type playlistBindingItem struct {
	bindBase

	// may want a context in here ?? later problem...
	playlist streaming.SoundCloudPlaylist
	state    playlistState
	err      error
}

func (i *playlistBindingItem) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *playlistBindingItem) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

/*
playlistWidget displays a playlist in the ui
*/
type playlistWidget struct {
	widget.BaseWidget

	findingContent fyne.CanvasObject
	foundContent   fyne.CanvasObject
	failedContent  fyne.CanvasObject

	searchUrl *widget.Hyperlink

	err *widget.Label

	progress *widget.ProgressBarInfinite

	name            *widget.Label
	openPlaylistBtn *widget.Button

	ctxCancel func() // used to cancel a downloading context
}

/*
newPlaylistWidget returns a new instance of playlistWidget
*/
func newPlaylistWidget() *playlistWidget {
	i := &playlistWidget{
		name:            widget.NewLabel(""),
		openPlaylistBtn: widget.NewButton("Open Playlist", func() {}),
		searchUrl:       widget.NewHyperlink("", nil),
		err:             widget.NewLabel(""),
		progress:        widget.NewProgressBarInfinite(),
	}

	i.findingContent = i.progress
	i.failedContent = i.err
	i.foundContent = container.NewBorder(
		i.name,
		nil,
		nil,
		nil,
		i.openPlaylistBtn,
	)

	i.err.Importance = widget.WarningImportance

	// i.findingContent.Hide()
	// i.failedContent.Hide()
	// i.foundContent.Hide()

	i.ExtendBaseWidget(i)

	return i
}

func (i *playlistWidget) CreateRenderer() fyne.WidgetRenderer {

	c := container.NewBorder(
		i.searchUrl, nil, nil, nil,
		container.NewStack(
			i.findingContent,
			i.foundContent,
			i.failedContent,
		),
	)

	return widget.NewSimpleRenderer(c)
}

/*
addPlaylistWidget displays a section used for adding a playlist to the ui
*/
type addPlaylistWidget struct {
	widget.BaseWidget

	urlEntry        *widget.Entry
	submitButton    *widget.Button
	validationLabel *widget.Label
}

/*
newAddPlaylistWidget returns a new instance of addPlaylistWidget
*/
func newAddPlaylistWidget(addPlaylist func(string)) *addPlaylistWidget {

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("SoundCloud Playlist URL")
	urlEntry.Validator = validation.NewRegexp(`soundcloud\.com\/.*\/sets`, "not a valid SoundCloud playlist url")

	urlEntry.OnSubmitted = func(s string) {
		if s == "" {
			return
		}
		err := urlEntry.Validate()
		if err != nil {
			urlEntry.SetValidationError(err)
		} else {
			urlEntry.SetText("")
			addPlaylist(s)
		}
		urlEntry.Refresh()
	}

	submitBtn := widget.NewButton("Add Playlist", func() {
		s := urlEntry.Text
		if s == "" {
			return
		}
		err := urlEntry.Validate()
		if err != nil {
			urlEntry.SetValidationError(err)
		} else {
			urlEntry.SetText("")
			addPlaylist(s)
		}
		urlEntry.Refresh()
	})

	validationLabel := widget.NewLabel("")

	urlEntry.SetOnValidationChanged(func(err error) {
		if err != nil {
			validationLabel.SetText(err.Error())
		} else {
			validationLabel.SetText("")
		}
	})

	i := &addPlaylistWidget{
		submitButton:    submitBtn,
		urlEntry:        urlEntry,
		validationLabel: validationLabel,
	}

	widget.NewForm()

	i.ExtendBaseWidget(i)

	return i
}

func (i *addPlaylistWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewBorder(
		widget.NewLabel("Add playlist"),
		nil, nil, nil,
		container.NewBorder(
			nil, i.validationLabel, nil, i.submitButton, i.urlEntry,
		),
	)
	return widget.NewSimpleRenderer(c)
}

/*
getAddPlaylistCallback returns a function that can be used to add a playlist.

This function calls SoundCloud, and adds the playlist to the database, it is attached to the
'add playlist' button
*/
func (e *guiEnv) getAddPlaylistCallback(playlistBindVals *playlistBindingList, refreshFunc func()) func(string) {
	ctx := context.Background()
	ctx, ctxClose := context.WithCancel(ctx)
	opEnv := e.opEnv()
	opEnv.RegisterStepHandler(streamingStepHandler{
		stepFunc:     func() {},
		finishedFunc: func() { ctxClose() },
	})

	return func(urlRaw string) {

		pbi := &playlistBindingItem{
			err:   nil,
			state: Finding,
		}

		netUrl, err := url.Parse(urlRaw)

		if err != nil {
			pbi.state = Failed
			pbi.err = err
			pbi.playlist = streaming.SoundCloudPlaylist{SearchUrl: urlRaw}
			playlistBindVals.Append(pbi)
			refreshFunc()
			opEnv.Logger.Error(fault.Flatten(fault.Wrap(
				err,
				fmsg.With("error parsing url"),
			)))
			return
		}

		netUrl.RawQuery = ""

		pbi.playlist = streaming.SoundCloudPlaylist{
			SearchUrl: netUrl.String(),
		}

		playlistBindVals.Append(pbi)
		refreshFunc()

		go opEnv.GetSoundCloudPlaylist(ctx, operations.GetSoundCloudPlaylistOpts{
			PlaylistURL: netUrl.String(),
		}, func(p streaming.SoundCloudPlaylist, err error) {
			if err != nil {
				pbi.state = Failed
				pbi.err = err
				opEnv.Logger.Error(
					"error getting playlist",
					fault.Flatten(err),
				)
				return
			}

			pbi.playlist = p
			pbi.state = Found
			pbi.err = nil

			refreshFunc()
		})
	}
}

/*
updatePlaylistsList updates the playlistWidget for a given playlist

This is called when a playlist is added, or when the state of the playlist changes,
for example when the request to get information about the playlist fails, or is completed
*/
func (e *guiEnv) updatePlaylistsList(playlistWidget *playlistWidget, playlistBindingItem *playlistBindingItem) {

	playlist := playlistBindingItem.playlist

	playlistWidget.searchUrl.SetText(playlist.SearchUrl)
	playlistWidget.searchUrl.SetURLFromString(playlist.SearchUrl)

	switch playlistBindingItem.state {
	case Finding:
		playlistWidget.findingContent.Show()
		playlistWidget.foundContent.Hide()
		playlistWidget.failedContent.Hide()
	case Failed:
		playlistWidget.findingContent.Hide()
		playlistWidget.foundContent.Hide()
		playlistWidget.failedContent.Show()
		playlistWidget.err.SetText(playlistBindingItem.err.Error())
	case Found:
		playlistWidget.findingContent.Hide()
		playlistWidget.foundContent.Show()
		playlistWidget.failedContent.Hide()
		playlistWidget.name.SetText(playlist.Name)
		playlistWidget.openPlaylistBtn.OnTapped = func() {
			if e.guiState.busy {
				e.showErrorDialog(helpers.ErrBusyPleaseFinishFirst)
				return
			}
			e.openPlaylistPopup(playlist)
		}
	}

	playlistWidget.Refresh()
}

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
		e.mainWindow,
		&trackListBinding,
		selectedTrack,
		iwidget.TrackListFuncs{
			RefreshSoundCloudPlaylist: e.getRefreshSoundCloudPlaylistFunc(playlist, &trackListBinding),
		},
	)

	// Build the track section widget (displays info about selected track)
	trackSection := iwidget.NewTrackSection(
		e.mainWindow,
		iwidget.TrackFuncs{
			DownloadSoundCloudTrack: e.getDownloadSoundCloudTrackFunc(selectedTrack, playlist.Name),
			SaveSoundCloudTrackToDB: e.getSaveSoundCloudTrackFunc(selectedTrack),
			OnError:                 e.displayErrorDialog,
		},
		e.resizeEvents,
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

	go func(tlb *iwidget.TrackListBinding) {
		tracks, err := e.SerenDB.ListSoundCloudTracksByPlaylistExternalID(
			context.Background(),
			sql.NullInt64{Valid: true, Int64: playlist.ExternalID},
		)
		if err != nil {
			playlistPopup.Hide()
			e.showErrorDialog(err)
			return
		}
		streamTracks := make([]*streaming.SoundCloudTrack, len(tracks))
		for i, track := range tracks {
			streamTrack := streaming.SoundCloudTrack{}
			streamTrack.LoadFromDB(track)
			streamTracks[i] = &streamTrack
		}
		tlb.Set(streamTracks)
		tlb.ApplyFilterSort()
		loading.Hide()
	}(&trackListBinding)

	e.guiState.busy = true

	playlistPopup.Show()
}
