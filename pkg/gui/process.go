package gui

import (
	"context"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

/*
stepHandler Implements methods defined in the StepHandler interface defined in operations/operations.go

The callbacks are used to update the UI as the process runs
*/
type convertStepHandler struct {
	ctxClose             context.CancelCauseFunc
	bindVals             *progressBindingList
	progressBarBindValue binding.Float
	stepFunc             func()
	finishedFunc         func()
}

/*
StepCallback is executed each time a step finishes inside the operations package
*/
func (o convertStepHandler) StepCallback(stepInfo operations.StepInfo) {
	if stepInfo.Progress != 0 {
		o.progressBarBindValue.Set(stepInfo.Progress)
	}
	if !stepInfo.SkipLog {
		o.bindVals.Append(&progressBindingItem{
			id:         len(o.bindVals.Items),
			message:    stepInfo.Message,
			err:        stepInfo.Error,
			importance: stepInfo.Importance.GetFyneImportance(),
		})
	}
	o.stepFunc()
}

/*
ExitCallback is executed when the process finishes in the operations package
*/
func (o convertStepHandler) ExitCallback() {
	o.finishedFunc()
}

/*
Custom bindings for the log to allow for appending log messages with additional information
*/
type bindBase struct {
	sync.RWMutex
	listeners sync.Map // map[DataListener]bool
}

type progressBindingItem struct {
	bindBase

	id         int
	message    string
	err        error
	importance widget.Importance
}

type progressBindingList struct {
	bindBase

	Items []*progressBindingItem
}

func (i *progressBindingItem) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *progressBindingItem) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (i *progressBindingList) AddListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Store(l, true)
}

func (i *progressBindingList) RemoveListener(l binding.DataListener) {
	i.Lock()
	defer i.Unlock()
	i.listeners.Delete(l)
}

func (i *progressBindingList) GetItem(index int) (binding.DataItem, error) {
	i.Lock()
	defer i.Unlock()
	if index < 0 || index >= len(i.Items) {
		return nil, helpers.ErrIndexOutOfBounds
	}

	return i.Items[index], nil
}

func (i *progressBindingList) Length() int {
	i.Lock()
	defer i.Unlock()
	return len(i.Items)
}

func (i *progressBindingList) Append(message *progressBindingItem) {
	i.Lock()
	defer i.Unlock()
	i.Items = append(i.Items, message)
}

/*
stepsContainer is used to store widgets associated with a running process

stepInfo will be displayed here through the log and progress bar widgets
*/
type stepsContainer struct {
	container   *fyne.Container
	stopButton  *widget.Button
	progressBar *widget.ProgressBar
	log         *widget.List
}

/*
buildProcessContainer builds the processContainer widget,
which is used to display the progress of a running process and allow the user to stop it
*/
func buildProcessContainer(cancelCauseFunc context.CancelCauseFunc, progressBarBindVal binding.Float, logBindVals *progressBindingList) stepsContainer {
	processContainerTop := container.NewVBox()
	progressBar := widget.NewProgressBarWithData(progressBarBindVal)
	stopButton := widget.NewButton("Stop", func() {
		cancelCauseFunc(helpers.ErrUserStoppedProcess)
	})
	var log *widget.List
	log = widget.NewListWithData(logBindVals,
		func() fyne.CanvasObject {
			// Template function for the list, called when the list is created
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			// Update function for the list, called when the bind value changes
			label := o.(*widget.Label)

			id := i.(*progressBindingItem).id
			msg := i.(*progressBindingItem).message
			err := i.(*progressBindingItem).err

			var labelText string
			if err != nil {
				labelText = err.Error()
			} else {
				labelText = msg
			}

			label.Bind(binding.BindString(&labelText))
			label.Importance = i.(*progressBindingItem).importance
			log.SetItemHeight(
				id,
				label.MinSize().Height,
			)
		},
	)

	processContainerTop.Add(progressBar)
	processContainerTop.Add(stopButton)

	processContainer := container.NewBorder(
		container.NewBorder(nil, nil, nil, stopButton, progressBar),
		nil, nil, nil,
		log,
	)

	return stepsContainer{
		container:   processContainer,
		progressBar: progressBar,
		log:         log,
		stopButton:  stopButton,
	}
}
