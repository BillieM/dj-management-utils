package cli

import (
	"fmt"

	"github.com/billiem/seren-management/pkg/operations"
)

type operationProcess struct {
}

func (o operationProcess) StepCallback(stepInfo operations.StepInfo) {
	fmt.Println(stepInfo)
}

func (o operationProcess) ExitCallback() {
	fmt.Println("Finished")
}
