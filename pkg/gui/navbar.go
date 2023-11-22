package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func (e *guiEnv) makeNavBar(a fyne.App, w fyne.Window) *fyne.MainMenu {

	fileMenu := e.makeFileNav(a, w)
	helpMenu := e.makeHelpNav(a, w)

	return fyne.NewMainMenu(
		fileMenu,
		helpMenu,
	)

}

func (e *guiEnv) makeFileNav(a fyne.App, w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("File",
		fyne.NewMenuItem("Settings", func() {
			alreadyOpen := e.openSettingsWindow(a)
			if alreadyOpen {
				dialog.ShowInformation("Settings", "Settings window is already open", w)
			}
		}),
		fyne.NewMenuItem("Quit", func() { a.Quit() }),
	)
}

func (e *guiEnv) makeHelpNav(a fyne.App, w fyne.Window) *fyne.Menu {
	return fyne.NewMenu("Help",
		fyne.NewMenuItem("About", func() {}),
	)
}
