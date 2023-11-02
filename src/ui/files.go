package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/helpers"
)

/*
TODO:
Make this more generic so it can be used to select dirs too & used outside of settings
*/
func (d *Data) openFileCanvas(w fyne.Window, label string, updateVal *string, fileFilter []string) fyne.CanvasObject {

	infoLabel := widget.NewLabel(label)
	pathLabel := widget.NewLabel(*updateVal)

	buttonWidget := widget.NewButton("Open", func() {
		f := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialogErr(w, err)
				return
			}
			if reader == nil {
				helpers.WriteToLog("no writer")
				return
			}
			*updateVal = reader.URI().Path()
			pathLabel.SetText(*updateVal)
		}, w)

		fileURI := storage.NewFileURI(*updateVal)
		dirURI, err := storage.Parent(fileURI)
		if err != nil {
			dialogErr(w, err)
			return
		}
		dirListableURI, err := storage.ListerForURI(dirURI)
		if err != nil {
			dialogErr(w, err)
			return
		}
		f.SetLocation(dirListableURI)
		f.SetFilter(storage.NewExtensionFileFilter(fileFilter))
		f.Show()
	})
	// can use this to set a file open icon
	// buttonWidget.SetIcon()
	// or just call newButtonWithIcon

	return formatOpenCanvas(infoLabel, pathLabel, buttonWidget)
}

func (d *Data) openDirCanvas(w fyne.Window, label string, updateVal *string) fyne.CanvasObject {

	infoLabel := widget.NewLabel(label)
	pathLabel := widget.NewLabel(*updateVal)

	buttonWidget := widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), func() {
		f := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				dialogErr(w, err)
				return
			}
			if reader == nil {
				helpers.WriteToLog("no writer")
				return
			}
			*updateVal = reader.Path()
			pathLabel.SetText(*updateVal)
		}, w)

		dirURI := storage.NewFileURI(*updateVal)
		dirListableURI, err := storage.ListerForURI(dirURI)
		if err != nil {
			dialogErr(w, err)
			return
		}
		f.SetLocation(dirListableURI)
		f.Show()
	})
	// can use this to set a file open icon
	// buttonWidget.SetIcon()
	// or just call newButtonWithIcon
	// need to get a button icon first though

	return formatOpenCanvas(infoLabel, pathLabel, buttonWidget)
}

func formatOpenCanvas(infoLabel *widget.Label, pathLabel *widget.Label, buttonWidget *widget.Button) fyne.CanvasObject {

	container := container.NewBorder(infoLabel, widget.NewSeparator(), nil, buttonWidget, pathLabel)

	return container
}
