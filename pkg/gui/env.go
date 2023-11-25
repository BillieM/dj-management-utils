package gui

import (
	"fyne.io/fyne/v2"
	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
)

/*
guiEnv holds the environment for the GUI
*/
type guiEnv struct {
	*helpers.Config
	*database.SerenDB
	*guiState
	tmpConfig   *helpers.Config
	views       map[string]guiView
	viewIndices map[string][]string
	mainWindow  fyne.Window
}

func (e *guiEnv) opEnv() *operations.OpEnv {
	return &operations.OpEnv{
		Config:  *e.Config,
		SerenDB: e.SerenDB,
	}
}

/*
buildGuiEnv builds the *guiEnv struct
*/
func buildGuiEnv(w fyne.Window) (*guiEnv, error) {

	cfg, err := helpers.LoadGUIConfig()

	if err != nil {
		return nil, err
	}

	db, err := database.Connect()

	if err != nil {
		return nil, err
	}

	e := &guiEnv{cfg, db, nil, nil, nil, nil, w}

	s := &guiState{}
	operations := e.getViewList()
	operationIndex := e.getViewIndex()

	e.guiState = s
	e.views = operations
	e.viewIndices = operationIndex

	return e, nil
}

type guiState struct {
	settingsAlreadyOpen bool
	busy                bool
}
