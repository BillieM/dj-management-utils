package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

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

func (d *Data) separateTrackView(w fyne.Window) fyne.CanvasObject {

	ok, canvas := d.checkConfig([]func() (bool, string){d.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	var trackPath string
	processContainer := container.NewVBox()
	processContainer.Hidden = true

	trackPathCanvas := d.openFileCanvas(w, "Track Path", &trackPath, []string{".mp3"}, func() { processContainer.Hidden = false })

	processContainer.Add(widget.NewButton("Separate Track", func() {
		progressBar := widget.NewProgressBar()
		d.processing = true
		processContainer.Add(widget.NewLabel("Processing..."))
		processContainer.Add(progressBar)

	}))

	return container.NewVBox(trackPathCanvas, processContainer)
}

func (d *Data) separateFolderView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("separateFolderView")
}

func (d *Data) separateCollectionView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("separateCollectionView")
}

func (d *Data) convertMp3sView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertMp3sView")
}

func (d *Data) convertSingleMp3View(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertSingleMp3View")
}

func (d *Data) convertFolderMp3View(w fyne.Window) fyne.CanvasObject {

	ok, canvas := d.checkConfig([]func() (bool, string){d.Config.CheckTmpDir})

	if !ok {
		return canvas
	}

	var dirPath string
	optionsContainer := container.NewVBox()
	optionsContainer.Hide()

	// show the options container when a valid dir path is selected
	trackPathCanvas := d.openDirCanvas(w, "Folder Path", &dirPath, func() { optionsContainer.Show() })

	processContainerOuter := container.NewStack()

	startFunc := func() {
		d.startConvertFolderMp3(w, processContainerOuter, startConvertFolderMp3Options{dirPath: &dirPath})
	}
	startButton := widget.NewButton("Convert folder to mp3", startFunc)
	optionsContainer.Add(startButton)

	return container.NewBorder(
		container.NewVBox(trackPathCanvas, optionsContainer), nil, nil, nil,
		processContainerOuter,
	)
}

func (d *Data) convertCollectionMp3View(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertCollectionMp3View")
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

/*
Checks the config for any issues for a given set of checks

# Returns true if there are no issues, false if there are issues

If there are issues, it will return a fyne.CanvasObject containing the issues

	TODO: abstract the checking into a seperate function
		can then add unit tests surrounding it
		and create a seperate function for the generation of the canvas object
*/
func (d *Data) checkConfig(checks []func() (bool, string)) (bool, fyne.CanvasObject) {

	configIssues := []string{}

	for _, check := range checks {
		pass, msg := check()
		if !pass {
			configIssues = append(configIssues, msg)
		}
	}

	if len(configIssues) > 0 {
		issuesContainer := container.NewVBox(
			widget.NewLabel("Please fix the following issues with your config:"),
		)
		for _, issue := range configIssues {
			issuesContainer.Add(widget.NewLabel(issue))
		}
		return false, issuesContainer
	}

	return true, nil
}
