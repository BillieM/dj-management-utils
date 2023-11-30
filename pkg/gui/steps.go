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
	stepCallback    func(operations.StepInfo)
	successCallback func(any)
	errorCallback   func(error)
}

func (s streamingStepHandlerNew) StepCallback(step operations.StepInfo) {
	fmt.Println(step)
	s.stepCallback(step)
}

func (s streamingStepHandlerNew) SuccessCallback(i any) {
	fmt.Println(i)
	s.successCallback(i)
}

func (s streamingStepHandlerNew) ErrorCallback(err error) {
	fmt.Println(err)
	s.errorCallback(err)
}
