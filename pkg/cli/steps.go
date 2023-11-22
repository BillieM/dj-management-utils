package cli

import (
	"fmt"

	"github.com/billiem/seren-management/pkg/operations"
)

type stepHandler struct {
}

func (o stepHandler) StepCallback(stepInfo operations.StepInfo) {
	fmt.Println(stepInfo)
}

func (o stepHandler) ExitCallback() {
	fmt.Println("Finished")
}

/*
going to put some stuff in here to handle progress bar/ log in console

*/
