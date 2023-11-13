package ui

import (
	"fyne.io/fyne/v2"
	"github.com/billiem/seren-management/src/helpers"
)

type Data struct {
	*helpers.Config
	*State
	TmpConfig      *helpers.Config
	Operations     map[string]Operation
	OperationIndex map[string][]string
}

/*
builds the main data object for the application
*/
func buildData(c *helpers.Config) *Data {
	d := &Data{c, nil, nil, nil, nil}

	s := &State{}
	operations := d.getOperationsList()
	operationIndex := d.getOperationIndex()

	d.State = s
	d.Operations = operations
	d.OperationIndex = operationIndex

	return d
}

type State struct {
	settingsAlreadyOpen bool
	processing          bool
}

/*
Operations are the main views of the application
*/

type Operation struct {
	Name string
	View func(w fyne.Window) fyne.CanvasObject
}

func (d *Data) getOperationsList() map[string]Operation {

	return map[string]Operation{
		"home": {
			Name: "Home",
			View: d.homeView,
		},
		"stems": {
			Name: "Stems",
			View: d.stemsView,
		},
		"separateTrack": {
			Name: "Separate Track",
			View: d.separateSingleStemView,
		},
		"separateFolder": {
			Name: "Separate Folder",
			View: d.separateFolderStemView,
		},
		"separateCollection": {
			Name: "Separate Collection",
			View: d.separateCollectionStemView,
		},
		"mp3s": {
			Name: "Convert MP3s",
			View: d.convertMp3sView,
		},
		"convertSingleMp3": {
			Name: "Convert Single",
			View: d.convertSingleMp3View,
		},
		"convertFolderMp3": {
			Name: "Convert Folder",
			View: d.convertFolderMp3View,
		},
		"convertCollectionMp3": {
			Name: "Convert Collection",
			View: d.convertCollectionMp3View,
		},
		"tags": {
			Name: "Process Tags",
			View: d.tagsView,
		},
		"rereadTags": {
			Name: "Reread Tags",
			View: d.rereadTagsView,
		},
		"cleanTags": {
			Name: "Clean Tags",
			View: d.cleanTagsView,
		},
		"conversion": {
			Name: "Conversion",
			View: d.conversionView,
		},
		"playlistMatching": {
			Name: "Playlist Matching",
			View: d.playlistMatchingView,
		},
	}
}

func (d *Data) getOperationIndex() map[string][]string {
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
