package iwidget

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

type PathType int

const (
	File PathType = iota
	Directory
)

type OpenPath struct {
	widget.BaseWidget

	parentWindow fyne.Window

	Type PathType

	BaseDir string
	URI     fyne.URI

	ExtensionFilter []string

	OnValidCallback func()
	OnErrorCallback func(error)

	PathLabel *ClickableLabel

	Dialog *dialog.FileDialog
}

func NewOpenPath(w fyne.Window, startingPath string, pathType PathType) *OpenPath {

	openPath := &OpenPath{
		PathLabel:    NewClickableLabel(startingPath, func() {}),
		Type:         pathType,
		parentWindow: w,
	}

	if startingPath != "" {
		URI, err := storage.ParseURI(startingPath)
		if err != nil {
			fmt.Println("Failed to parse starting path", err)
		}
		openPath.URI = URI
	}

	openPath.BaseDir = "/"

	openPath.SetDialog()

	openPath.PathLabel.OnTapped = func() {
		listableURI, err := getListableURI(openPath.URI.Path(), openPath.BaseDir)

		if err != nil {
			fmt.Println("Failed to get listable URI", err)
			return
		}

		openPath.Dialog.SetLocation(listableURI)
		openPath.Dialog.Show()
	}

	openPath.ExtendBaseWidget(openPath)

	return openPath

}

func (i *OpenPath) SetBaseDir(dir string) {
	i.BaseDir = dir
}

func (i *OpenPath) SetURI(path string) {
	URI, err := storage.ParseURI(path)
	if err != nil {
		fmt.Println("Failed to parse path", err)
	}
	i.URI = URI
	pathStr := URI.Path()
	i.PathLabel.SetText(pathStr)
}

func (i *OpenPath) SetExtensionFilter(filter []string) {
	i.ExtensionFilter = filter
}

func (i *OpenPath) SetOnValidCallback(callback func()) {
	i.OnValidCallback = callback
}

func (i *OpenPath) SetOnErrorCallback(callback func(error)) {
	i.OnErrorCallback = callback
}

func (i *OpenPath) CreateRenderer() fyne.WidgetRenderer {

	c := container.NewBorder(
		nil, nil,
		widget.NewIcon(theme.FolderOpenIcon()),
		nil,
		i.PathLabel,
	)

	return widget.NewSimpleRenderer(c)
}

func (i *OpenPath) SetDialog() {

	switch i.Type {
	case File:
		i.Dialog = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				i.OnErrorCallback(err)
				return
			}
			if reader == nil {
				return
			}
			// Below runs if file selection was valid
			i.URI = reader.URI()
			pathStr := reader.URI().Path()
			i.PathLabel.SetText(pathStr)
			i.OnValidCallback()
		}, i.parentWindow)
		i.Dialog.SetFilter(storage.NewExtensionFileFilter(i.ExtensionFilter))
	case Directory:
		i.Dialog = dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				i.OnErrorCallback(err)
				return
			}
			if reader == nil {
				return
			}
			// Below runs if directory selection was valid
			i.URI = reader
			pathStr := reader.Path()
			i.PathLabel.SetText(pathStr)
			i.OnValidCallback()
		}, i.parentWindow)
	}
}

func getListableURI(path string, baseDir string) (fyne.ListableURI, error) {

	var recursionCount int
	dirPath, err := helpers.GetClosestDir(path, baseDir, &recursionCount)
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
