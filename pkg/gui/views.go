package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/operations"
)

/*
setMainContent sets the main content of the window to the provided content

Called on tab change on the main menu
*/
func (e *guiEnv) setMainContent(w fyne.Window, contentStack *fyne.Container, view guiView) {

	labelContainer := container.NewVBox(widget.NewLabel(view.name), widget.NewSeparator())

	contentContainer := container.NewBorder(labelContainer, nil, nil, nil, view.render(w))

	contentStack.Objects = []fyne.CanvasObject{contentContainer}
	contentStack.Refresh()
}

func (e *guiEnv) homeView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("Welcome to Seren Library Management!")
}

func (e *guiEnv) stemsView(w fyne.Window) fyne.CanvasObject {
	content := widget.NewLabel("Contains a selection of utilities for separating stems from audio files.")

	return container.NewVBox(content)
}

func (e *guiEnv) separateSingleStemView(w fyne.Window) fyne.CanvasObject {
	ok, canvas := e.checkConfig([]func() (bool, string){e.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.SeparateSingleStemOpts{}

	startFunc := func() {
		e.startSeparateSingleStem(w, processContainerOuter, opts)
	}

	startButton := widget.NewButton("Separate stem", startFunc)
	startButton.Disable()

	trackPathCanvas := e.openFileCanvas(w, "Track Path", &opts.InFilePath, []string{".wav", ".mp3"}, func() { enableBtnIfOptsOkay(opts, startButton) })
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

func (e *guiEnv) separateFolderStemView(w fyne.Window) fyne.CanvasObject {
	ok, canvas := e.checkConfig([]func() (bool, string){e.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.SeparateFolderStemOpts{}

	startFunc := func() {
		e.startSeparateFolderStem(w, processContainerOuter, opts)
	}

	startButton := widget.NewButton("Separate folder", startFunc)
	startButton.Disable()

	trackPathCanvas := e.openDirCanvas(w, "Folder Path", &opts.InDirPath, func() { enableBtnIfOptsOkay(opts, startButton) })
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

func (e *guiEnv) separateCollectionStemView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("separateCollectionView")
}

/*
Convert Mp3s Section
*/

// convertMp3sView returns the view for the convert mp3s info section
func (e *guiEnv) convertMp3sView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertMp3sView")
}

// convertSingleMp3View returns the view for the convert single mp3 operation
func (e *guiEnv) convertSingleMp3View(w fyne.Window) fyne.CanvasObject {
	ok, canvas := e.checkConfig([]func() (bool, string){e.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.ConvertSingleMp3Opts{}

	startFunc := func() {
		e.startConvertSingleMp3(w, processContainerOuter, opts)
	}

	startButton := widget.NewButton("Convert to mp3", startFunc)
	startButton.Disable()

	trackPathCanvas := e.openFileCanvas(w, "Track Path", &opts.InFilePath, []string{".wav", ".flac"}, func() { startButton.Enable() })

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
func (e *guiEnv) convertFolderMp3View(w fyne.Window) fyne.CanvasObject {
	ok, canvas := e.checkConfig([]func() (bool, string){e.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.ConvertFolderMp3Opts{}

	startFunc := func() {
		e.startConvertFolderMp3(w, processContainerOuter, opts)
	}

	startButton := widget.NewButton("Convert folder to mp3", startFunc)
	startButton.Disable()

	trackPathCanvas := e.openDirCanvas(w, "Folder Path", &opts.InDirPath, func() { startButton.Enable() })

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
func (e *guiEnv) convertCollectionMp3View(w fyne.Window) fyne.CanvasObject {

	path := widget.NewLabel("collection path: " + e.Config.TraktorCollectionPath)
	btn := widget.NewButton("do collection things :D", func() {
		// collection.ReadCollection(e.Config.TraktorCollectionPath)
	})

	return container.NewBorder(
		nil, nil, nil, nil,
		container.NewVBox(
			path,
			btn,
		),
	)
}

func (e *guiEnv) tagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("tagsView")
}

func (e *guiEnv) rereadTagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("rereadTagsView")
}

func (e *guiEnv) cleanTagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("cleanTagsView")
}

func (e *guiEnv) conversionView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("conversionView")
}

func (e *guiEnv) playlistMatchingView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("playlistMatchingView")
}
