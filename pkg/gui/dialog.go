package gui

import (
	"errors"
	"strings"

	"fyne.io/fyne/v2/dialog"
	"github.com/Southclaws/fault/fmsg"
)

func (e *guiEnv) showErrorDialog(err error) {

	e.logger.NonFatalError(err)

	// user readable error messages (description of fmsg.withDesc())
	issues := fmsg.GetIssues(err)

	dialog.ShowError(errors.New(strings.Join(issues, "/n")), e.mainWindow)

}

func (e *guiEnv) showInfoDialog(title, message string) {
	dialog.ShowInformation(title, message, e.mainWindow)
}
