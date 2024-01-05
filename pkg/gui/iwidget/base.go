package iwidget

import (
	"fyne.io/fyne/v2"
	"github.com/billiem/seren-management/pkg/gui/uihelpers"
	"github.com/billiem/seren-management/pkg/helpers"
)

/*
Base is a struct to be embedded inside of custom widgets, this
allows us to provide a common set of functionality to all widgets,
including logging, app config, and access to the main application window
for dialogs
*/
type Base struct {
	Logger       helpers.SerenLogger
	Config       helpers.Config
	MainWindow   fyne.Window
	App          fyne.App
	ResizeEvents *uihelpers.ResizeEvents
}
