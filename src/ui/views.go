package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func (d *Data) setMainContent(w fyne.Window, contentStack *fyne.Container, operation Operation) {

	contentContainer := container.NewVBox(widget.NewLabel(operation.Name), widget.NewSeparator(), operation.View(w))

	contentStack.Objects = []fyne.CanvasObject{contentContainer}
	contentStack.Refresh()
}

func homeView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("Welcome to Seren Library Management!")
}

func stemsView(w fyne.Window) fyne.CanvasObject {
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

	processContainer.Add(widget.NewLabel("Processing..."))

	return container.NewVBox(trackPathCanvas, processContainer)
}

func separateFolderView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("separateFolderView")
}

func separateCollectionView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("separateCollectionView")
}

func convertMp3sView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertMp3sView")
}

func convertSingleMp3View(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertSingleMp3View")
}

func convertFolderMp3View(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertFolderMp3View")
}

func convertCollectionMp3View(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("convertCollectionMp3View")
}

func tagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("tagsView")
}

func rereadTagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("rereadTagsView")
}

func cleanTagsView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("cleanTagsView")
}

func conversionView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("conversionView")
}

func playlistMatchingView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("playlistMatchingView")
}

/*
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
