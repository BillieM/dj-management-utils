package iwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/database"
)

type Track struct {
	widget.BaseWidget

	// track info
	*TrackInfo

	// get track
	*GetTrack

	// link track
	*LinkTrack
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

	c := container.NewVBox(
		t.TrackInfo,
		t.GetTrack,
		t.LinkTrack,
	)

	return widget.NewSimpleRenderer(c)
}
