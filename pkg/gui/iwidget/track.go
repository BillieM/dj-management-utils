package iwidget

import (
	"fmt"
	"net/url"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/database"
)

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

func NewTrack(t database.SoundCloudTrack) *Track {

	i := &Track{
		TrackInfo: NewTrackInfo(t),
		GetTrack:  NewGetTrack(t.PurchaseTitle),
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
			widget.NewSeparator(),
			t.LinkTrack,
		),
	)

	return widget.NewSimpleRenderer(c)
}

func (trackWidget *Track) Bind(selectedTrack *SelectedTrackBinding) {
	listener := binding.NewDataListener(func() {
		trackWidget.updateFromData(selectedTrack.trackBinding)
	})

	selectedTrack.AddListener(listener)
}

func (t *Track) updateFromData(b *TrackBinding) {
	t.TrackInfo.Update(*b.track)
}

type SelectedTrackBinding struct {
	bindBase

	trackBinding *TrackBinding
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

	i.TrackLinkButton.OnTapped = func() {

	}

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

func (i *TrackInfo) Update(t database.SoundCloudTrack) {
	i.TrackNameLink.SetURLFromString(t.Permalink)
	i.TrackNameLink.SetText(t.Name)
	i.TrackLinkButton.SetURLFromString(t.Permalink)
	i.TrackProperties.Update(t)
}

type OpenInBrowserButton struct {
	*widget.Button

	URL *url.URL
}

func NewOpenInBrowserButton(text string, urlString string) *OpenInBrowserButton {

	openInBrowserBtn := &OpenInBrowserButton{}

	btn := widget.NewButton(text, func() {})

	if urlString != "" {
		openInBrowserBtn.SetURLFromString(urlString)
	}
	openInBrowserBtn.Button = btn

	return openInBrowserBtn
}

func (i *OpenInBrowserButton) setOpenFunc() {
	i.OnTapped = func() {
		err := fyne.CurrentApp().OpenURL(i.URL)
		if err != nil {
			fyne.LogError("Failed to open url", err)
		}
	}
}

func (i *OpenInBrowserButton) SetURLFromString(urlStr string) {
	u, err := url.Parse(urlStr)
	if err != nil {
		fyne.LogError("Failed to parse url", err)
	}
	i.URL = u
	i.setOpenFunc()
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

func (i *TrackProperties) Update(t database.SoundCloudTrack) {
	i.GenrePropertyLabel.Update(t.Genre)
	i.TagListPropertyLabel.Update(t.TagList)
	i.PublisherPropertyLabel.Update(t.PublisherArtist)
	i.SoundCloudUserPropertyLabel.Update(t.SoundCloudUser)
}

type TrackProperty struct {
	widget.BaseWidget
	labelLabel    *widget.Label
	PropertyLabel *EmphasizedLabel
}

func NewTrackProperty(propertyName string, propertyValue string) *TrackProperty {

	labelLabel := widget.NewLabel(fmt.Sprintf("%s:", propertyName))
	propertyLabel := NewEmphasizedLabel(propertyValue)
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
	TrackPurchase *widget.Button
}

func NewGetTrack(purchaseBtnText string) *GetTrack {
	i := &GetTrack{
		TrackDownload: NewTrackDownload(),
		TrackPurchase: widget.NewButton("purchaseBtnText", func() {}),
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *GetTrack) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewGridWithColumns(
			2,
			i.TrackDownload,
			i.TrackPurchase,
		),
	)
}

/*
TrackDownload widget handles downloads directly from SoundCloud, i.e.
when a track has a 'download file' option within the SoundCloud UI.

This is different from the 'free download'/ 'buy' options, which are
handled by 'TrackPurchase'
*/

type TrackDownload struct {
	widget.BaseWidget

	TrackDownloadButton   *widget.Button
	TrackDownloadProgress *widget.ProgressBarInfinite
}

func NewTrackDownload() *TrackDownload {

	trackDownloadButton := widget.NewButton("Download Track", func() {})

	i := &TrackDownload{
		TrackDownloadButton:   trackDownloadButton,
		TrackDownloadProgress: widget.NewProgressBarInfinite(),
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *TrackDownload) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			i.TrackDownloadButton,
			i.TrackDownloadProgress,
		),
	)
}

func (i *TrackDownload) Update() {

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
