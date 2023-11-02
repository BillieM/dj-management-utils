package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func setMainContent(w fyne.Window, contentStack *fyne.Container, operation Operation) {

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

func separateTrackView(w fyne.Window) fyne.CanvasObject {
	return widget.NewLabel("separateTrackView")
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
