package ui

import (
	"context"
	"errors"
	"sync"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/operations"
)

/*
Operation Process

Implements methods defined in the OperationsProcess interface defined in operations/operations.go

The callbacks are used to update the UI as the process runs
*/

type OperationProcess struct {
	ctxClose             context.CancelCauseFunc
	bindVals             *progressBindingList
	progressBarBindValue binding.Float
	stepFunc             func()
	finishedFunc         func()
}

func (o OperationProcess) StepCallback(stepInfo operations.StepInfo) {
	o.progressBarBindValue.Set(stepInfo.Progress)
	o.bindVals.Append(&progressBindingItem{
		message: stepInfo.Message,
	})
	o.stepFunc()
}

func (o OperationProcess) ExitCallback() {
	o.finishedFunc()
}

/*
Custom bindings for the log to allow for appending log messages with additional information

TODO: can probably make these private
*/

type bindBase struct {
	sync.RWMutex
	listeners sync.Map // map[DataListener]bool
}

type progressBindingItem struct {
	bindBase

	message string
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
		return nil, errors.New("index out of bounds")
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
ProcessContainer is used to store widgets associated with a running process

The container will be added to the main content stack, and removed when the process is finished
*/

type operationProcessContainer struct {
	container   *fyne.Container
	progressBar *widget.ProgressBar
	log         *widget.List
}

/*
TODO: rename all the list calls, should be called 'log' or something

Builds the processContainer widget, which is used to display the progress of a running process
*/
func buildProcessContainer(cancelCauseFunc context.CancelCauseFunc, progressBarBindVal binding.Float, logBindVals *progressBindingList) operationProcessContainer {
	processContainerTop := container.NewVBox()
	progressBar := widget.NewProgressBarWithData(progressBarBindVal)
	stopButton := widget.NewButton("Stop", func() {
		cancelCauseFunc(errors.New("user stopped process"))
	})
	log := widget.NewListWithData(logBindVals,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			msg := i.(*progressBindingItem).message
			o.(*widget.Label).Bind(binding.BindString(&msg))
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
	}
}
