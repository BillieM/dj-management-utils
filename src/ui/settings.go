package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (d *Data) openSettingsWindow(a fyne.App) bool {

	if d.settingsAlreadyOpen {
		return true
	} else {
		d.settingsAlreadyOpen = true
		// clone config state so we can discard changes if the user closes the window without saving
		tmpConfig := *d.Config
		d.TmpConfig = &tmpConfig
	}

	w := a.NewWindow("Settings")

	w.SetOnClosed(func() {
		d.settingsAlreadyOpen = false
		d.TmpConfig = nil
	})

	// Create a new container with a vertical layout
	container := container.NewVScroll(container.NewVBox(
		d.settingsList(w)...,
	))

	// Set the window content to the container
	w.SetContent(container)

	w.Resize(fyne.NewSize(800, 600))

	w.Show()

	return false
}

/*
TODO: implement 'tabs' for settings, will likely require some refactoring

settingsList generates a list of canvas objects for the settings window

any altered settings are stored in the TmpConfig struct, which is discarded if the user closes the window without saving
*/
func (d *Data) settingsList(w fyne.Window) []fyne.CanvasObject {

	objs := []fyne.CanvasObject{}

	// TODO: may want to add a settings wrapper around these canvas calls in order to add a label to each call?

	objs = append(objs, d.openFileCanvas(
		w, "Traktor Collection Filepath", &d.TmpConfig.TraktorCollectionPath, []string{".nml"}, func() {},
	))
	objs = append(objs, d.openDirCanvas(
		w, "Temporary Content Directory", &d.TmpConfig.TmpDir, func() {},
	))

	objs = append(objs, d.saveButton(w))

	return objs

}

/*
saveButton returns a button that saves the current state of the TmpConfig struct to the Config struct
and then saves the Config struct to the config file
*/
func (d *Data) saveButton(w fyne.Window) *widget.Button {
	btn := widget.NewButton("Save", func() {
		if d.State.processing {
			pleaseWaitForProcess(w)
			return
		}

		d.Config = d.TmpConfig
		err := d.Config.SaveConfig()
		if err != nil {
			dialogErr(w, err)
			return
		}
		dialog.ShowInformation("Settings", "Settings saved", w)
		w.Close()
	})
	btn.Importance = widget.HighImportance
	return btn
}
