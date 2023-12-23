package gui

import (
	"fyne.io/fyne/v2"
	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/gui/uihelpers"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

/*
guiEnv holds the environment for the GUI
*/
type guiEnv struct {
	*helpers.Config
	*data.SerenDB
	*helpers.AppLogger
	*guiState
	tmpConfig    *helpers.Config
	views        map[string]guiView
	viewIndices  map[string][]string
	mainWindow   fyne.Window
	app          fyne.App
	resizeEvents *uihelpers.ResizeEvents
}

func (e *guiEnv) opEnv() *operations.OpEnv {
	return &operations.OpEnv{
		Config:  *e.Config,
		Logger:  e.AppLogger.OPLogger,
		SerenDB: e.SerenDB,
	}
}

/*
buildGuiEnv builds the *guiEnv struct
*/
func buildGuiEnv(a fyne.App, w fyne.Window) (*guiEnv, error) {

	cfg, err := helpers.LoadGUIConfig()

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error loading GUI config"))
	}

	logger, err := helpers.BuildAppLogger(*cfg)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error building logger"))
	}

	queries, err := data.Connect(*cfg, *logger)

	if err != nil {
		return nil, fault.Wrap(err, fmsg.With("Error connecting to database"))
	}

	e := &guiEnv{cfg, queries, logger, nil, nil, nil, nil, w, a, nil}

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
