package operations

import (
	"github.com/billiem/seren-management/pkg/helpers"
)

// stepHandler should be registered at the start of each operation
// it handles callbacks & logging for each step
type stepHandler struct {
	stepCallback func(StepInfo)
	exitCallback func()
}

/*
Experiment to replace the existing stepHandler with a more generic one
*/
type stepHandlerNew struct {
	stepCallback     func(StepInfoNew)
	finishedCallback func(FinishedInfo)
}

/*
StepHandler is used to provide callbacks to the operations package
*/
type StepHandler interface {
	StepCallback(StepInfo)
	ExitCallback()
}

type StepHandlerNew interface {
	StepCallback(StepInfoNew)
	FinishedCallback(FinishedInfo)
}

/*
StepInfo is returned to the StepCallback after each step

It provides information about the step that just finished
*/
type StepInfo struct {
	SkipLog    bool
	Progress   float64
	Message    string
	Error      error
	Importance helpers.Importance
}

type StepInfoNew struct {
	Progress   float64 // value between 0 and 1
	Message    string
	Importance helpers.Importance
}

type FinishedInfo struct {
	Success bool
	Data    any // TODO: consider changing this to an interface
	err     error
}

/*
Helper functions for building StepInfo objects
*/

func processFinishedStepInfo(msg string) StepInfo {
	return StepInfo{
		Message:    msg,
		Progress:   1,
		Importance: helpers.HighImportance,
	}
}

func stageStepInfo(msg string) StepInfo {
	return StepInfo{
		Message:    msg,
		Importance: helpers.HighImportance,
	}
}

func stepStartedStepInfo(msg string) StepInfo {
	return StepInfo{
		Message:    msg,
		Importance: helpers.MediumImportance,
	}
}

func stepFinishedStepInfo(msg string, progress float64) StepInfo {
	return StepInfo{
		Message:    msg,
		Progress:   progress,
		Importance: helpers.MediumImportance,
	}
}

func trackFinishedStepInfo(msg string, progress float64) StepInfo {
	return StepInfo{
		Message:    msg,
		Progress:   progress,
		Importance: helpers.SuccessImportance,
	}
}

func stepWarningStepInfo(err error, progress float64) StepInfo {
	return StepInfo{
		Error:      err,
		Progress:   progress,
		Importance: helpers.WarningImportance,
	}
}

func warningStepInfo(err error) StepInfo {
	return stepWarningStepInfo(err, 0)
}

func dangerStepInfo(err error) StepInfo {
	return StepInfo{
		Error:      err,
		Importance: helpers.DangerImportance,
	}
}

func progressOnlyStepInfo(progress float64) StepInfo {
	return StepInfo{
		Progress: progress,
		SkipLog:  true,
	}
}
