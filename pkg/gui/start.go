package gui

import (
	"context"
	"io"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

func (e *guiEnv) sharedStartBuild(processContainerOuter *fyne.Container) (context.Context, convertStepHandler, error) {

	processContainerOuter.Objects = nil

	ctx := context.Background()

	if e.guiState.busy {
		return ctx, convertStepHandler{}, helpers.ErrBusyPleaseFinishFirst
	}

	e.guiState.busy = true

	ctx, cancelCauseFunc := context.WithCancelCause(ctx)

	progressBarBinding := binding.NewFloat()
	listBinding := progressBindingList{}

	processContainer := buildProcessContainer(cancelCauseFunc, progressBarBinding, &listBinding)

	stepFunc := func() {
		processContainer.log.ScrollToBottom()
	}

	finishedFunc := func() {
		processContainer.container.Remove(processContainer.stopButton)
		e.guiState.busy = false
		e.showInfoDialog("Finished", "Process finished")
	}

	context.AfterFunc(ctx, func() {
		listBinding.Append(&progressBindingItem{message: "Process cancelled, finishing already running steps, please wait..."})
		stepFunc()
	})

	processContainerOuter.Add(processContainer.container)

	return ctx, convertStepHandler{
		ctxClose:             cancelCauseFunc,
		bindVals:             &listBinding,
		progressBarBindValue: progressBarBinding,
		stepFunc:             stepFunc,
		finishedFunc:         finishedFunc,
	}, nil
}

/*
startSeparateSingleStem is the entrypoint for the SeperateSingleStem operation from the UI
*/
func (e *guiEnv) startSeparateSingleStem(processContainerOuter *fyne.Container, opts operations.SeparateSingleStemOpts) {

	ctx := context.Background()
	ctx, ctxClose := context.WithCancel(ctx)

	// build the running operation widget, which will be used to display the progress of the operation & allow it to be cancelled
	runningOperation := iwidget.NewRunningOperation(e.getWidgetBase(), ctxClose)

	processContainerOuter.Add(runningOperation)

	pr, pw := io.Pipe()

	opEnv := e.opEnv()
	opEnv.RegisterOperationHandler(
		func(i operations.OperationProgressInfo) {
			runningOperation.ProgressBar.SetValue(i.Progress)
		},
		func(i operations.OperationFinishedInfo) {
			pw.Close()
			if i.Err != nil {
				e.displayErrorDialog(i.Err)
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

	go opEnv.SeparateSingleStem(
		ctx,
		opts,
	)
}

/*
startSeparateFolderStem is the entrypoint for the SeperateFolderStem operation from the UI
*/
func (e *guiEnv) startSeparateFolderStem(processContainerOuter *fyne.Container, opts operations.SeparateFolderStemOpts) {

	ctx, sh, err := e.sharedStartBuild(processContainerOuter)

	if err != nil {
		e.showErrorDialog(err)
		return
	}

	opEnv := e.opEnv()
	opEnv.RegisterStepHandler(sh)

	go opEnv.SeparateFolderStem(
		ctx,
		opts,
	)
}

/*
startConvertSingleMp3 is the entrypoint for the ConvertSingleMp3 operation from the UI
*/
func (e *guiEnv) startConvertSingleMp3(processContainerOuter *fyne.Container, opts operations.ConvertSingleMp3Opts) {

	ctx, sh, err := e.sharedStartBuild(processContainerOuter)

	if err != nil {
		e.showErrorDialog(err)
		return
	}

	opEnv := e.opEnv()
	opEnv.RegisterStepHandler(sh)

	go opEnv.ConvertSingleMp3(
		ctx,
		opts,
	)
}

/*
startConvertFolderMp3 is the entrypoint for the ConvertFolderMp3 operation from the UI
*/
func (e *guiEnv) startConvertFolderMp3(processContainerOuter *fyne.Container, opts operations.ConvertFolderMp3Opts) {

	ctx, sh, err := e.sharedStartBuild(processContainerOuter)

	if err != nil {
		e.showErrorDialog(err)
		return
	}

	opEnv := e.opEnv()
	opEnv.RegisterStepHandler(sh)

	go opEnv.ConvertFolderMp3(
		ctx,
		opts,
	)
}
