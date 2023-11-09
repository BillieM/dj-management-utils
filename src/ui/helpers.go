package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func dialogErr(w fyne.Window, err error) {
	dialog.ShowError(err, w)
}

func dialogInfo(w fyne.Window, title, message string) {
	dialog.ShowInformation(title, message, w)
}

func pleaseWaitForProcess(w fyne.Window) {
	dialogInfo(w, "Please Wait", "Please wait for the current process to finish")
}
