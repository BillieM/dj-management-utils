package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func (d *Data) makeMenu(a fyne.App, w fyne.Window) *fyne.MainMenu {

	fileMenu := d.makeFileMenu(a, w)
	helpMenu := d.makeHelpMenu(a, w)

	return fyne.NewMainMenu(
		fileMenu,
		helpMenu,
	)

}

func (d *Data) makeFileMenu(a fyne.App, w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("File",
		fyne.NewMenuItem("Settings", func() {
			alreadyOpen := d.openSettingsWindow(a)
			if alreadyOpen {
				dialog.ShowInformation("Settings", "Settings window is already open", w)
			}
		}),
		fyne.NewMenuItem("Quit", func() { a.Quit() }),
	)
}

func (d *Data) makeHelpMenu(a fyne.App, w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {}),
	)
}
