package iwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
main section widget

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
when track has a 'Download' link on SoundCloud

*/

type TrackDownload struct {
	widget.BaseWidget

	TrackDownloadButton   *TrackDownloadButton
	TrackDownloadProgress *widget.ProgressBarInfinite
}

func NewTrackDownload() *TrackDownload {
	i := &TrackDownload{
		TrackDownloadButton:   NewTrackDownloadButton(),
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

type TrackDownloadButton struct {
	widget.Button

	HasDownloadsLeft bool
	onTapped         func()
}

func NewTrackDownloadButton() *TrackDownloadButton {
	i := &TrackDownloadButton{
		Button: *widget.NewButton("Download Track", func() {}),
	}

	i.ExtendBaseWidget(i)

	return i
}
