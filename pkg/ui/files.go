package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
)

func (d *Data) openFileCanvas(w fyne.Window, title string, updateVal *string, fileFilter []string, callbackFn func()) fyne.CanvasObject {

	pathCard := buildPathCard(*updateVal, "file")

	buttonWidget := widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), func() {
		if d.State.processing {
			showErrorDialog(w, helpers.ErrPleaseWaitForProcess)
			return
		}

		f := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				showErrorDialog(w, err)
				return
			}
			if reader == nil {
				return
			}
			// Below runs if file selection was valid
			*updateVal = reader.URI().Path()
			pathCard.SetSubTitle(*updateVal)
			callbackFn()
		}, w)
		// Set properties of the file open dialog
		location, err := d.getListableURI(*updateVal)
		if err != nil {
			showErrorDialog(w, err)
			return
		}
		f.SetLocation(location)
		f.Resize(fyne.NewSize(640, 480))
		f.SetFilter(storage.NewExtensionFileFilter(fileFilter))
		f.Show()
	})

	return formatOpenCanvas(title, pathCard, buttonWidget)
}

func (d *Data) openDirCanvas(w fyne.Window, title string, updateVal *string, callbackFn func()) fyne.CanvasObject {

	pathCard := buildPathCard(*updateVal, "directory")

	buttonWidget := widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), func() {

		if d.State.processing {
			showErrorDialog(w, helpers.ErrPleaseWaitForProcess)
			return
		}

		f := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				showErrorDialog(w, err)
				return
			}
			if reader == nil {
				return
			}
			// Below runs if directory selection was valid
			*updateVal = reader.Path()
			pathCard.SetSubTitle(*updateVal)
			callbackFn()
		}, w)
		// Set properties of the folder open dialog
		location, err := d.getListableURI(*updateVal)
		if err != nil {
			showErrorDialog(w, err)
			return
		}
		f.SetLocation(location)
		f.Resize(fyne.NewSize(640, 480))
		f.Show()
	})

	return formatOpenCanvas(title, pathCard, buttonWidget)
}

/*
Builds a fyne card given a path and pathType

pathType is used to determine the default card text if path is empty
*/
func buildPathCard(path string, pathType string) *widget.Card {

	var cardText string

	if path == "" {
		cardText = fmt.Sprintf("Please select a valid %s", pathType)
	} else {
		cardText = path
	}

	pathCard := widget.NewCard("", cardText, nil)
	return pathCard
}

/*
Accepts a path and returns a listable URI for the closest directory

If no 'close directory' is found, it will return the base directory
If the base directory is not found, it will return the root directory (i.e. /)
*/
func (d *Data) getListableURI(path string) (fyne.ListableURI, error) {

	var recursionCount int
	dirPath, err := helpers.GetClosestDir(path, d.Config.BaseDir, &recursionCount)
	if err != nil {
		return nil, helpers.GenErrGettingClosestDir(err)
	}
	dirURI := storage.NewFileURI(dirPath)
	dirListableURI, err := storage.ListerForURI(dirURI)
	if err != nil {
		return nil, helpers.GenErrGettingListableURI(err)
	}
	return dirListableURI, nil
}

/*
Creates a fyne container to be used for opening files/directories
*/
func formatOpenCanvas(title string, pathLabel fyne.CanvasObject, buttonWidget fyne.CanvasObject) fyne.CanvasObject {

	titleCard := widget.NewCard(title, "", nil)
	sep := widget.NewSeparator()

	container := container.NewBorder(titleCard, sep, nil, buttonWidget, pathLabel)

	return container
}
