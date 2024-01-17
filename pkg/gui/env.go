package gui

import (
	"io"

	"fyne.io/fyne/v2"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/gui/uihelpers"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

/*
guiEnv serves as the environment for the GUI, and is passed around much of the graphical
code in the application.

We do not pass this directly to other packages, instead we generate a subset of this environment,
containing only relevant fields for the package.
*/
type guiEnv struct {
	*helpers.Config
	*data.SerenDB
	logger helpers.SerenLogger
	*guiState
	tmpConfig    *helpers.Config
	views        map[string]guiView
	viewIndices  map[string][]string
	mainWindow   fyne.Window
	app          fyne.App
	resizeEvents *uihelpers.ResizeEvents
	termSink     *helpers.TermSink
}

/*
opEnv returns an OpEnv struct for use in operations,
this is generated from the guiEnv struct
*/
func (e *guiEnv) opEnv() *operations.OpEnv {
	return &operations.OpEnv{
		Config:  *e.Config,
		Logger:  helpers.BuildOperationLogger(*e.Config, e.termSink),
		SerenDB: e.SerenDB,
	}
}

/*
getWidgetBase returns a *iwidget.Base struct for use in custom widgets,
this is generated from the guiEnv struct
*/
func (e *guiEnv) getWidgetBase() *iwidget.Base {
	return &iwidget.Base{
		Logger:       e.logger,
		Config:       *e.Config,
		MainWindow:   e.mainWindow,
		App:          e.app,
		ResizeEvents: e.resizeEvents,
	}
}

/*
buildGuiEnv returns a pointer to an instance of the guiEnv struct with
the environment loaded.
*/
func buildGuiEnv(a fyne.App, w fyne.Window) (*guiEnv, error) {

	cfg, err := helpers.LoadGUIConfig()

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error loading GUI config"))
	}

	loggers, err := helpers.BuildAppLoggers(*cfg)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error building logger"))
	}

	queries, err := data.Connect(*cfg, loggers.DBLogger)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error connecting to database"))
	}

	termSink := helpers.BuildTermSink(io.Discard)

	e := &guiEnv{cfg, queries, loggers.AppLogger, nil, nil, nil, nil, w, a, nil, termSink}

	s := &guiState{}
	operations := e.getViewList()
	operationIndex := e.getViewIndex()
	resizeEvents := uihelpers.BuildResizeEvents()

	e.guiState = s
	e.views = operations
	e.viewIndices = operationIndex
	e.resizeEvents = resizeEvents

	return e, nil
}

type guiState struct {
	settingsAlreadyOpen bool
	busy                bool
}

/*
isBusy returns true if the GUI is currently busy with an operation

It also displays an error dialog to the user if the GUI is busy
*/
func (e *guiEnv) isBusy() bool {
	if e.guiState.busy {
		e.showErrorDialog(fault.Wrap(
			fault.New("gui state busy"),
			fmsg.WithDesc(
				"wait for current operation to finish",
				"Busy, please wait for the current operation to finish",
			),
		), false)
		return true
	}
	return false
}
