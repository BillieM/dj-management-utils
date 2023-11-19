package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/billiem/seren-management/pkg/helpers"
)

func Entry() {

	c, err := helpers.LoadConfig()

	if err != nil {
		helpers.HandleFatalError(err)
	}

	d := buildData(c)

	a := app.New()
	w := a.NewWindow("Library Utilities")

	// Seems strange this method is called SetMainMenu as it really defines the top bar of the application, but hey :)
	w.SetMainMenu(d.makeNavBar(a, w))

	w.Resize(fyne.NewSize(960, 720))

	contentStack := container.NewStack()
	d.setMainContent(w, contentStack, d.getOperationsList()["home"])

	split := container.NewHSplit(d.makeNavMenu(w, contentStack), contentStack)
	split.Offset = 0.25

	w.SetContent(split)
	w.ShowAndRun()
}
