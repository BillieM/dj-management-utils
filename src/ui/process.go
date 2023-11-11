package ui

import (
	"context"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/helpers"
	"github.com/billiem/seren-management/src/operations"
)

/*
operationProcess Implements methods defined in the OperationsProcess interface defined in operations/operations.go

The callbacks are used to update the UI as the process runs
*/
type operationProcess struct {
	ctxClose             context.CancelCauseFunc
	bindVals             *progressBindingList
	progressBarBindValue binding.Float
	stepFunc             func()
	finishedFunc         func()
}

/*
StepCallback is executed each time a step finishes inside the operations package
*/
func (o operationProcess) StepCallback(stepInfo operations.StepInfo) {
	if stepInfo.Progress != 0 {
		o.progressBarBindValue.Set(stepInfo.Progress)
	}
	o.bindVals.Append(&progressBindingItem{
		message:    stepInfo.Message,
		err:        stepInfo.Error,
		importance: stepInfo.Importance.GetFyneImportance(),
	})
	o.stepFunc()
}

/*
ExitCallback is executed when the process finishes in the operations package
*/
func (o operationProcess) ExitCallback() {
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

func (l *progressBindingList) GetItem(index int) (binding.DataItem, error) {
	if index < 0 || index >= len(l.Items) {
		return nil, helpers.ErrIndexOutOfBounds
	}

	return l.Items[index], nil
}

func (l *progressBindingList) Length() int {
	return len(l.Items)
}

func (l *progressBindingList) Append(message *progressBindingItem) {
	l.Items = append(l.Items, message)
}

/*
operationsProcessContainer is used to store widgets associated with a running process
*/
type operationProcessContainer struct {
	container   *fyne.Container
	stopButton  *widget.Button
	progressBar *widget.ProgressBar
	log         *widget.List
}

/*
buildProcessContainer builds the processContainer widget,
which is used to display the progress of a running process and allow the user to stop it
*/
func buildProcessContainer(cancelCauseFunc context.CancelCauseFunc, progressBarBindVal binding.Float, logBindVals *progressBindingList) operationProcessContainer {
	processContainerTop := container.NewVBox()
	progressBar := widget.NewProgressBarWithData(progressBarBindVal)
	stopButton := widget.NewButton("Stop", func() {
		cancelCauseFunc(helpers.ErrUserStoppedProcess)
	})
	log := widget.NewListWithData(logBindVals,
		func() fyne.CanvasObject {
			// Template function for the list, called when the list is created
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			// Update function for the list, called when the bind value changes
			msg := i.(*progressBindingItem).message
			err := i.(*progressBindingItem).err
			if err != nil {
				errMsg := err.Error()
				o.(*widget.Label).Bind(binding.BindString(&errMsg))
			} else {
				o.(*widget.Label).Bind(binding.BindString(&msg))
			}
			o.(*widget.Label).Importance = i.(*progressBindingItem).importance
		},
	)

	processContainerTop.Add(progressBar)
	processContainerTop.Add(stopButton)

	processContainer := container.NewBorder(
		container.NewBorder(nil, nil, nil, stopButton, progressBar),
		nil, nil, nil,
		log,
	)

	return operationProcessContainer{
		container:   processContainer,
		progressBar: progressBar,
		log:         log,
		stopButton:  stopButton,
	}
}
