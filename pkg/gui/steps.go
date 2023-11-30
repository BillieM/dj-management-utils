package gui

import (
	"fmt"

	"github.com/billiem/seren-management/pkg/operations"
)

type streamingStepHandler struct {
	stepFunc     func()
	finishedFunc func()
}

func (s streamingStepHandler) StepCallback(step operations.StepInfo) {
	fmt.Println(step)
	s.stepFunc()
}

func (s streamingStepHandler) ExitCallback() {
	fmt.Println("finished")
	s.finishedFunc()
}

type streamingStepHandlerNew struct {
	stepCallback     func(operations.StepInfoNew)
	finishedCallback func(operations.FinishedInfo)
}

func (s streamingStepHandlerNew) StepCallback(i operations.StepInfoNew) {
	fmt.Println(i)
	s.stepCallback(i)
}

func (s streamingStepHandlerNew) FinishedCallback(i operations.FinishedInfo) {
	fmt.Println(i)
	s.finishedCallback(i)
}
