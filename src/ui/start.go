package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/billiem/seren-management/src/operations"
)

func sharedStartBuild() {

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

	if d.State.processing {
		pleaseWaitForProcess(w)
		return
	}

	d.State.processing = true

	// new context & cancel function
	ctx := context.Background()
	ctx, cancelCauseFunc := context.WithCancelCause(ctx)

	progressBarBinding := binding.NewFloat()
	listBinding := progressBindingList{}

	processContainer := buildProcessContainer(cancelCauseFunc, progressBarBinding, &listBinding)

	stepFunc := func() {
		processContainer.log.ScrollToBottom()
	}

	finishedFunc := func() {
		processContainerOuter.Remove(processContainer.container)
		d.State.processing = false
		dialogInfo(w, "Finished", "Process finished")
	}

	context.AfterFunc(ctx, func() {
		listBinding.Append(&progressBindingItem{message: "Process cancelled"})
		stepFunc()
	})

	processContainerOuter.Add(processContainer.container)

	go operations.ConvertFolderMp3(
		ctx,
		*d.Config,
		OperationProcess{
			ctxClose:             cancelCauseFunc,
			bindVals:             &listBinding,
			progressBarBindValue: progressBarBinding,
			stepFunc:             stepFunc,
			finishedFunc:         finishedFunc,
		},
		operations.ConvertFolderMp3Params{
			InDirPath: *opts.dirPath,
		},
	)
}
