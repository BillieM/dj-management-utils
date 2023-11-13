package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

/*
Checks the config for any issues for a given set of checks

# Returns true if there are no issues, false if there are issues

If there are issues, it will return a fyne.CanvasObject containing the issues

	TODO: abstract the checking into a seperate function
		can then add unit tests surrounding it
		and create a seperate function for the generation of the canvas object
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
