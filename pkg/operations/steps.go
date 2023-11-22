package operations

import (
	"github.com/billiem/seren-management/pkg/helpers"
	"gorm.io/gorm"
)

// stepHandler should be registered at the start of each operation
// it handles callbacks & logging for each step
type stepHandler struct {
	db       *gorm.DB
	callback func(StepInfo)
}

func (s stepHandler) step(stepInfo StepInfo) {
	s.callback(stepInfo)
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
		Importance: helpers.LowImportance,
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
