package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/src/operations"
)

/*
Checks the config for any issues for a given set of checks

# Returns true if there are no issues, false if there are issues

If there are issues, it will return a fyne.CanvasObject containing the issues
*/
func (d *Data) checkConfig(checks []func() (bool, string)) (bool, fyne.CanvasObject) {

	configIssues := []string{}

	for _, check := range checks {
		pass, msg := check()
		if !pass {
			configIssues = append(configIssues, msg)
		}
	}

	if len(configIssues) > 0 {
		issuesContainer := container.NewVBox(
			widget.NewLabel("Please fix the following issues with your config:"),
		)
		for _, issue := range configIssues {
			issuesContainer.Add(widget.NewLabel(issue))
		}
		return false, issuesContainer
	}

	return true, nil
}

func buildStemTypeSelect(t *operations.StemSeparationType, callbackFn func()) *widget.Select {
	w := widget.NewSelect(
		[]string{"Traktor Stem File", "4 Stem Files"},
		func(s string) {
			if s == "Traktor Stem File" {
				*t = operations.Traktor
			} else if s == "4 Stem Files" {
				*t = operations.FourTrack
			}
			callbackFn()
		},
	)
	w.PlaceHolder = "Please select the type of stem extraction you would like to perform"

	return w
}

func enableBtnIfOptsOkay(o operations.OperationOptions, btn *widget.Button) {
	ok, _ := o.Check()
	if ok {
		btn.Enable()
	}
}
