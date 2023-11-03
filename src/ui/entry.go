package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"

	"github.com/billiem/seren-management/src/helpers"
)

type Data struct {
	*helpers.Config
	*State
	TmpConfig      *helpers.Config
	Operations     map[string]Operation
	OperationIndex map[string][]string
}

type State struct {
	settingsAlreadyOpen bool
	processing          bool
}

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

/*
list
	fyne has a list type, this could be useful for displaying a list of tracks ?

tree
	fyne has a tree type, assuming i can make selections here, this could be useful for selecting folders to convert to stems ?

table
	fyne has a table type, this could be useful for displaying a list of tracks, kinda worried about performance with a very large libary though
	perhaps it'll be alright if we're just dealing with tracks within playlists ?

*/
