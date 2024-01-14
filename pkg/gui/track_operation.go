package gui

import (
	"context"
	"io"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/operations"
)

type trackOperation struct {
	ctx              context.Context
	opEnv            *operations.OpEnv
	runningOperation *iwidget.RunningOperation
}

/*
prepareTrackOperation prepares a track operation, trackOperations are used for stem separation and mp3 conversion
*/
func (e *guiEnv) prepareTrackOperation() trackOperation {
	e.guiState.busy = true

	ctx := context.Background()
	ctx, ctxClose := context.WithCancel(ctx)

	opEnv := e.opEnv()

	// build the running operation widget, which will be used to display the progress of the operation & allow it to be cancelled
	runningOperation := iwidget.NewRunningOperation(e.getWidgetBase(), func() {
		ctxClose()
		opEnv.Logger.Info("Cancelled operation, finishing up running tasks, please wait...")
	})

	pr, pw := io.Pipe()

	opEnv.RegisterOperationHandler(
		func(i float64) {
			runningOperation.ProgressBar.SetValue(i)
		},
		func(i operations.OperationFinishedInfo) {
			// pw.Close()
			runningOperation.ProgressBar.SetValue(1)
			e.guiState.busy = false
			if i.Err != nil {
				e.showErrorDialog(i.Err)
				return
			}
			e.showInfoDialog("Finished", "Process finished")
		},
	)

	// add the terminal writer to the logger
	opEnv.Logger.AddTermCore(pw)

	go func() {

		err := runningOperation.Log.RunWithConnection(nil, pr)

		if err != nil {
			e.logger.NonFatalError(fault.Wrap(
				err,
				fmsg.With("error running terminal"),
			))
		}
	}()

	return trackOperation{
		ctx:              ctx,
		opEnv:            opEnv,
		runningOperation: runningOperation,
	}
}
