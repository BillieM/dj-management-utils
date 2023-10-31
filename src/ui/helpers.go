package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func dialogErr(w fyne.Window, err error) {
	dialog.ShowError(err, w)
	// helpers.HandleFatalError(err)
}
