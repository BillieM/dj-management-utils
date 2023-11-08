package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/operations"
)

func (d *Data) setMainContent(w fyne.Window, contentStack *fyne.Container, operation Operation) {

	contentContainer := container.NewVBox(widget.NewLabel(operation.Name), widget.NewSeparator(), operation.View(w))

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
	processContainer := container.NewVBox()
	processContainer.Hidden = true

	trackPathCanvas := d.openDirCanvas(w, "Folder Path", &dirPath, func() { processContainer.Hidden = false })

	processContainer.Add(widget.NewButton("Convert folder to mp3", func() {
		progressBar := d.buildProgressBar()
		d.State.processing = true
		processContainer.Add(widget.NewLabel("Processing..."))
		processContainer.Add(progressBar)

		// new context
		ctx := context.Background()
		// create cancelFunc from context
		ctx, cancelFunc := context.WithCancel(ctx)

		context.AfterFunc(ctx, func() {
			/*
				rather than pass a callback through, we should just be able to do post processing steps here i.e.:
					restoring interface state
					reporting any failures
					reporting progress if early shutdown
						possibly also for cleaning up any residue files on early shutdown ?
			*/

			d.State.processing = false
		})

		stopButton := widget.NewButton("Stop", func() {
			cancelFunc()
		})
		processContainer.Add(stopButton)

		operations.ConvertFolderMp3Params{
			BaseOperationParams: operations.BaseOperationParams{
				Config: d.Config,
				StepCallback: func(value float64) {
					progressBar.updateProgressBar(value)
				},
				Context: ctx,
			},
			InDirPath: dirPath,
		}.ExecuteOperation()
	}))

	return container.NewVBox(trackPathCanvas, processContainer)
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
