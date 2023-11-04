package ui

import (
	"fyne.io/fyne/v2"
	"github.com/billiem/seren-management/src/helpers"
)

type Operation struct {
	Name string
	View func(w fyne.Window) fyne.CanvasObject
}

type Data struct {
	*helpers.Config
	*State
	TmpConfig      *helpers.Config
	Operations     map[string]Operation
	OperationIndex map[string][]string
}

type State struct {
	settingsAlreadyOpen bool
	processing          bool
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
			View: d.separateTrackView,
		},
		"separateFolder": {
			Name: "Separate Folder",
			View: d.separateFolderView,
		},
		"separateCollection": {
			Name: "Separate Collection",
			View: d.separateCollectionView,
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
