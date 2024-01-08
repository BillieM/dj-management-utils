package iwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
AppLoading is a widget that is displayed when the application is starting up, it displays a simple message
and an infinite progress bar

It is used to prevent the user from interacting with the application while it is starting up, and in the future
to provide updates on the loading process of the application to the user

This will be necessary once indexing functionality is added to the application, as the application will need to
index music files on the user's system, as well as checking the database for DJ software applications for changes
*/
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
