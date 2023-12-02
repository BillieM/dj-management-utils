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

	onValidCallback func(string)
	onErrorCallback func(error)
	onOpenCallback  func()
	onCloseCallback func()

	PathCard *ClickablePathCard

	Dialog *dialog.FileDialog
}

func NewOpenPath(w fyne.Window, startingPath string, pathType PathType) *OpenPath {

	openPath := &OpenPath{
		PathCard:     NewClickablePathCard("", theme.FolderOpenIcon(), func() {}),
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

	openPath.setDialog()

	openPath.PathCard.OnTapped = func() {

		func() {
			fmt.Println("Opening dialog")
			if openPath.onOpenCallback != nil {
				openPath.onOpenCallback()
			}
		}()

		defer func() {
			fmt.Println("Closing dialog")
			if openPath.onCloseCallback != nil {
				openPath.onCloseCallback()
			}
		}()

		listableURI, err := getListableURI(openPath.URI.Path(), openPath.BaseDir)

		if err != nil {
			fmt.Println("Failed to get listable URI", err)
			return
		}

		if openPath.Type == File && len(openPath.ExtensionFilter) > 0 {
			openPath.SetExtensionFilter(openPath.ExtensionFilter)
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

/*
Sets the path of the OpenPath widget from a string

Automatically adds the file:// prefix if it is not present
*/
func (i *OpenPath) SetURIFromPathString(uriString string) {
	URI, err := storage.ParseURI(fmt.Sprintf("file://%s", uriString))
	if err != nil {
		fmt.Println("Failed to parse path", err)
		return
	}
	i.SetURI(URI)
}

func (i *OpenPath) SetURI(uri fyne.URI) {
	i.URI = uri
	pathStr := uri.Path()
	i.PathCard.SetText(pathStr)
}

func (i *OpenPath) SetExtensionFilter(filter []string) {
	i.ExtensionFilter = filter
}

/*
Callback for when a valid file or directory is selected

The path is passed to the callback as a string
*/

func (i *OpenPath) SetOnValidCallback(callback func(string)) {
	i.onValidCallback = callback
}

/*
Callback for when an error occurs

The error is passed to the callback
*/
func (i *OpenPath) SetOnErrorCallback(callback func(error)) {
	i.onErrorCallback = callback
}

func (i *OpenPath) SetOnOpenCallback(callback func()) {
	i.onOpenCallback = callback
}

func (i *OpenPath) SetOnCloseCallback(callback func()) {
	i.onCloseCallback = callback
}

func (i *OpenPath) CreateRenderer() fyne.WidgetRenderer {

	c := i.PathCard

	return widget.NewSimpleRenderer(c)
}

func (i *OpenPath) setDialog() {

	switch i.Type {
	case File:
		i.Dialog = dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				i.onErrorCallback(err)
				return
			}
			if reader == nil {
				return
			}
			// Below runs if file selection was valid
			i.SetURI(reader.URI())
			if i.onValidCallback != nil {
				i.onValidCallback(reader.URI().Path())
			}
		}, i.parentWindow)
	case Directory:
		i.Dialog = dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				i.onErrorCallback(err)
				return
			}
			if reader == nil {
				return
			}
			// Below runs if directory selection was valid
			i.SetURI(reader)
			if i.onValidCallback != nil {
				i.onValidCallback(reader.Path())
			}
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

type ClickablePathCard struct {
	ClickableCard

	PathLabel *widget.Label
	Icon      *widget.Icon
}

func NewClickablePathCard(text string, icon fyne.Resource, onTapped func()) *ClickablePathCard {

	clickablePathCard := &ClickablePathCard{
		PathLabel: widget.NewLabel(text),
		Icon:      widget.NewIcon(icon),
	}

	clickablePathCard.ExtendBaseWidget(clickablePathCard)

	clickablePathCard.OnTapped = onTapped

	clickablePathCard.SetContent(container.NewBorder(
		nil, nil,
		clickablePathCard.Icon,
		nil,
		clickablePathCard.PathLabel,
	))

	return clickablePathCard
}

func (i *ClickablePathCard) SetText(text string) {
	i.PathLabel.SetText(text)

	i.SetContent(container.NewBorder(
		nil, nil,
		i.Icon,
		nil,
		i.PathLabel,
	))

}
