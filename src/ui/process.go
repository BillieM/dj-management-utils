package ui

import (
	"context"
	"errors"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*


 */

type OperationProcess struct {
	progressBar MyProgressBar
	ctxClose    context.CancelCauseFunc
}

func (o OperationProcess) StepCallback(progress float64) {
	fmt.Println("step callback progress", progress)

	o.progressBar.updateProgressBar(progress)
}

func (o OperationProcess) ExitCallback() {
	o.ctxClose(errors.New("Operation finished"))
}

/*


 */

type MyProcessContainer struct {
	Container   *fyne.Container
	StopButton  MyStopButton
	ProgressBar MyProgressBar
}

func buildProcessContainer(cancelCauseFunc context.CancelCauseFunc) MyProcessContainer {
	processContainer := container.NewVBox()

	progressBar := buildProgressBar()

	stopButton := buildStopButton(cancelCauseFunc)

	processContainer.Add(widget.NewLabel("Processing..."))
	processContainer.Add(progressBar)
	processContainer.Add(stopButton)

	return MyProcessContainer{
		Container:   processContainer,
		ProgressBar: progressBar,
		StopButton:  stopButton,
	}
}

/*


 */

type MyProgressBar struct {
	*widget.ProgressBar
}

func buildProgressBar() MyProgressBar {
	return MyProgressBar{
		widget.NewProgressBar(),
	}
}

func (p MyProgressBar) updateProgressBar(value float64) {
	p.SetValue(value)
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
