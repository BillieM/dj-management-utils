package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/collection"
	"github.com/billiem/seren-management/src/operations"
)

/*
setMainContent sets the main content of the window to the provided content

Called on tab change on the main menu
*/
func (d *Data) setMainContent(w fyne.Window, contentStack *fyne.Container, operation Operation) {

	labelContainer := container.NewVBox(widget.NewLabel(operation.Name), widget.NewSeparator())

	contentContainer := container.NewBorder(labelContainer, nil, nil, nil, operation.View(w))

	contentStack.Objects = []fyne.CanvasObject{contentContainer}
	contentStack.Refresh()
}

func (d *Data) homeView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("Welcome to Seren Library Management!")
}

func (d *Data) stemsView(w fyne.Window) fyne.CanvasObject {
	content := widget.NewLabel("Contains a selection of utilities for separating stems from audio files.")

	return container.NewVBox(content)
}

func (d *Data) separateSingleStemView(w fyne.Window) fyne.CanvasObject {
	ok, canvas := d.checkConfig([]func() (bool, string){d.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.SeparateSingleStemOpts{}

	startFunc := func() {
		d.startSeparateSingleStem(w, processContainerOuter, opts)
	}

	startButton := widget.NewButton("Separate stem", startFunc)
	startButton.Disable()

	trackPathCanvas := d.openFileCanvas(w, "Track Path", &opts.InFilePath, []string{".wav", ".mp3"}, func() { enableBtnIfOptsOkay(opts, startButton) })
	stemTypeSelect := buildStemTypeSelect(&opts.Type, func() { enableBtnIfOptsOkay(opts, startButton) })

	return container.NewBorder(
		container.NewVBox(
			container.NewVBox(
				trackPathCanvas,
				stemTypeSelect,
			),
			startButton,
		), nil, nil, nil,
		processContainerOuter,
	)
}

func (d *Data) separateFolderStemView(w fyne.Window) fyne.CanvasObject {
	ok, canvas := d.checkConfig([]func() (bool, string){d.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.SeparateFolderStemOpts{}

	startFunc := func() {
		d.startSeparateFolderStem(w, processContainerOuter, opts)
	}

	startButton := widget.NewButton("Separate folder", startFunc)
	startButton.Disable()

	trackPathCanvas := d.openDirCanvas(w, "Folder Path", &opts.InDirPath, func() { enableBtnIfOptsOkay(opts, startButton) })
	stemTypeSelect := buildStemTypeSelect(&opts.Type, func() { enableBtnIfOptsOkay(opts, startButton) })

	return container.NewBorder(
		container.NewVBox(
			container.NewVBox(
				trackPathCanvas,
				stemTypeSelect,
			),
			startButton,
		), nil, nil, nil,
		processContainerOuter,
	)
}

func (d *Data) separateCollectionStemView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("separateCollectionView")
}

/*
Convert Mp3s Section
*/

// convertMp3sView returns the view for the convert mp3s info section
func (d *Data) convertMp3sView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertMp3sView")
}

// convertSingleMp3View returns the view for the convert single mp3 operation
func (d *Data) convertSingleMp3View(w fyne.Window) fyne.CanvasObject {
	ok, canvas := d.checkConfig([]func() (bool, string){d.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.ConvertSingleMp3Opts{}

	startFunc := func() {
		d.startConvertSingleMp3(w, processContainerOuter, opts)
	}

	startButton := widget.NewButton("Convert to mp3", startFunc)
	startButton.Disable()

	trackPathCanvas := d.openFileCanvas(w, "Track Path", &opts.InFilePath, []string{".wav", ".flac"}, func() { startButton.Enable() })

	return container.NewBorder(
		container.NewVBox(
			container.NewVBox(
				trackPathCanvas,
			),
			startButton,
		), nil, nil, nil,
		processContainerOuter,
	)
}

// convertFolderMp3View returns the view for the convert folder mp3 operation
func (d *Data) convertFolderMp3View(w fyne.Window) fyne.CanvasObject {
	ok, canvas := d.checkConfig([]func() (bool, string){d.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.ConvertFolderMp3Opts{}

	startFunc := func() {
		d.startConvertFolderMp3(w, processContainerOuter, opts)
	}

	startButton := widget.NewButton("Convert folder to mp3", startFunc)
	startButton.Disable()

	trackPathCanvas := d.openDirCanvas(w, "Folder Path", &opts.InDirPath, func() { startButton.Enable() })

	return container.NewBorder(
		container.NewVBox(
			container.NewVBox(
				trackPathCanvas,
			),
			startButton,
		), nil, nil, nil,
		processContainerOuter,
	)
}

// convertCollectionMp3View returns the view for the convert collection mp3 operation
func (d *Data) convertCollectionMp3View(w fyne.Window) fyne.CanvasObject {

	path := widget.NewLabel("collection path: " + d.Config.TraktorCollectionPath)
	btn := widget.NewButton("do collection things :D", func() {
		collection.ReadCollection(d.Config.TraktorCollectionPath)
	})

	return container.NewBorder(
		nil, nil, nil, nil,
		container.NewVBox(
			path,
			btn,
		),
	)
}

func (d *Data) tagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("tagsView")
}

func (d *Data) rereadTagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("rereadTagsView")
}

func (d *Data) cleanTagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("cleanTagsView")
}

func (d *Data) conversionView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("conversionView")
}

func (d *Data) playlistMatchingView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("playlistMatchingView")
}
