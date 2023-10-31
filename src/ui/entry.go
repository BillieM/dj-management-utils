package ui

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/billiem/seren-management/src/helpers"
)

type Data struct {
	*helpers.Config
	*State
	TmpConfig *helpers.Config
}

type State struct {
	settingsAlreadyOpen bool
}

func Entry(c *helpers.Config) {

	s := &State{}

	d := &Data{c, s, nil}

	fmt.Println("Hello, world!")

	a := app.New()
	w := a.NewWindow("Seren Library Management")

	w.SetMainMenu(d.makeMenu(a, w))

	w.Resize(fyne.NewSize(960, 720))

	hello := widget.NewLabel("Hello Fyne!")
	w.SetContent(container.NewVBox(
		hello,
		widget.NewButton("Hi!", func() {
			hello.SetText("Welcome :)")
		}),
	))

	w.ShowAndRun()
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
