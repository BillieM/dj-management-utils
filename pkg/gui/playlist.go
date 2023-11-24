package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
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
		content = container.NewHBox(
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

	submitButton *widget.Button
	urlEntry     *widget.Entry
}

/*
newAddPlaylistWidget returns a new instance of addPlaylistWidget
*/
func newAddPlaylistWidget(p *playlistBindingList, onSubmit func()) *addPlaylistWidget {

	urlEntry := widget.NewEntry()
	submitBtn := widget.NewButton("add playlist", func() {
		err := urlEntry.Validate()
		fmt.Println(err)
		p.addPlaylist(urlEntry.Text, onSubmit)
	})

	urlEntry.SetPlaceHolder("SoundCloud playlist url")
	urlEntry.Validator = validation.NewRegexp(`soundcloud\.com\/.*\/sets`, "not a valid SoundCloud playlist url")

	urlEntry.OnSubmitted = func(s string) { p.addPlaylist(s, onSubmit) }

	i := &addPlaylistWidget{
		submitButton: submitBtn,
		urlEntry:     urlEntry,
	}

	i.ExtendBaseWidget(i)

	return i
}

func (item *addPlaylistWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewBorder(
		widget.NewLabel("Add playlist"),
		nil, nil, nil,
		container.NewBorder(
			nil, nil, nil, item.submitButton, item.urlEntry,
		),
	)
	return widget.NewSimpleRenderer(c)
}

func (p *playlistBindingList) addPlaylist(url string, callback func()) {

	fmt.Println("addPlaylist", url)
	p.Append(&playlistBindingItem{
		url:         url,
		downloading: true,
		failed:      false,
	})
	callback()
}
