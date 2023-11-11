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
startConvertFolderMp3Options serves as a way to pass arguments to startConvertFolderMp3
*/
type startConvertFolderMp3Options struct {
	dirPath *string
}

/*
startConvertFolderMp3 is the entrypoint for the ConvertFolderMp3 operation from the UI
*/
func (d *Data) startConvertFolderMp3(w fyne.Window, processContainerOuter *fyne.Container, opts startConvertFolderMp3Options) {

	ctx, op, err := d.sharedStartBuild(w, processContainerOuter)

	if err != nil {
		showErrorDialog(w, err)
		return
	}

	go operations.ConvertFolderMp3(
		ctx,
		*d.Config,
		op,
		operations.ConvertFolderMp3Params{
			InDirPath: *opts.dirPath,
		},
	)
}
