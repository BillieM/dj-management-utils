package ui

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/helpers"
)

func (d *Data) openFileCanvas(w fyne.Window, title string, updateVal *string, fileFilter []string, callbackFn func()) fyne.CanvasObject {

	pathCard := buildPathCard(*updateVal, "file")

	buttonWidget := widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), func() {
		f := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialogErr(w, err)
				return
			}
			if reader == nil {
				helpers.WriteToLog("no writer")
				return
			}
			// Below runs if file selection was valid
			*updateVal = reader.URI().Path()
			pathCard.SetSubTitle(*updateVal)
			callbackFn()
		}, w)
		// Set properties of the file open dialog
		f.SetLocation(d.getListableURI(*updateVal))
		f.SetFilter(storage.NewExtensionFileFilter(fileFilter))
		f.Show()
	})

	return formatOpenCanvas(title, pathCard, buttonWidget)
}

func (d *Data) openDirCanvas(w fyne.Window, title string, updateVal *string, callbackFn func()) fyne.CanvasObject {

	pathCard := buildPathCard(*updateVal, "directory")

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
			// Below runs if directory selection was valid
			*updateVal = reader.Path()
			pathCard.SetSubTitle(*updateVal)
			callbackFn()
		}, w)
		// Set properties of the folder open dialog
		f.SetLocation(d.getListableURI(*updateVal))
		f.Show()
	})

	return formatOpenCanvas(title, pathCard, buttonWidget)
}

/*
Builds a fyne card given a path and pathType

pathType is used to determine the default card text if path is empty

used as a way of displaying paths kinda nicely?
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
func (d *Data) getListableURI(path string) fyne.ListableURI {

	var recursionCount int
	dirPath := d.GetClosestDir(path, &recursionCount)
	dirURI := storage.NewFileURI(dirPath)
	dirListableURI, err := storage.ListerForURI(dirURI)
	if err != nil {
		helpers.HandleFatalError(errors.New("Something went wrong getting the listable URI, err: " + err.Error()))
	}
	return dirListableURI
}

/*
Provides a recursive way of finding the closest directory to a given path,
or the base directory if no 'close directory' is found within 4 recursive calls
If no base directory is found, it will return the root directory (i.e. /)
*/
func (d *Data) GetClosestDir(path string, rCnt *int) string {
	*rCnt++
	fi, err := os.Stat(path)
	// fmt.Println(*rCnt, path)
	if err != nil {
		if *rCnt <= 4 {
			return d.GetClosestDir(filepath.Join(path, ".."), rCnt)
		} else if *rCnt == 5 {
			return d.GetClosestDir(d.Config.BaseDir, rCnt)
		} else if *rCnt == 6 {
			return d.GetClosestDir(filepath.Join("/"), rCnt)
		} else {
			helpers.HandleFatalError(errors.New("Something went very wrong getting the cloest dir, err: " + err.Error()))
		}
	}
	if fi.IsDir() {
		return path
	} else {
		return d.GetClosestDir(filepath.Join(path, ".."), rCnt)
	}
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
