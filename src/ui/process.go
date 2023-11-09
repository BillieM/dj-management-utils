package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*


 */

type OperationProcess struct {
	progressBar MyProgressBar
}

func (o OperationProcess) StepCallback(progress float64) {
	fmt.Println("step callback progress", progress)

	o.progressBar.updateProgressBar(progress)
}

/*


 */

type MyProcessContainer struct {
	Container   *fyne.Container
	StopButton  MyStopButton
	ProgressBar MyProgressBar
}

func buildProcessContainer(cancelFunc context.CancelFunc) MyProcessContainer {
	processContainer := container.NewVBox()

	progressBar := buildProgressBar()

	stopButton := buildStopButton(cancelFunc)

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
	fmt.Println("update progress bar progress", value)
	p.SetValue(value)
}

/*


 */

type MyStopButton struct {
	*widget.Button
}

func buildStopButton(cancelFunc context.CancelFunc) MyStopButton {
	return MyStopButton{
		widget.NewButton("Stop", func() {
			cancelFunc()
		}),
	}
}
