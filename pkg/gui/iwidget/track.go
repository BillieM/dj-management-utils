package iwidget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/database"
)

type TrackSection struct {
	widget.BaseWidget

	Track       *Track
	Placeholder *widget.Label
}

func NewTrackSection(t database.SoundCloudTrack, downloadFunc func(database.SoundCloudTrack)) *TrackSection {

	track := NewTrack(t, downloadFunc)

	i := &TrackSection{
		Track:       track,
		Placeholder: widget.NewLabel("Please select a track..."),
	}

	i.Track.Hide()
	i.Placeholder.Show()

	i.ExtendBaseWidget(i)

	return i
}

func (t *TrackSection) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewStack(
			t.Track,
			t.Placeholder,
		),
	)
}

func (trackSection *TrackSection) Bind(selectedTrack *SelectedTrackBinding) {
	listener := binding.NewDataListener(func() {
		trackSection.updateFromData(selectedTrack.TrackBinding)
	})

	selectedTrack.AddListener(listener)
}

func (trackSection *TrackSection) updateFromData(b *TrackBinding) {
	if b.track != nil {
		trackSection.Track.updateFromData(b)
		trackSection.Track.Show()
		trackSection.Placeholder.Hide()
	} else {
		trackSection.Track.Hide()
		trackSection.Placeholder.Show()
	}
}

/*
Track is the main widget used for displaying a track selected from the track list
inside of a playlist view
*/
type Track struct {
	widget.BaseWidget

	// track info
	TrackInfo *TrackInfo

	// get track
	GetTrack *GetTrack

	// link track
	LinkTrack *LinkTrack
}

func NewTrack(t database.SoundCloudTrack, downloadFunc func(database.SoundCloudTrack)) *Track {

	i := &Track{
		TrackInfo: NewTrackInfo(t),
		GetTrack:  NewGetTrack(downloadFunc),
		LinkTrack: NewLinkTrack(),
	}

	i.ExtendBaseWidget(i)

	return i
}

func (t *Track) CreateRenderer() fyne.WidgetRenderer {

	c := container.NewVScroll(
		container.NewVBox(
			t.TrackInfo,
			widget.NewSeparator(),
			t.GetTrack,
			t.LinkTrack,
		),
	)

	return widget.NewSimpleRenderer(c)
}

func (t *Track) updateFromData(b *TrackBinding) {
	t.TrackInfo.updateFromData(*b.track)

	if b.track.HasDownloadsLeft || b.track.PurchaseTitle != "" {
		t.GetTrack.Show()
	} else {
		t.GetTrack.Hide()
	}
	t.GetTrack.updateFromData(*b.track)
	t.LinkTrack.updateFromData(*b.track)
}

type SelectedTrackBinding struct {
	bindBase

	TrackBinding *TrackBinding
}

func (i *SelectedTrackBinding) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *SelectedTrackBinding) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (i *SelectedTrackBinding) trigger() (interface{}, error) {

	i.listeners.Range(func(key, _ interface{}) bool {
		key.(binding.DataListener).DataChanged()
		return true
	})

	return nil, nil
}

/*
TrackInfo displays the track name, link to the track, and track properties
(i.e. genre, tags, publisher, soundcloud user)
*/
type TrackInfo struct {
	widget.BaseWidget

	TrackNameLink   *widget.Hyperlink
	TrackLinkButton *OpenInBrowserButton
	TrackProperties *TrackProperties
}

func NewTrackInfo(t database.SoundCloudTrack) *TrackInfo {

	trackNameLink := widget.NewHyperlink(t.Name, nil)
	trackLinkButton := NewOpenInBrowserButton("Open in browser", "")

	if trackLinkButton.URL != nil {
		err := fyne.CurrentApp().OpenURL(trackLinkButton.URL)
		if err != nil {
			fyne.LogError("Failed to open url", err)
		}
	}

	i := &TrackInfo{
		TrackNameLink:   trackNameLink,
		TrackLinkButton: trackLinkButton,
		TrackProperties: NewTrackProperties(t),
	}

	i.TrackNameLink.TextStyle.Bold = true

	i.TrackLinkButton.OnTapped = func() {}

	i.ExtendBaseWidget(i)

	return i
}

func (i *TrackInfo) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			container.NewBorder(
				nil, nil, i.TrackNameLink, i.TrackLinkButton,
			),
			widget.NewSeparator(),
			i.TrackProperties,
		),
	)
}

func (i *TrackInfo) updateFromData(t database.SoundCloudTrack) {
	i.TrackNameLink.SetURLFromString(t.Permalink)
	i.TrackNameLink.SetText(t.Name)
	i.TrackLinkButton.SetContent("Open in browser", t.Permalink)
	i.TrackProperties.updateFromData(t)
}

type TrackProperties struct {
	widget.BaseWidget

	GenrePropertyLabel          *TrackProperty
	TagListPropertyLabel        *TrackProperty
	PublisherPropertyLabel      *TrackProperty
	SoundCloudUserPropertyLabel *TrackProperty
}

func NewTrackProperties(t database.SoundCloudTrack) *TrackProperties {

	genrePropertyLabel := NewTrackProperty("Genre", t.Genre)
	tagListPropertyLabel := NewTrackProperty("Tags", t.TagList)
	publisherPropertyLabel := NewTrackProperty("Publisher", t.PublisherArtist)
	soundCloudUserPropertyLabel := NewTrackProperty("SoundCloud User", t.SoundCloudUser)

	i := &TrackProperties{
		GenrePropertyLabel:          genrePropertyLabel,
		TagListPropertyLabel:        tagListPropertyLabel,
		PublisherPropertyLabel:      publisherPropertyLabel,
		SoundCloudUserPropertyLabel: soundCloudUserPropertyLabel,
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *TrackProperties) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			i.GenrePropertyLabel,
			i.TagListPropertyLabel,
			i.PublisherPropertyLabel,
			i.SoundCloudUserPropertyLabel,
		),
	)
}

func (i *TrackProperties) updateFromData(t database.SoundCloudTrack) {
	i.GenrePropertyLabel.Update(t.Genre)
	i.TagListPropertyLabel.Update(t.TagList)
	i.PublisherPropertyLabel.Update(t.PublisherArtist)
	i.SoundCloudUserPropertyLabel.Update(t.SoundCloudUser)
}

type TrackProperty struct {
	widget.BaseWidget
	labelLabel    *widget.Label
	PropertyLabel *widget.Label
}

func NewTrackProperty(propertyName string, propertyValue string) *TrackProperty {

	labelLabel := widget.NewLabel(fmt.Sprintf("%s:", propertyName))
	propertyLabel := widget.NewLabelWithStyle(propertyValue, fyne.TextAlignLeading, fyne.TextStyle{Bold: true})
	i := &TrackProperty{
		labelLabel:    labelLabel,
		PropertyLabel: propertyLabel,
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *TrackProperty) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewHBox(
			i.labelLabel,
			i.PropertyLabel,
		),
	)
}

func (l *TrackProperty) Update(text string) {
	l.PropertyLabel.SetText(text)
}

type TrackNameLink struct {
	widget.Hyperlink
}

func NewTrackNameLink(text string) *TrackNameLink {
	i := &TrackNameLink{}

	i.ExtendBaseWidget(i)

	return i
}

func (i *TrackNameLink) SetNameLinkFromString(name string, url string) {
	i.SetURLFromString(url)
	i.SetText(name)
}

/*
 */
type GetTrack struct {
	widget.BaseWidget

	TrackDownload *TrackDownload
	TrackPurchase *TrackPurchase
}

func NewGetTrack(downloadFunc func(database.SoundCloudTrack)) *GetTrack {
	i := &GetTrack{
		TrackDownload: NewTrackDownload(downloadFunc),
		TrackPurchase: NewTrackPurchase(),
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *GetTrack) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			i.TrackDownload,
			i.TrackPurchase,
			widget.NewSeparator(),
		),
	)
}

func (i *GetTrack) updateFromData(t database.SoundCloudTrack) {

	if t.HasDownloadsLeft {
		i.TrackDownload.Show()
		i.TrackDownload.updateFromData(t)
	} else {
		i.TrackDownload.Hide()
	}

	if t.PurchaseTitle != "" {
		i.TrackPurchase.Show()
		i.TrackPurchase.updateFromData(t)
	} else {
		i.TrackPurchase.Hide()
	}
}

/*
TrackPurchase widget handles purchase/ free links visible on SoundCloud
*/
type TrackPurchase struct {
	widget.BaseWidget

	TrackPurchaseButton *OpenInBrowserButton
}

func NewTrackPurchase() *TrackPurchase {
	i := &TrackPurchase{
		TrackPurchaseButton: NewOpenInBrowserButton("Purchase Track", ""),
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *TrackPurchase) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewHBox(
			widget.NewLabel("Get track (opens in browser)"),
			i.TrackPurchaseButton,
		),
	)
}

func (i *TrackPurchase) updateFromData(t database.SoundCloudTrack) {
	i.TrackPurchaseButton.SetContent(t.PurchaseTitle, t.PurchaseURL)
}

/*
TrackDownload widget handles downloads directly from SoundCloud, i.e.
when a track has a 'download file' option within the SoundCloud UI.

This is different from the 'free download'/ 'buy' options, which are
handled by 'TrackPurchase'
*/

type TrackDownload struct {
	widget.BaseWidget

	DownloadFunc func(database.SoundCloudTrack)

	TrackDownloadButton   *widget.Button
	TrackDownloadProgress *widget.ProgressBarInfinite
}

func NewTrackDownload(downloadFunc func(database.SoundCloudTrack)) *TrackDownload {

	trackDownloadButton := widget.NewButton("Download Track", func() {})

	trackDownloadProgress := widget.NewProgressBarInfinite()
	trackDownloadProgress.Hide()

	i := &TrackDownload{
		TrackDownloadButton:   trackDownloadButton,
		TrackDownloadProgress: trackDownloadProgress,
		DownloadFunc:          downloadFunc,
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *TrackDownload) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewBorder(
			nil, nil,
			container.NewHBox(
				widget.NewLabel("Get track (save to output directory)"),
				i.TrackDownloadButton,
			),
			nil,
			i.TrackDownloadProgress,
		),
	)
}

func (i *TrackDownload) updateFromData(t database.SoundCloudTrack) {
	i.TrackDownloadButton.OnTapped = func() {
		i.TrackDownloadProgress.Show()
		go func() {
			i.DownloadFunc(t)
			i.TrackDownloadProgress.Hide()
		}()
	}
}

/*
LinkTrack allows for the user to establish a link between a SoundCloud track and
a track within their DJ libary/ local filesystem.
*/
type LinkTrack struct {
	widget.BaseWidget
}

func NewLinkTrack() *LinkTrack {
	i := &LinkTrack{}

	i.ExtendBaseWidget(i)

	return i
}

func (i *LinkTrack) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		widget.NewLabel("Link Track"),
	)
}

func (i *LinkTrack) updateFromData(t database.SoundCloudTrack) {

}
