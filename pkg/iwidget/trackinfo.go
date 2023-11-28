package iwidget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/database"
)

type TrackInfo struct {
	widget.BaseWidget

	TrackNameLink   *widget.Hyperlink
	TrackProperties *TrackProperties
}

func NewTrackInfo(t database.SoundCloudTrack) *TrackInfo {

	i := &TrackInfo{
		TrackNameLink:   widget.NewHyperlink(t.Name, nil),
		TrackProperties: NewTrackProperties(t),
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *TrackInfo) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			i.TrackNameLink,
			i.TrackProperties,
		),
	)
}

type TrackProperties struct {
	widget.BaseWidget

	GenrePropertyLabel          *TrackProperty
	TagListPropertyLabel        *TrackProperty
	PublisherPropertyLabel      *TrackProperty
	SoundCloudUserPropertyLabel *TrackProperty
}

func NewTrackProperties(t database.SoundCloudTrack) *TrackProperties {

	fmt.Println(t.Genre)

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

func (l *TrackProperty) SetText(text string) {
	l.PropertyLabel.Text = text
	l.PropertyLabel.Refresh()
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
