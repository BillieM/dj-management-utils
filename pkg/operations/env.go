package operations

import (
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
)

type OpEnv struct {
	helpers.Config
	Logger helpers.SerenLogger
	operationHandler
	*data.SerenDB
	*stepHandler
}

func (e *OpEnv) RegisterStepHandler(sh StepHandler) {
	e.stepHandler = &stepHandler{
		stepCallback: sh.StepCallback,
		exitCallback: sh.ExitCallback,
	}
}

func (e *OpEnv) step(stepInfo StepInfo) {
	e.stepHandler.stepCallback(stepInfo)
}

func (e *OpEnv) exit() {
	e.stepHandler.exitCallback()
}
