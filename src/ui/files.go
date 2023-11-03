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

/*
TODO:
Make this more generic so it can be used to select dirs too & used outside of settings
*/
func (d *Data) openFileCanvas(w fyne.Window, title string, updateVal *string, fileFilter []string) fyne.CanvasObject {

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
			*updateVal = reader.URI().Path()
			pathCard.SetSubTitle(*updateVal)
		}, w)
		f.SetLocation(d.getListableURI(*updateVal))
		f.SetFilter(storage.NewExtensionFileFilter(fileFilter))
		f.Show()
	})

	return formatOpenCanvas(title, pathCard, buttonWidget)
}

func (d *Data) openDirCanvas(w fyne.Window, title string, updateVal *string) fyne.CanvasObject {

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
			*updateVal = reader.Path()
			pathCard.SetSubTitle(*updateVal)
		}, w)
		f.SetLocation(d.getListableURI(*updateVal))
		f.Show()
	})

	return formatOpenCanvas(title, pathCard, buttonWidget)
}

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
Returns a fyne.ListableURI for a given path

If a directory path is given, returns a fyne.ListableURI for that directory
If a file path is given, returns a fyne.ListableURI for the parent directory

If the path does not exist, returns
*/
func (d *Data) getListableURI(path string) fyne.ListableURI {

	var recursionCount int
	dirPath := d.GetClosestDir(path, &recursionCount)
	fmt.Println(dirPath)
	dirURI := storage.NewFileURI(dirPath)
	fmt.Println(dirURI)
	dirListableURI, err := storage.ListerForURI(dirURI)
	fmt.Println(dirListableURI)
	fmt.Println("exists")
	exists, err2 := storage.Exists(dirURI)
	fmt.Println("exists", exists, err2)
	if err != nil {
		helpers.HandleFatalError(errors.New("Something went wrong getting the listable URI, err: " + err.Error()))
	}
	return dirListableURI
}

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

func formatOpenCanvas(title string, pathLabel fyne.CanvasObject, buttonWidget fyne.CanvasObject) fyne.CanvasObject {

	titleCard := widget.NewCard(title, "", nil)
	sep := widget.NewSeparator()

	container := container.NewBorder(titleCard, sep, nil, buttonWidget, pathLabel)

	return container
}
