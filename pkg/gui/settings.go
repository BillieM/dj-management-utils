package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func (e *guiEnv) openSettingsWindow(a fyne.App) bool {

	if e.busy {
		return true
	} else {
		e.busy = true
		// clone config state so we can discard changes if the user closes the window without saving
		tmpConfig := *e.Config
		e.tmpConfig = &tmpConfig
	}

	w := a.NewWindow("Settings")

	w.SetOnClosed(func() {
		e.busy = false
		e.tmpConfig = nil
	})

	// Create a new container
	// Use bordered layout to show save button at the bottom &
	// tabs at the top
	tabsContainer := container.NewAppTabs(
		container.NewTabItem("General", e.generalTab()),
		container.NewTabItem("Stems", e.stemsTab()),
		container.NewTabItem("SoundCloud", e.soundCloudTab()),
		container.NewTabItem("Traktor", e.traktorTab()),
		container.NewTabItem("Rekordbox", e.rekordboxTab()),
	)

	container := container.NewBorder(
		nil,
		container.NewBorder(
			nil, nil, nil,
			e.saveButton(w),
			widget.NewLabel("Save settings"),
		),
		nil, nil,

		tabsContainer,
	)

	// Set the window content to the container
	w.SetContent(container)

	w.Resize(fyne.NewSize(800, 600))

	w.Show()

	return false
}

func (e *guiEnv) generalTab() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("General settings"),
	)
}

func (e *guiEnv) stemsTab() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Stems settings"),
	)
}

func (e *guiEnv) soundCloudTab() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("SoundCloud settings"),
	)
}

func (e *guiEnv) traktorTab() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Traktor settings"),
	)
}

func (e *guiEnv) rekordboxTab() *fyne.Container {
	return container.NewVBox(
		widget.NewLabel("Rekordbox settings"),
	)
}

/*
saveButton returns a button that saves the current state of the TmpConfig struct to the Config struct
and then saves the Config struct to the config file
*/
func (e *guiEnv) saveButton(w fyne.Window) *widget.Button {
	btn := widget.NewButton("Save", func() {
		if e.isBusy() {
			return
		}

		e.Config = e.tmpConfig
		err := e.Config.SaveConfig()
		if err != nil {
			e.showErrorDialog(err, true)
			return
		}
		dialog.ShowInformation("Settings", "Settings saved", w)
		w.Close()
	})
	btn.Importance = widget.HighImportance
	return btn
}
