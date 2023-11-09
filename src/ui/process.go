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


 */

type OperationProcess struct {
	// TODO: may need to use external bindings? idk what they do
	listBindValue        binding.StringList
	progressBarBindValue binding.Float
	ctxClose             context.CancelCauseFunc
}

func (o OperationProcess) StepCallback(progress float64, message string) {
	fmt.Println("step callback", progress, message)
	o.progressBarBindValue.Set(progress)
	// TODO: may want to add a nicer handler for this
	// i.e. to show only a limited number of items at any one time
	// o.listBindValue.Append(message)
}

func (o OperationProcess) ExitCallback() {
	o.ctxClose(errors.New("Operation finished"))
}

/*


 */

type MyProcessContainer struct {
	Container   *fyne.Container
	StopButton  MyStopButton
	ProgressBar *widget.ProgressBar
	List        *widget.List
}

// TODO: rename all the list calls, should be called 'log' or something
func buildProcessContainer(cancelCauseFunc context.CancelCauseFunc, progressBarBindVal binding.Float, listBindVal binding.StringList) MyProcessContainer {
	processContainer := container.NewVBox()

	progressBar := buildProgressBar(progressBarBindVal)
	list := buildList(listBindVal)
	stopButton := buildStopButton(cancelCauseFunc)

	processContainer.Add(widget.NewLabel("Processing..."))
	processContainer.Add(progressBar)
	processContainer.Add(stopButton)
	processContainer.Add(list)

	return MyProcessContainer{
		Container:   processContainer,
		ProgressBar: progressBar,
		StopButton:  stopButton,
		List:        list,
	}
}

/*


 */

func buildProgressBar(bindVal binding.Float) *widget.ProgressBar {
	return widget.NewProgressBarWithData(bindVal)
}

/*


 */

func buildList(bindVal binding.StringList) *widget.List {
	list := widget.NewList(
		func() int {
			return 0
		},
		func() fyne.CanvasObject {
			return widget.NewLabel("template")
		},
		func(i widget.ListItemID, o fyne.CanvasObject) {
			fmt.Println("update item", i)
		},
	)

	return list
}

func updateListBinding(list binding.StringList) {

}

/*


 */

type MyStopButton struct {
	*widget.Button
}

func buildStopButton(cancelCauseFunc context.CancelCauseFunc) MyStopButton {
	return MyStopButton{
		widget.NewButton("Stop", func() {
			cancelCauseFunc(errors.New("user stopped process"))
		}),
	}
}
