package iwidget

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/database"
)

type Track struct {
	widget.BaseWidget

	// track info
	TrackInfo *TrackInfo

	// get track
	GetTrack *GetTrack

	// link track
	LinkTrack *LinkTrack

	binder basicBinder
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

func (t *Track) Bind(b *TrackBinding) {
	t.binder.SetCallback(t.updateFromData)
	t.binder.Bind(b)
}

func (t *Track) Unbind() {
	t.binder.Unbind()
}

func (t *Track) updateFromData(b binding.DataItem) {
	trackBind := b.(*TrackBinding)
	scTrack := *trackBind.track

	fmt.Println("updateFromData", scTrack)

	t.TrackInfo.Update(scTrack)
}
