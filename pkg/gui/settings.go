package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
)

func (e *guiEnv) openSettingsWindow(a fyne.App) bool {

	if e.settingsAlreadyOpen {
		return true
	} else {
		e.settingsAlreadyOpen = true
		// clone config state so we can discard changes if the user closes the window without saving
		tmpConfig := *e.Config
		e.tmpConfig = &tmpConfig
	}

	w := a.NewWindow("Settings")

	w.SetOnClosed(func() {
		e.settingsAlreadyOpen = false
		e.tmpConfig = nil
	})

	// Create a new container with a vertical layout
	container := container.NewVScroll(container.NewVBox(
		e.settingsList(w)...,
	))

	// Set the window content to the container
	w.SetContent(container)

	w.Resize(fyne.NewSize(800, 600))

	w.Show()

	return false
}

/*
settingsList generates a list of canvas objects for the settings window

any altered settings are stored in the TmpConfig struct, which is discarded if the user closes the window without saving
*/
func (e *guiEnv) settingsList(w fyne.Window) []fyne.CanvasObject {

	objs := []fyne.CanvasObject{}

	objs = append(objs, e.openFileCanvas(
		"Traktor Collection Filepath", &e.tmpConfig.TraktorCollectionPath, []string{".nml"}, func() {},
	))
	objs = append(objs, e.openDirCanvas(
		"Temporary Content Directory", &e.tmpConfig.TmpDir, func() {},
	))

	objs = append(objs, e.saveButton(w))

	return objs

}

/*
saveButton returns a button that saves the current state of the TmpConfig struct to the Config struct
and then saves the Config struct to the config file
*/
func (e *guiEnv) saveButton(w fyne.Window) *widget.Button {
	btn := widget.NewButton("Save", func() {
		if e.guiState.busy {
			e.showErrorDialog(helpers.ErrPleaseWaitForProcess)
			return
		}

		e.Config = e.tmpConfig
		err := e.Config.SaveConfig()
		if err != nil {
			e.showErrorDialog(err)
			return
		}
		dialog.ShowInformation("Settings", "Settings saved", w)
		w.Close()
	})
	btn.Importance = widget.HighImportance
	return btn
}
