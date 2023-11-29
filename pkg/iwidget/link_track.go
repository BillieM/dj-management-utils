package iwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

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
