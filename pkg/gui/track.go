package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
)

type trackWidget struct {
	widget.BaseWidget
	name *widget.Label
	url  *widget.Hyperlink

	purchaseBtn     *widget.Button
	downloadFileBtn *widget.Button
	downloadingFile *widget.ProgressBarInfinite

	downloadFileCanvas fyne.CanvasObject
	purchaseCanvas     fyne.CanvasObject

	genre   *widget.Label
	tagList *widget.Label

	publisherArtist *widget.Label
	soundCloudUser  *widget.Label
}

func newTrackWidget() *trackWidget {

	// just build a widget constructor or smtn for these

	nameLabel := widget.NewLabel("")
	nameLabel.TextStyle.Bold = true
	nameLabel.Importance = widget.HighImportance

	genreLabel := widget.NewLabel("")
	genreLabel.TextStyle.Bold = true
	genreLabel.Importance = widget.HighImportance

	tagListLabel := widget.NewLabel("")
	tagListLabel.TextStyle.Bold = true
	tagListLabel.Importance = widget.HighImportance

	publisherArtistLabel := widget.NewLabel("")
	publisherArtistLabel.TextStyle.Bold = true
	publisherArtistLabel.Importance = widget.HighImportance

	soundCloudUserLabel := widget.NewLabel("")
	soundCloudUserLabel.TextStyle.Bold = true
	soundCloudUserLabel.Importance = widget.HighImportance

	downloadFileBtn := widget.NewButton("Download File", func() {})
	purchaseBtn := widget.NewButton("", func() {})
	downloadingFile := widget.NewProgressBarInfinite()

	downloadingFile.Hide()

	i := &trackWidget{
		name:            nameLabel,
		url:             widget.NewHyperlink("", nil),
		purchaseBtn:     purchaseBtn,
		downloadFileBtn: downloadFileBtn,
		downloadingFile: downloadingFile,
		genre:           genreLabel,
		tagList:         tagListLabel,
		publisherArtist: publisherArtistLabel,
		soundCloudUser:  soundCloudUserLabel,
	}

	i.ExtendBaseWidget(i)

	return i
}

func (t *trackWidget) CreateRenderer() fyne.WidgetRenderer {

	c := container.NewBorder(
		container.NewVBox(
			t.name,
			t.url,
		),
		nil,
		nil,
		nil,
		container.NewVBox(
			container.NewGridWithColumns(
				4,
				container.NewHBox(widget.NewLabel("Genre: "), t.genre),
				container.NewHBox(widget.NewLabel("Tags: "), t.tagList),
				container.NewHBox(widget.NewLabel("Publisher: "), t.publisherArtist),
				container.NewHBox(widget.NewLabel("SoundCloud User: "), t.soundCloudUser),
			),
			container.NewGridWithColumns(
				2,
				container.NewHBox(t.downloadFileBtn, t.downloadingFile),
				t.purchaseBtn,
			),
		),
	)

	return widget.NewSimpleRenderer(c)
}

/*
playlistBindingList stores a list of playlistBindingItem structs

It is used to display a list of playlists as playlistWidgets in the UI
*/
type trackBindingList struct {
	bindBase

	Items []*trackBindingItem
}

func (i *trackBindingList) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *trackBindingList) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (i *trackBindingList) GetItem(index int) (binding.DataItem, error) {
	i.Lock()
	defer i.Unlock()
	if index < 0 || index >= len(i.Items) {
		return nil, helpers.ErrIndexOutOfBounds
	}

	return i.Items[index], nil
}

func (i *trackBindingList) Length() int {
	i.Lock()
	defer i.Unlock()
	return len(i.Items)
}

func (i *trackBindingList) Append(p *trackBindingItem) {
	i.Lock()
	defer i.Unlock()
	i.Items = append(i.Items, p)
}

/*
load loads all tracks from a given playlist id into the trackBindingList
*/
func (i *trackBindingList) load(s *database.SerenDB, externalID int64) {

	// TODO err handling...
	tracks, _ := s.GetSoundCloudTracks(externalID)

	for _, track := range tracks {
		i.Append(&trackBindingItem{
			track: track,
		})
	}
}

/*
playlistBindingItem is a struct that contains the data for a playlist

It is used to display a playlist as a playlistWidget in the UI
*/
type trackBindingItem struct {
	bindBase

	// may want a context in here ?? later problem...
	track database.SoundCloudTrack
}

func (i *trackBindingItem) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *trackBindingItem) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (e *guiEnv) updateTracksList(trackWidget *trackWidget, trackBindingItem *trackBindingItem, playlistName string) {
	track := trackBindingItem.track

	trackWidget.name.SetText(track.Name)

	trackWidget.url.SetText(track.Permalink)
	trackWidget.url.SetURLFromString(track.Permalink)

	trackWidget.genre.SetText(track.Genre)
	trackWidget.tagList.SetText(track.TagList)
	trackWidget.publisherArtist.SetText(track.PublisherArtist)
	trackWidget.soundCloudUser.SetText(track.SoundCloudUser)

	if track.HasDownloadsLeft {
		trackWidget.downloadFileBtn.SetText("Download File")
		trackWidget.downloadFileBtn.Enable()
		trackWidget.downloadFileBtn.OnTapped = func() {
			trackWidget.downloadFileBtn.Disable()

			opEnv := e.opEnv()

			go func() {
				trackWidget.downloadingFile.Show()
				opEnv.DownloadSoundCloudFile(track, playlistName)
				trackWidget.downloadingFile.Hide()
				trackWidget.downloadFileBtn.Enable()
			}()
		}
	} else {
		trackWidget.downloadFileBtn.SetText("No Download")
		trackWidget.downloadFileBtn.Disable()
	}

	if track.PurchaseTitle != "" {
		trackWidget.purchaseBtn.SetText(track.PurchaseTitle)
		trackWidget.purchaseBtn.Enable()
		trackWidget.purchaseBtn.OnTapped = func() {
			opEnv := e.opEnv()

			go func() {
				trackWidget.purchaseBtn.Disable()
				opEnv.OpenSoundCloudPurchase(track)
				trackWidget.purchaseBtn.Enable()
			}()
		}
	} else {
		trackWidget.purchaseBtn.SetText("No Purchase Link")
		trackWidget.purchaseBtn.Disable()
	}

	trackWidget.Refresh()
}
