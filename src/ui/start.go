package ui

import (
	"context"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/data/binding"
	"github.com/billiem/seren-management/src/operations"
)

type startConvertFolderMp3Options struct {
	dirPath *string
}

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
	listBinding := binding.NewStringList()

	processContainer := buildProcessContainer(cancelCauseFunc, progressBarBinding, listBinding)

	stepFunc := func() {
		processContainer.List.ScrollToBottom()
	}

	finishedFunc := func() {
		processContainerOuter.Remove(processContainer.Container)
		d.State.processing = false
		dialogInfo(w, "Finished", "Process finished")
	}

	context.AfterFunc(ctx, func() {
		listBinding.Append("Operation cancelled, finishing off running processes...")
		stepFunc()
	})

	processContainerOuter.Add(processContainer.Container)

	go operations.ConvertFolderMp3(
		ctx,
		*d.Config,
		OperationProcess{
			ctxClose:             cancelCauseFunc,
			listBindValue:        listBinding,
			progressBarBindValue: progressBarBinding,
			stepFunc:             stepFunc,
			finishedFunc:         finishedFunc,
		},
		operations.ConvertFolderMp3Params{
			InDirPath: *opts.dirPath,
		},
	)
}
