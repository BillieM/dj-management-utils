package ui

import (
	"context"
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/operations"
)

type startConvertFolderMp3Options struct {
	dirPath *string
}

func (d *Data) startConvertFolderMp3(processContainerOuter *fyne.Container, startButton *widget.Button, opts startConvertFolderMp3Options) {

	startButton.Disable()

	d.State.processing = true

	// new context
	ctx := context.Background()
	// create cancelFunc from context
	ctx, cancelFunc := context.WithCancel(ctx)

	processContainer := buildProcessContainer(cancelFunc)

	context.AfterFunc(ctx, func() {
		fmt.Println("context closed because: ", context.Cause(ctx))

		processContainerOuter.Remove(processContainer.Container)

		d.State.processing = false

		startButton.Enable()
	})

	processContainerOuter.Add(processContainer.Container)

	go operations.ConvertFolderMp3(
		ctx,
		*d.Config,
		OperationProcess{
			progressBar: processContainer.ProgressBar,
		},
		operations.ConvertFolderMp3Params{
			InDirPath: *opts.dirPath,
		},
	)
}
