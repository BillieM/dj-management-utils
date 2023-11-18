package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func showErrorDialog(w fyne.Window, err error) {
	dialog.ShowError(err, w)
}

func showInfoDialog(w fyne.Window, title, message string) {
	dialog.ShowInformation(title, message, w)
}
