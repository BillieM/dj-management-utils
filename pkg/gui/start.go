package gui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

func (e *guiEnv) sharedStartBuild(processContainerOuter *fyne.Container) (context.Context, convertStepHandler, error) {

	processContainerOuter.Objects = nil

	ctx := context.Background()

	if e.guiState.busy {
		return ctx, convertStepHandler{}, helpers.ErrPleaseWaitForProcess
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

	ctx, sh, err := e.sharedStartBuild(processContainerOuter)

	if err != nil {
		e.showErrorDialog(err)
		return
	}

	opEnv := e.opEnv()
	opEnv.RegisterStepHandler(sh)

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

	go e.opEnv().ConvertFolderMp3(
		ctx,
		opts,
	)
}
