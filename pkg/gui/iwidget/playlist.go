package iwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
Contains widgets related to playlists, these widgets are used as part of the 'playlist matching' functionality
in order to display an overview of a playlist, as well as providing a way of adding new playlists to the application
*/

/*
PlaylistState is an enum that represents the state of a playlist
*/

type PlaylistState int

const (
	NotSet PlaylistState = iota
	Found
	Finding
	Failed
)

func (p PlaylistState) String() string {
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
PlaylistBindingList stores a list of PlaylistBindingItem structs

It is used to display a list of playlists as playlistWidgets in the UI
*/
type PlaylistBindingList struct {
	*Base
	bindBase

	Items []*PlaylistBindingItem
}

func (i *PlaylistBindingList) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *PlaylistBindingList) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (i *PlaylistBindingList) GetItem(index int) (binding.DataItem, error) {
	i.Lock()
	defer i.Unlock()
	if index < 0 || index >= len(i.Items) {
		return nil, helpers.ErrIndexOutOfBounds
	}

	return i.Items[index], nil
}

func (i *PlaylistBindingList) Length() int {
	i.Lock()
	defer i.Unlock()
	return len(i.Items)
}

func (i *PlaylistBindingList) Append(p *PlaylistBindingItem) {
	i.Lock()
	defer i.Unlock()
	i.Items = append(i.Items, p)
}

/*
Load loads all playlists from the database into the PlaylistBindingList
*/
func (i *PlaylistBindingList) Load(playlists []data.SoundcloudPlaylist) {

	for _, playlist := range playlists {

		p := streaming.SoundCloudPlaylist{}
		p.LoadFromDB(playlist, nil)

		pbi := &PlaylistBindingItem{Base: i.Base}
		pbi.SetFound(p)

		i.Append(pbi)
	}
}

/*
PlaylistBindingItem is a struct that contains the data for a playlist

It is used to display a playlist as a playlistWidget in the UI
*/
type PlaylistBindingItem struct {
	*Base
	bindBase

	// may want a context in here ?? later problem...
	playlist streaming.SoundCloudPlaylist
	state    PlaylistState
	err      error
}

func (i *PlaylistBindingItem) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *PlaylistBindingItem) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (i *PlaylistBindingItem) SetFinding(playlist streaming.SoundCloudPlaylist) {
	i.Lock()
	defer i.Unlock()

	i.playlist = playlist
	i.state = Finding

	i.Logger.Debugf("set playlist to finding: %s", i.playlist)
}

func (i *PlaylistBindingItem) SetFound(playlist streaming.SoundCloudPlaylist) {
	i.Lock()
	defer i.Unlock()

	i.playlist = playlist
	i.state = Found

	i.Logger.Debugf("set playlist to found: %s", i.playlist)
}

func (i *PlaylistBindingItem) SetFailed(err error) {
	i.Lock()
	defer i.Unlock()

	i.err = err
	i.state = Failed

	i.Logger.Debugf("set playlist to failed: %s, with err %s", i.playlist, err)
}

/*
playlistWidget displays a playlist in the ui
*/
type Playlist struct {
	widget.BaseWidget

	findingContent fyne.CanvasObject
	foundContent   fyne.CanvasObject
	failedContent  fyne.CanvasObject

	searchUrl *widget.Hyperlink

	err *widget.Label

	progress *widget.ProgressBarInfinite

	name            *widget.Label
	openPlaylistBtn *widget.Button

	playlistData streaming.SoundCloudPlaylist

	ctxCancel func() // used to cancel a downloading context
}

/*
NewPlaylist returns a new instance of the Playlist widget
*/
func NewPlaylist(openPlaylistFunc func(streaming.SoundCloudPlaylist)) *Playlist {

	i := &Playlist{
		name:            widget.NewLabel(""),
		searchUrl:       widget.NewHyperlink("", nil),
		err:             widget.NewLabel(""),
		progress:        widget.NewProgressBarInfinite(),
		openPlaylistBtn: widget.NewButton("Open Playlist", func() {}),
	}

	i.openPlaylistBtn.OnTapped = func() {
		openPlaylistFunc(i.playlistData)
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

	i.ExtendBaseWidget(i)

	return i
}

func (i *Playlist) CreateRenderer() fyne.WidgetRenderer {

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

func (i *Playlist) SetFinding(playlist streaming.SoundCloudPlaylist) {
	i.playlistData = playlist

	i.findingContent.Show()
	i.foundContent.Hide()
	i.failedContent.Hide()
}

func (i *Playlist) SetFound(playlist streaming.SoundCloudPlaylist) {
	i.playlistData = playlist

	i.findingContent.Hide()
	i.foundContent.Show()
	i.failedContent.Hide()
	i.name.SetText(playlist.Name)
}

func (i *Playlist) SetFailed(err error) {
	i.findingContent.Hide()
	i.foundContent.Hide()
	i.failedContent.Show()
	i.err.SetText(err.Error()) // TODO: unwrapping of error
}

/*
AddPlaylist displays a section used for adding a playlist to the ui
*/
type AddPlaylist struct {
	widget.BaseWidget

	urlEntry        *widget.Entry
	submitButton    *widget.Button
	validationLabel *widget.Label

	OnAdd func()
}

/*
NewAddPlaylist returns a new instance of AddPlaylist widget
*/
func NewAddPlaylist(addPlaylist func(string)) *AddPlaylist {

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

	i := &AddPlaylist{
		submitButton:    submitBtn,
		urlEntry:        urlEntry,
		validationLabel: validationLabel,
	}

	widget.NewForm()

	i.ExtendBaseWidget(i)

	return i
}

func (i *AddPlaylist) CreateRenderer() fyne.WidgetRenderer {
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
UpdateFromData updates the playlistWidget for a given playlist

This is called when a playlist is added, or when the state of the playlist changes,
for example when the request to get information about the playlist fails, or is completed
*/
func (p *Playlist) UpdateFromData(playlistBindingItem *PlaylistBindingItem) {

	playlist := playlistBindingItem.playlist

	p.searchUrl.SetText(playlist.SearchUrl)
	p.searchUrl.SetURLFromString(playlist.SearchUrl)

	switch playlistBindingItem.state {
	case Finding:
		p.SetFinding(playlistBindingItem.playlist)
	case Failed:
		p.SetFailed(playlistBindingItem.err)
	case Found:
		p.SetFound(playlistBindingItem.playlist)
	}

	p.Refresh()
}
