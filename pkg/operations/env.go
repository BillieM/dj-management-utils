package operations

import (
	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
)

type OpEnv struct {
	helpers.Config
	*database.SerenDB
	*stepHandler
	*stepHandlerNew
}

func (e *OpEnv) RegisterStepHandler(sh StepHandler) {
	e.stepHandler = &stepHandler{
		stepCallback: sh.StepCallback,
		exitCallback: sh.ExitCallback,
	}
}

func (e *OpEnv) RegisterStepHandlerNew(sh StepHandlerNew) {
	e.stepHandlerNew = &stepHandlerNew{
		stepCallback:    sh.StepCallback,
		successCallback: sh.SuccessCallback,
		errorCallback:   sh.ErrorCallback,
	}
}

func (e *OpEnv) step(stepInfo StepInfo) {
	e.stepHandler.stepCallback(stepInfo)
}

func (e *OpEnv) exit() {
	e.stepHandler.exitCallback()
}

func (e *OpEnv) stepNew(stepInfo StepInfo) {
	e.stepHandlerNew.stepCallback(stepInfo)
}
