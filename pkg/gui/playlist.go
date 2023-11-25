package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
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
func (i *playlistBindingList) load(s *database.SerenDB) {

	// TODO err handling...
	playlists, _ := s.GetSoundCloudPlaylists()

	for _, playlist := range playlists {
		i.Append(&playlistBindingItem{
			playlist: playlist,
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
	playlist database.SoundCloudPlaylist
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
	name      *widget.Label
	err       *widget.Label
	progress  *widget.ProgressBarInfinite

	ctxCancel func() // used to cancel a downloading context
}

/*
newPlaylistWidget returns a new instance of playlistWidget
*/
func newPlaylistWidget() *playlistWidget {
	i := &playlistWidget{}

	i.ExtendBaseWidget(i)

	return i
}

func (i *playlistWidget) CreateRenderer() fyne.WidgetRenderer {

	i.searchUrl = widget.NewHyperlink("", nil)

	i.name = widget.NewLabel("")

	i.err = widget.NewLabel("")
	i.err.Importance = widget.WarningImportance

	i.progress = widget.NewProgressBarInfinite()

	i.findingContent = i.progress
	i.failedContent = i.err
	i.foundContent = container.NewBorder(
		i.name,
		nil,
		nil,
		nil,
	)

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

type streamingStepHandler struct {
	stepFunc     func()
	finishedFunc func()
}

func (s streamingStepHandler) StepCallback(step operations.StepInfo) {
	fmt.Println(step)
	s.stepFunc()
}

func (s streamingStepHandler) ExitCallback() {
	fmt.Println("finished")
	s.finishedFunc()
}
