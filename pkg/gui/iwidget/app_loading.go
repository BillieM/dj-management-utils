package iwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type AppLoading struct {
	widget.BaseWidget
}

func NewAppLoading() *AppLoading {

	appLoading := &AppLoading{}

	appLoading.ExtendBaseWidget(appLoading)

	return appLoading

}

func (t *AppLoading) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewVBox(
			widget.NewLabel("Application is starting, please wait..."),
			widget.NewProgressBarInfinite(),
		),
	)
}
