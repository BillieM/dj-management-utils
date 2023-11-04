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

/*
builds the main data object for the application

TODO: this may want to be moved out of the ui package
*/
func buildData(c *helpers.Config) *Data {
	d := &Data{c, nil, nil, nil, nil}

	s := &State{}
	operations := d.getOperationsList()
	operationIndex := d.getOperationIndex()

	d.State = s
	d.Operations = operations
	d.OperationIndex = operationIndex

	return d
}
