package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/helpers"
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

func (d *Data) settingsList(w fyne.Window) []fyne.CanvasObject {

	objs := []fyne.CanvasObject{}

	objs = append(objs, d.openFileCanvas(
		w, "Traktor Collection Path", &d.TmpConfig.TraktorCollectionPath, []string{".nml"},
	))
	objs = append(objs, d.openDirCanvas(
		w, "Tmp directory", &d.TmpConfig.TmpDir,
	))

	objs = append(objs, d.saveButton(w))

	return objs

}

func (d *Data) saveButton(w fyne.Window) *widget.Button {
	btn := widget.NewButton("Save", func() {
		d.Config = d.TmpConfig
		err := d.Config.SaveConfig()
		if err != nil {
			dialogErr(w, err)
			return
		}
	})
	btn.Importance = widget.HighImportance
	return btn
}

/*
TODO:
Make this more generic so it can be used to select dirs too & used outside of settings
*/
func (d *Data) openFileCanvas(w fyne.Window, label string, updateVal *string, fileFilter []string) fyne.CanvasObject {

	infoLabel := widget.NewLabel(label)
	pathLabel := widget.NewLabel(*updateVal)

	buttonWidget := widget.NewButton("Open", func() {
		f := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil {
				dialogErr(w, err)
				return
			}
			if reader == nil {
				helpers.WriteToLog("no writer")
				return
			}
			*updateVal = reader.URI().Path()
			pathLabel.SetText(*updateVal)
		}, w)

		fileURI := storage.NewFileURI(*updateVal)
		dirURI, err := storage.Parent(fileURI)
		if err != nil {
			dialogErr(w, err)
			return
		}
		dirListableURI, err := storage.ListerForURI(dirURI)
		if err != nil {
			dialogErr(w, err)
			return
		}
		f.SetLocation(dirListableURI)
		f.SetFilter(storage.NewExtensionFileFilter(fileFilter))
		f.Show()
	})
	// can use this to set a file open icon
	// buttonWidget.SetIcon()
	// or just call newButtonWithIcon

	return formatOpenCanvas(infoLabel, pathLabel, buttonWidget)
}

func (d *Data) openDirCanvas(w fyne.Window, label string, updateVal *string) fyne.CanvasObject {

	infoLabel := widget.NewLabel(label)
	pathLabel := widget.NewLabel(*updateVal)

	buttonWidget := widget.NewButton("Open", func() {
		f := dialog.NewFolderOpen(func(reader fyne.ListableURI, err error) {
			if err != nil {
				dialogErr(w, err)
				return
			}
			if reader == nil {
				helpers.WriteToLog("no writer")
				return
			}
			*updateVal = reader.Path()
			pathLabel.SetText(*updateVal)
		}, w)

		dirURI := storage.NewFileURI(*updateVal)
		dirListableURI, err := storage.ListerForURI(dirURI)
		if err != nil {
			dialogErr(w, err)
			return
		}
		f.SetLocation(dirListableURI)
		f.Show()
	})
	// can use this to set a file open icon
	// buttonWidget.SetIcon()
	// or just call newButtonWithIcon
	// need to get a button icon first though

	return formatOpenCanvas(infoLabel, pathLabel, buttonWidget)
}

func formatOpenCanvas(infoLabel *widget.Label, pathLabel *widget.Label, buttonWidget *widget.Button) fyne.CanvasObject {

	container := container.NewBorder(infoLabel, nil, nil, buttonWidget, pathLabel)

	return container
}
