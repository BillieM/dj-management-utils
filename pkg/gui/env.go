package gui

import (
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
}

/*
opEnv returns an OpEnv struct for use in operations,
this is generated from the guiEnv struct
*/
func (e *guiEnv) opEnv() *operations.OpEnv {
	return &operations.OpEnv{
		Config:  *e.Config,
		Logger:  helpers.BuildOperationLogger(*e.Config),
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

	e := &guiEnv{cfg, queries, loggers.AppLogger, nil, nil, nil, nil, w, a, nil}

	s := &guiState{}
	operations := e.getViewList()
	operationIndex := e.getViewIndex()
	resizeEvents := uihelpers.NewResizeEvents()

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
