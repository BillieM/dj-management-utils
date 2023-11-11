package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func showErrorDialog(w fyne.Window, err error) {
	/*
		TODO: may want to add capitalisation to the error message?

		TODO: may want to add a "more info" button to the dialog
			can use .Unwrap() to get the underlying error ?
	*/

	dialog.ShowError(err, w)
}

func showInfoDialog(w fyne.Window, title, message string) {
	dialog.ShowInformation(title, message, w)
}
