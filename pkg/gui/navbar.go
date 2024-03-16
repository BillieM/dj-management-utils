package gui

import (
	"fyne.io/fyne/v2"
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
			pleaseFinish := e.openSettingsWindow(a)
			if pleaseFinish {
				e.showInfoDialog("Error opening settings!", "Please finish what you're doing first")
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
