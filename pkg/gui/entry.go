package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"github.com/billiem/seren-management/pkg/helpers"
)

func Entry() {

	e, err := buildGuiEnv()

	if err != nil {
		helpers.HandleFatalError(err)
	}

	a := app.New()
	w := a.NewWindow("Library Utilities")

	// Seems strange this method is called SetMainMenu as it really defines the top bar of the application, but hey :)
	w.SetMainMenu(e.makeNavBar(a, w))

	w.Resize(fyne.NewSize(960, 720))

	contentStack := container.NewStack()
	e.setMainContent(w, contentStack, e.getViewList()["home"])

	split := container.NewHSplit(e.makeNavMenu(w, contentStack), contentStack)
	split.Offset = 0.25

	w.SetContent(split)
	w.ShowAndRun()
}
