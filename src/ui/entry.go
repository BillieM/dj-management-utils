package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/billiem/seren-management/src/helpers"
)

func Entry(c *helpers.Config) {

	d := buildData(c)

	a := app.New()
	w := a.NewWindow("Seren Library Management")

	w.SetMainMenu(d.makeMenu(a, w))

	w.Resize(fyne.NewSize(960, 720))

	contentStack := container.NewStack()
	d.setMainContent(w, contentStack, d.getOperationsList()["home"])

	split := container.NewHSplit(d.makeNavMenu(w, contentStack), contentStack)
	split.Offset = 0.25

	w.SetContent(split)
	w.ShowAndRun()
}
