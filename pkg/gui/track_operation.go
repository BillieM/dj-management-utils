package gui

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

/*
prepareTrackOperation prepares a track operation, trackOperations are used for stem separation and mp3 conversion
*/
func (e *guiEnv) prepareTrackOperation() (*operations.OpEnv, *iwidget.RunningOperation) {

	opEnv := e.opEnv()

	runningOperation := iwidget.NewRunningOperation(e.getWidgetBase())
	runningOperation.ProgressBar.Hide()
	runningOperation.StopButton.Hide()

	e.termSink.SetIO(io.Pipe())

	opEnv.BuildOperationHandler(
		func(i float64) {
			runningOperation.ProgressBar.SetValue(i)
		},
		func(_ map[string]any) {
			runningOperation.StopButton.Disable()
			e.guiState.busy = false
			e.showInfoDialog("Finished", "Process finished")
		},
		func(err error) {
			runningOperation.StopButton.Disable()
			e.guiState.busy = false
			e.showErrorDialog(err, true)
		},
	)

	go func() {

		err := runningOperation.Log.RunWithConnection(
			helpers.NewDiscardCloser(),
			e.termSink.Reader,
		)

		if err != nil {
			e.logger.NonFatalError(fault.Wrap(
				err,
				fmsg.With("error running terminal"),
			))
		}
	}()

	return opEnv, runningOperation
}

type execTrackOperationOpts struct {
	execFunc         func(context.Context)
	runningOperation *iwidget.RunningOperation
	opEnv            *operations.OpEnv
}

func (e *guiEnv) executeTrackOperation(opts *execTrackOperationOpts) {
	if e.isBusy() {
		return
	}

	e.guiState.busy = true

	ctx, ctxClose := context.WithCancel(context.Background())

	opts.runningOperation.SetCancelFunc(func() {
		opts.opEnv.Logger.Info("Cancelled operation, finishing up running tasks, please wait...")
		ctxClose()
	})

	opts.runningOperation.StopButton.Show()
	opts.runningOperation.StopButton.Enable()

	opts.runningOperation.ProgressBar.Show()
	opts.runningOperation.ProgressBar.SetValue(0)

	go func() {
		opts.execFunc(ctx)
	}()
}
