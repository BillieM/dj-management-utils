package ui

import (
	"context"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
)

/*
Operation Process

Implements methods defined in the OperationsProcess interface defined in operations/operations.go

The callbacks are used to update the UI as the process runs
*/

type OperationProcess struct {
	ctxClose             context.CancelCauseFunc
	listBindValue        binding.StringList
	progressBarBindValue binding.Float
	stepFunc             func()
	finishedFunc         func()
}

func (o OperationProcess) StepCallback(progress float64, message string) {
	fmt.Println("hit StepCallback() in process.go", progress, message)
	o.progressBarBindValue.Set(progress)
	o.listBindValue.Append(message)
	o.stepFunc()
}

func (o OperationProcess) ExitCallback() {
	fmt.Println("hit ExitCallback() in process.go")
	o.finishedFunc()
}

/*
ProcessContainer is used to store widgets associated with a running process

The container will be added to the main content stack, and removed when the process is finished
*/

type MyProcessContainer struct {
	Container   *fyne.Container
	StopButton  *widget.Button
	ProgressBar *widget.ProgressBar
	List        *widget.List
}

/*
TODO: rename all the list calls, should be called 'log' or something

Builds the processContainer widget, which is used to display the progress of a running process
*/
func buildProcessContainer(cancelCauseFunc context.CancelCauseFunc, progressBarBindVal binding.Float, listBindVal binding.StringList) MyProcessContainer {
	processContainerTop := container.NewVBox()

	progressBar := widget.NewProgressBarWithData(progressBarBindVal)
	stopButton := widget.NewButton("Stop", func() {
		cancelCauseFunc(errors.New("user stopped process"))
	})
	list := widget.NewListWithData(listBindVal,
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i binding.DataItem, o fyne.CanvasObject) {
			o.(*widget.Label).Bind(i.(binding.String))
		},
	)

	processContainerTop.Add(progressBar)
	processContainerTop.Add(stopButton)

	processContainer := container.NewBorder(
		container.NewBorder(nil, nil, nil, stopButton, progressBar),
		nil, nil, nil,
		list,
	)

	return MyProcessContainer{
		Container:   processContainer,
		ProgressBar: progressBar,
		StopButton:  stopButton,
		List:        list,
	}
}
