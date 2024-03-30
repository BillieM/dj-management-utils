package gui

import (
	"fmt"

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
			widget.NewSeparator(),
			nil, nil,
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
		widget.NewLabel("General Settings"),
	)
}

func (e *guiEnv) stemsTab() *fyne.Container {

	// build input widgets
	batchSizeSlider := widget.NewSlider(1, 10)
	mergeSlider := widget.NewSlider(1, 10)
	cleanUpSlider := widget.NewSlider(1, 10)
	cudaCheckbox := widget.NewCheck("", func(useCuda bool) {
		e.tmpConfig.CudaEnabled = useCuda
	})

	// build form items
	batchSizeFormItem := widget.NewFormItem("", batchSizeSlider)
	mergeFormItem := widget.NewFormItem("", mergeSlider)
	cleanUpFormItem := widget.NewFormItem("", cleanUpSlider)
	cudaFormItem := widget.NewFormItem("Process stems with CUDA", cudaCheckbox)

	// set form item tooltips
	batchSizeFormItem.HintText = "The batch size to call demucs (the stem separation library) with. Higher values may use more memory."
	mergeFormItem.HintText = "The number of workers to use for merging demucs output to m4a."
	cleanUpFormItem.HintText = "The number of workers to use for cleaning up demucs output."
	cudaFormItem.HintText = `Use CUDA for demucs processing. This requires a Nvidia GPU with CUDA support to work. Usage of CUDA will speed up demucs processing significantly.`

	form := widget.NewForm(
		batchSizeFormItem,
		mergeFormItem,
		cleanUpFormItem,
		cudaFormItem,
	)

	// set slider change callback
	batchSizeSlider.OnChanged = func(val float64) {
		e.tmpConfig.DemucsBatchSize = int(val)
		batchSizeFormItem.Text = fmt.Sprintf("Demucs batch size: %d", e.tmpConfig.DemucsBatchSize)
		form.Refresh()
	}

	mergeSlider.OnChanged = func(val float64) {
		e.tmpConfig.MergeWorkers = int(val)
		mergeFormItem.Text = fmt.Sprintf("Merge workers: %d", e.tmpConfig.MergeWorkers)
		form.Refresh()
	}

	cleanUpSlider.OnChanged = func(val float64) {
		e.tmpConfig.CleanUpWorkers = int(val)
		cleanUpFormItem.Text = fmt.Sprintf("Clean up workers: %d", e.tmpConfig.CleanUpWorkers)
		form.Refresh()
	}

	// set form item values
	batchSizeSlider.SetValue(float64(e.tmpConfig.DemucsBatchSize))
	mergeSlider.SetValue(float64(e.tmpConfig.MergeWorkers))
	cleanUpSlider.SetValue(float64(e.tmpConfig.CleanUpWorkers))
	cudaCheckbox.SetChecked(e.tmpConfig.CudaEnabled)

	return container.NewBorder(
		widget.NewLabel("Warning, changing these settings may cause the application to crash or behave unexpectedly"),
		nil, nil, nil,
		form,
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
