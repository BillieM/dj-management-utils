package gui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
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
}

/*
newPlaylistWidget returns a new instance of playlistWidget
*/
func newPlaylistWidget(name string) *playlistWidget {
	i := &playlistWidget{
		name: widget.NewLabel(name),
	}
	i.ExtendBaseWidget(i)

	return i
}

func (item *playlistWidget) CreateRenderer() fyne.WidgetRenderer {
	c := container.NewBorder(nil, nil, nil, nil, item.name)
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
	i := &addPlaylistWidget{
		submitButton: widget.NewButton("add playlist", func() {}),
		urlEntry:     widget.NewEntry(),
	}
	i.ExtendBaseWidget(i)

	i.urlEntry.OnSubmitted = func(s string) {
		p.addPlaylist(s, onSubmit)
	}
	i.submitButton.OnTapped = func() {
		p.addPlaylist(i.urlEntry.Text, onSubmit)
	}

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
	if url == "" {
		return
	}
	fmt.Println("addPlaylist", url)
	p.Append(&playlistBindingItem{
		url: url,
	})
	callback()
}
