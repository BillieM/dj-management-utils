package iwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
ViewLoading is a simple widget that is displayed when the application is loading a view,
it displays a message and an infinite progress bar
*/
type ViewLoading struct {
	widget.BaseWidget

	loadingLabel *widget.Label
	loadingBar   *widget.ProgressBarInfinite
}

func NewViewLoading(loadingText string) *ViewLoading {
	i := &ViewLoading{
		loadingLabel: widget.NewLabel(loadingText),
		loadingBar:   widget.NewProgressBarInfinite(),
	}

	i.ExtendBaseWidget(i)

	return i
}

func (i *ViewLoading) CreateRenderer() fyne.WidgetRenderer {

	c := container.NewBorder(
		i.loadingLabel, nil, nil, nil,
		container.NewBorder(i.loadingBar, nil, nil, nil),
	)

	return widget.NewSimpleRenderer(c)
}
