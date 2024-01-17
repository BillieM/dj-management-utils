package iwidget

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/fyne-io/terminal"
)

/*
Contains a collection of widgets used by in progress operations in order
to display the progress of the operation to the user, and to allow the user to cancel the operation

Progress is displayed as a progress bar, and a collection of messages reported back to the interface
by the in progress operation

The main widget is the RunningOperation widget, which serves as a container for the other widgets within this package
*/

/*
RunningOperation is the main widget for displaying the progress of an operation to the user,
it serves as a container for the other widgets within this package
*/
type RunningOperation struct {
	*Base
	widget.BaseWidget

	cancelFunc func()

	ProgressBar *widget.ProgressBar
	StopButton  *widget.Button
	Log         *terminal.Terminal
}

func NewRunningOperation(widgetBase *Base) *RunningOperation {

	runningOperation := &RunningOperation{
		Base:        widgetBase,
		ProgressBar: widget.NewProgressBar(),
		StopButton:  widget.NewButton("Stop", func() {}),
		Log:         terminal.New(),
	}

	runningOperation.StopButton.OnTapped = func() {
		runningOperation.StopButton.Disable()
		runningOperation.cancelFunc()
	}

	runningOperation.ExtendBaseWidget(runningOperation)

	return runningOperation

}

func (r *RunningOperation) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(
		container.NewBorder(
			container.NewBorder(nil, nil, nil, r.StopButton, r.ProgressBar),
			nil, nil, nil,
			container.NewVScroll(r.Log),
		),
	)
}

func (r *RunningOperation) SetCancelFunc(f func()) {
	r.cancelFunc = f
}

// TODO: custom layout method
// terminal is set as visible whereas progress bar and stop button are set as invisible
// this is a workaround to the terminal widget taking more time to render than the other widgets
// may aswell render the terminal where it should be and make the other widgets visible later
