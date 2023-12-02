package gui

import (
	"fyne.io/fyne/v2/dialog"
)

func (e *guiEnv) showErrorDialog(err error) {
	dialog.ShowError(err, e.mainWindow)
}

func (e *guiEnv) showInfoDialog(title, message string) {
	dialog.ShowInformation(title, message, e.mainWindow)
}
