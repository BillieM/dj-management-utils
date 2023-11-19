package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func (d *Data) makeNavBar(a fyne.App, w fyne.Window) *fyne.MainMenu {

	fileMenu := d.makeFileNav(a, w)
	helpMenu := d.makeHelpNav(a, w)

	return fyne.NewMainMenu(
		fileMenu,
		helpMenu,
	)

}

func (d *Data) makeFileNav(a fyne.App, w fyne.Window) *fyne.Menu {
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

func (d *Data) makeHelpNav(a fyne.App, w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {}),
	)
}
