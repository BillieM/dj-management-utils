package iwidget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/gui/uihelpers"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/streaming"
)

type TrackSection struct {
	widget.BaseWidget

	Track       *Track
	Placeholder *widget.Label
}

func NewTrackSection(w fyne.Window, trackFuncs TrackFuncs, resizeEvents *uihelpers.ResizeEvents) *TrackSection {

	track := NewTrack(w, trackFuncs, resizeEvents)

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

/*
Bind binds the TrackSection to a SelectedTrackBinding
This allows the TrackSection to update when the selected track changes

note: currently a bit of a hack, and the Trigger() method should be called
to make the updates reflect
*/
func (trackSection *TrackSection) Bind(selectedTrack *SelectedTrackBinding) {
	listener := binding.NewDataListener(func() {
		trackSection.updateFromData(selectedTrack)
	})

	selectedTrack.AddListener(listener)
}

func (trackSection *TrackSection) updateFromData(b *SelectedTrackBinding) {
	if b.TrackBinding.Track != nil {
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

func NewTrack(w fyne.Window, trackFuncs TrackFuncs, resizeEvents *uihelpers.ResizeEvents) *Track {

	i := &Track{
		TrackInfo: NewTrackInfo(w),
		GetTrack:  NewGetTrack(w, trackFuncs.DownloadSoundCloudTrack),
		LinkTrack: NewLinkTrack(
			w,
			trackFuncs.SaveSoundCloudTrackToDB,
			trackFuncs.OnError,
			resizeEvents,
		),
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

func (t *Track) updateFromData(b *SelectedTrackBinding) {
	scTrack := *b.TrackBinding.Track

	t.TrackInfo.updateFromData(scTrack)

	if scTrack.HasDownloadsLeft || scTrack.PurchaseTitle != "" {
		t.GetTrack.Show()
	} else {
		t.GetTrack.Hide()
	}
	t.GetTrack.updateFromData(scTrack)
	t.LinkTrack.updateFromData(b)
}

type SelectedTrackBinding struct {
	bindBase

	Locked       bool
	ListID       widget.ListItemID
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

/*
Trigger calls the DataChanged method on all listeners

Somewhat of a hack, but it saves using Set methods
*/
func (i *SelectedTrackBinding) Trigger() {

	i.listeners.Range(func(key, _ interface{}) bool {
		key.(binding.DataListener).DataChanged()
		return true
	})

}

func (i *SelectedTrackBinding) LockSelected() {
	i.Lock()
	defer i.Unlock()
	i.Locked = true
}

func (i *SelectedTrackBinding) UnlockSelected() {
	i.Lock()
	defer i.Unlock()
	i.Locked = false
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

func NewTrackInfo(w fyne.Window) *TrackInfo {

	trackNameLink := widget.NewHyperlink("", nil)
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
		TrackProperties: NewTrackProperties(),
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

func (i *TrackInfo) updateFromData(t streaming.SoundCloudTrack) {
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

func NewTrackProperties() *TrackProperties {

	genrePropertyLabel := NewTrackProperty("Genre", "")
	tagListPropertyLabel := NewTrackProperty("Tags", "")
	publisherPropertyLabel := NewTrackProperty("Publisher", "")
	soundCloudUserPropertyLabel := NewTrackProperty("SoundCloud User", "")

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

func (i *TrackProperties) updateFromData(t streaming.SoundCloudTrack) {
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

func NewGetTrack(w fyne.Window, downloadFunc func()) *GetTrack {
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

func (i *GetTrack) updateFromData(t streaming.SoundCloudTrack) {

	if t.HasDownloadsLeft {
		i.TrackDownload.Show()
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

func (i *TrackPurchase) updateFromData(t streaming.SoundCloudTrack) {
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

	TrackDownloadButton   *widget.Button
	TrackDownloadProgress *widget.ProgressBarInfinite
}

func NewTrackDownload(downloadFunc func()) *TrackDownload {
	trackDownloadProgress := widget.NewProgressBarInfinite()
	trackDownloadProgress.Hide()

	trackDownloadButton := widget.NewButton("Download Track",
		func() {
			go func() {
				trackDownloadProgress.Show()
				downloadFunc()
				trackDownloadProgress.Hide()
			}()
		})

	i := &TrackDownload{
		TrackDownloadButton:   trackDownloadButton,
		TrackDownloadProgress: trackDownloadProgress,
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

/*
LinkTrack allows for the user to establish a link between a SoundCloud track and
a track within their DJ libary/ local filesystem.
*/
type LinkTrack struct {
	widget.BaseWidget

	parentWindow fyne.Window

	LinkTrackFileSelect *LinkTrackFileSelect
}

func NewLinkTrack(w fyne.Window, saveSoundCloudTrackFunc func(), onError func(error), resizeEvents *uihelpers.ResizeEvents) *LinkTrack {
	i := &LinkTrack{
		parentWindow:        w,
		LinkTrackFileSelect: NewLinkTrackFileSelect(w, saveSoundCloudTrackFunc, onError, resizeEvents),
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *LinkTrack) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			i.LinkTrackFileSelect,
		),
	)
}

func (i *LinkTrack) updateFromData(t *SelectedTrackBinding) {
	i.LinkTrackFileSelect.updateFromData(t)
}

/*
LinkTrackFileSelect allows for the user to link a SoundCloud track to a file
on their local file system via a file selection dialog
*/
type LinkTrackFileSelect struct {
	widget.BaseWidget

	parentWindow fyne.Window

	resizeEvents            *uihelpers.ResizeEvents
	saveSoundCloudTrackFunc func()

	OpenPath *OpenPath
}

func NewLinkTrackFileSelect(w fyne.Window, saveSoundCloudTrackFunc func(), onError func(error), resizeEvents *uihelpers.ResizeEvents) *LinkTrackFileSelect {

	openPath := NewOpenPath(w, "", File)

	openPath.SetExtensionFilter(helpers.GetAudioExtensions())
	openPath.SetOnError(onError)

	resizeFunc := func() {
		openPath.Dialog.Resize(uihelpers.CanvasPercentSize(w, 0.75, 0.75, fyne.NewSize(480, 320), fyne.NewSize(1280, 0)))
	}

	var key string

	openPath.SetOnOpen(func() {
		resizeFunc()
		key = resizeEvents.Add(resizeFunc)
	})
	openPath.SetOnClose(func() {
		resizeEvents.Remove(key)
	})

	i := &LinkTrackFileSelect{
		parentWindow:            w,
		OpenPath:                openPath,
		resizeEvents:            resizeEvents,
		saveSoundCloudTrackFunc: saveSoundCloudTrackFunc,
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *LinkTrackFileSelect) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			widget.NewLabel("Select track file location from local file system"),
			container.NewBorder(
				nil, nil,
				layout.NewSpacer(), layout.NewSpacer(),
				i.OpenPath,
			),
		),
	)
}

func (i *LinkTrackFileSelect) updateFromData(t *SelectedTrackBinding) {
	scTrack := *t.TrackBinding.Track

	if scTrack.LocalPath != "" {
		i.OpenPath.SetURIFromPathString(scTrack.LocalPath)
	} else {
		i.OpenPath.SetURIFromPathString("/")
	}

	i.OpenPath.SetOnValid(func(uri string) {
		t.TrackBinding.Track.LocalPath = uri
		i.saveSoundCloudTrackFunc()
		t.Trigger()
	})
}

type TrackFuncs struct {
	DownloadSoundCloudTrack func()
	SaveSoundCloudTrackToDB func()
	OnError                 func(error)
}
