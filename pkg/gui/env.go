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
}

func (e *guiEnv) opEnv() operations.OpEnv {
	return operations.OpEnv{
		Config:  *e.Config,
		SerenDB: e.SerenDB,
	}
}

/*
buildGuiEnv builds the *guiEnv struct
*/
func buildGuiEnv() (*guiEnv, error) {

	cfg, err := helpers.LoadGUIConfig()

	if err != nil {
		return nil, err
	}

	db, err := database.Connect()

	if err != nil {
		return nil, err
	}

	e := &guiEnv{cfg, db, nil, nil, nil, nil}

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
	processing          bool
}

/*
guiView is a struct that contains the name of the view and the function that returns the view

Views define the content for the main content area of the application
*/
type guiView struct {
	name   string
	render func(w fyne.Window) fyne.CanvasObject
}

func (e *guiEnv) getViewList() map[string]guiView {

	return map[string]guiView{
		"home": {
			name:   "Home",
			render: e.homeView,
		},
		"stems": {
			name:   "Stems",
			render: e.stemsView,
		},
		"separateTrack": {
			name:   "Separate Track",
			render: e.separateSingleStemView,
		},
		"separateFolder": {
			name:   "Separate Folder",
			render: e.separateFolderStemView,
		},
		"separateCollection": {
			name:   "Separate Collection",
			render: e.separateCollectionStemView,
		},
		"mp3s": {
			name:   "Convert MP3s",
			render: e.convertMp3sView,
		},
		"convertSingleMp3": {
			name:   "Convert Single",
			render: e.convertSingleMp3View,
		},
		"convertFolderMp3": {
			name:   "Convert Folder",
			render: e.convertFolderMp3View,
		},
		"convertCollectionMp3": {
			name:   "Convert Collection",
			render: e.convertCollectionMp3View,
		},
		"tags": {
			name:   "Process Tags",
			render: e.tagsView,
		},
		"rereadTags": {
			name:   "Reread Tags",
			render: e.rereadTagsView,
		},
		"cleanTags": {
			name:   "Clean Tags",
			render: e.cleanTagsView,
		},
		"conversion": {
			name:   "Conversion",
			render: e.conversionView,
		},
		"playlistMatching": {
			name:   "Playlist Matching",
			render: e.playlistMatchingView,
		},
	}
}

func (e *guiEnv) getViewIndex() map[string][]string {
	return map[string][]string{
		"": {"home", "stems", "mp3s", "tags", "conversion", "playlistMatching"},
		"stems": {
			"separateTrack",
			"separateFolder",
			"separateCollection",
		},
		"mp3s": {
			"convertSingleMp3",
			"convertFolderMp3",
			"convertCollectionMp3",
		},
		"tags": {
			"rereadTags",
			"cleanTags",
		},
	}
}
