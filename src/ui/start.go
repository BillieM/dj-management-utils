package ui

import (
	"context"
	"fmt"

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

	context.AfterFunc(ctx, func() {
		fmt.Println("context closed because: ", context.Cause(ctx))

		processContainerOuter.Remove(processContainer.Container)

		d.State.processing = false
	})

	processContainerOuter.Add(processContainer.Container)

	go operations.ConvertFolderMp3(
		ctx,
		*d.Config,
		OperationProcess{
			progressBarBindValue: progressBarBinding,
			ctxClose:             cancelCauseFunc,
		},
		operations.ConvertFolderMp3Params{
			InDirPath: *opts.dirPath,
		},
	)
}
