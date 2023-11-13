package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/billiem/seren-management/src/helpers"
	"github.com/billiem/seren-management/src/operations"
)

func (d *Data) sharedStartBuild(w fyne.Window, processContainerOuter *fyne.Container) (context.Context, operationProcess, error) {

	processContainerOuter.Objects = nil

	ctx := context.Background()

	if d.State.processing {
		return ctx, operationProcess{}, helpers.ErrPleaseWaitForProcess
	}

	d.State.processing = true

	ctx, cancelCauseFunc := context.WithCancelCause(ctx)

	progressBarBinding := binding.NewFloat()
	listBinding := progressBindingList{}

	processContainer := buildProcessContainer(cancelCauseFunc, progressBarBinding, &listBinding)

	stepFunc := func() {
		processContainer.log.ScrollToBottom()
	}

	finishedFunc := func() {
		processContainer.container.Remove(processContainer.stopButton)
		d.State.processing = false
		showInfoDialog(w, "Finished", "Process finished")
	}

	context.AfterFunc(ctx, func() {
		listBinding.Append(&progressBindingItem{message: "Process cancelled, finishing already running steps, please wait..."})
		stepFunc()
	})

	processContainerOuter.Add(processContainer.container)

	return ctx, operationProcess{
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
func (d *Data) startSeparateSingleStem(w fyne.Window, processContainerOuter *fyne.Container, opts operations.SeparateSingleStemOpts) {

	ctx, op, err := d.sharedStartBuild(w, processContainerOuter)

	if err != nil {
		showErrorDialog(w, err)
		return
	}

	go operations.SeparateSingleStem(
		ctx,
		*d.Config,
		op,
		opts,
	)
}

/*
startSeparateFolderStem is the entrypoint for the SeperateFolderStem operation from the UI
*/
func (d *Data) startSeparateFolderStem(w fyne.Window, processContainerOuter *fyne.Container, opts operations.SeparateFolderStemOpts) {

	ctx, op, err := d.sharedStartBuild(w, processContainerOuter)

	if err != nil {
		showErrorDialog(w, err)
		return
	}

	go operations.SeparateFolderStem(
		ctx,
		*d.Config,
		op,
		opts,
	)
}

/*
startConvertSingleMp3 is the entrypoint for the ConvertSingleMp3 operation from the UI
*/
func (d *Data) startConvertSingleMp3(w fyne.Window, processContainerOuter *fyne.Container, opts operations.ConvertSingleMp3Opts) {

	ctx, op, err := d.sharedStartBuild(w, processContainerOuter)

	if err != nil {
		showErrorDialog(w, err)
		return
	}

	go operations.ConvertSingleMp3(
		ctx,
		*d.Config,
		op,
		opts,
	)
}

/*
startConvertFolderMp3 is the entrypoint for the ConvertFolderMp3 operation from the UI
*/
func (d *Data) startConvertFolderMp3(w fyne.Window, processContainerOuter *fyne.Container, opts operations.ConvertFolderMp3Opts) {

	ctx, op, err := d.sharedStartBuild(w, processContainerOuter)

	if err != nil {
		showErrorDialog(w, err)
		return
	}

	go operations.ConvertFolderMp3(
		ctx,
		*d.Config,
		op,
		opts,
	)
}
