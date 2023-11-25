package gui

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

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
playlistBindingItem is a struct that contains the data for a playlist

It is used to display a playlist as a playlistWidget in the UI
*/
type playlistBindingItem struct {
	bindBase

	// may want a context in here ?? later problem...

	downloading bool
	failed      bool

	name string
	url  string
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
	name *widget.Label
	url  *widget.Hyperlink

	ctxCancel func() // used to cancel a downloading context

	downloading bool
	failed      bool
}

/*
newPlaylistWidget returns a new instance of playlistWidget
*/
func newPlaylistWidget(name string) *playlistWidget {
	i := &playlistWidget{
		name: widget.NewLabel(name),
		url:  widget.NewHyperlink("playlist url", nil),
	}
	i.ExtendBaseWidget(i)

	return i
}

func (i *playlistWidget) CreateRenderer() fyne.WidgetRenderer {

	var content *fyne.Container

	if i.failed {
		content = container.NewHBox(
			widget.NewLabel("download failed, click to retry"),
		)
	} else if i.downloading {
		cancelBtn := widget.NewButton("cancel", func() {
			i.ctxCancel()
		})

		content = container.NewBorder(
			nil, nil, nil, cancelBtn,
			widget.NewProgressBarInfinite(),
		)
	} else {
		content = container.NewHBox(
			i.name,
		)
	}

	c := container.NewBorder(i.url, nil, nil, nil, content)

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
func newAddPlaylistWidget(p *playlistBindingList, onSubmit func(*playlistBindingItem)) *addPlaylistWidget {

	urlEntry := widget.NewEntry()
	urlEntry.SetPlaceHolder("SoundCloud Playlist URL")
	urlEntry.Validator = validation.NewRegexp(`soundcloud\.com\/.*\/sets`, "not a valid SoundCloud playlist url")

	urlEntry.OnSubmitted = func(s string) {

		var err error

		if err = urlEntry.Validate(); err == nil {
			err = p.addPlaylist(s, onSubmit)
		}

		if err != nil {
			urlEntry.SetValidationError(err)
			urlEntry.Refresh()
		}
	}

	submitBtn := widget.NewButton("Add Playlist", func() {
		var err error

		if err = urlEntry.Validate(); err == nil {
			err = p.addPlaylist(urlEntry.Text, onSubmit)
		}

		if err != nil {
			urlEntry.SetValidationError(err)
			urlEntry.Refresh()
		}
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

func (p *playlistBindingList) addPlaylist(urlStr string, callback func(*playlistBindingItem)) error {

	u, err := url.Parse(urlStr)

	if err != nil {
		return err
	}

	u.RawQuery = ""

	i := &playlistBindingItem{
		url:         fmt.Sprint(u),
		downloading: true,
		failed:      false,
	}

	p.Append(i)
	callback(i)

	return nil
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
