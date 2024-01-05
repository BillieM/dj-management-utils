package iwidget

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/helpers"
)

type PathType int

const (
	File PathType = iota
	Directory
)

type OpenPath struct {
	*Base
	widget.BaseWidget

	Type PathType

	BaseDir string
	URI     fyne.URI

	ExtensionFilter []string

	onValid func(string)
	onError func(error)
	onOpen  func()
	onClose func()

	PathCard *ClickablePathCard

	Dialog *dialog.FileDialog
}

func NewOpenPath(widgetBase *Base, startingPath string, pathType PathType) *OpenPath {

	openPath := &OpenPath{
		Base:     widgetBase,
		PathCard: NewClickablePathCard(widgetBase, "", theme.FolderOpenIcon(), func() {}),
		Type:     pathType,
	}

	if startingPath != "" {
		URI, err := storage.ParseURI(startingPath)
		if err != nil {
			openPath.Logger.NonFatalError(fault.Wrap(
				err,
				fctx.With(fctx.WithMeta(
					context.Background(),
					"starting_path", startingPath,
				)),
				fmsg.With("error parsing starting path"),
			))
		}
		openPath.URI = URI
	} else {
		openPath.URI = storage.NewFileURI("")
	}

	openPath.SetBaseDir(openPath.Config.BaseDir)

	openPath.setDialog()

	openPath.PathCard.OnTapped = func() {

		listableURI, err := getListableURI(openPath.URI.Path(), openPath.BaseDir)

		if err != nil {
			ctx := fctx.WithMeta(
				context.Background(),
				"path", openPath.URI.Path(),
				"base_dir", openPath.BaseDir,
				"listable_uri_path", listableURI.Path(),
			)

			openPath.Logger.NonFatalError(
				fault.Wrap(
					err,
					fctx.With(ctx),
					fmsg.With("error getting listable URI"),
				),
			)
			return
		}

		func() {
			openPath.Logger.Debug(
				"Opening path dialog",
				"path", openPath.URI.Path(),
				"listable_uri_path", listableURI.Path(),
			)
			if openPath.onOpen != nil {
				openPath.onOpen()
			}
		}()

		defer func() {
			openPath.Logger.Debug(
				"Closing path dialog",
				"path", openPath.URI.Path(),
				"listable_uri_path", listableURI.Path(),
			)
			if openPath.onClose != nil {
				openPath.onClose()
			}
		}()

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
SetURIFromPathString sets the path of the OpenPath widget from a string

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
SetOnValid sets the callback for when a valid file or directory is selected

The path is passed to the callback as a string
*/

func (i *OpenPath) SetOnValid(callback func(string)) {
	i.onValid = callback
}

/*
SetOnError sets the callback for when an error occurs

The error is passed to the callback
*/
func (i *OpenPath) SetOnError(callback func(error)) {
	i.onError = callback
}

func (i *OpenPath) SetOnOpen(callback func()) {
	i.onOpen = callback
}

func (i *OpenPath) SetOnClose(callback func()) {
	i.onClose = callback
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
				i.onError(err)
				return
			}
			if reader == nil {
				return
			}
			// Below runs if file selection was valid
			i.SetURI(reader.URI())
			if i.onValid != nil {
				i.onValid(reader.URI().Path())
			}
		}, i.MainWindow)
	case Directory:
		i.Dialog = dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				i.onError(err)
				return
			}
			if reader == nil {
				return
			}
			// Below runs if directory selection was valid
			i.SetURI(reader)
			if i.onValid != nil {
				i.onValid(reader.Path())
			}
		}, i.MainWindow)
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
	*Base
	ClickableCard

	PathLabel *widget.Label
	Icon      *widget.Icon
}

func NewClickablePathCard(widgetBase *Base, text string, icon fyne.Resource, onTapped func()) *ClickablePathCard {

	clickablePathCard := &ClickablePathCard{
		Base:      widgetBase,
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
