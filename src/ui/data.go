package ui

import "fyne.io/fyne/v2"

type Operation struct {
	Name string
	View func(w fyne.Window) fyne.CanvasObject
}

var (
	Operations = map[string]Operation{
		"home": {
			Name: "Home",
			View: homeView,
		},
		"stems": {
			Name: "Stems",
			View: stemsView,
		},
		"separateTrack": {
			Name: "Separate Track",
			View: separateTrackView,
		},
		"separateFolder": {
			Name: "Separate Folder",
			View: separateFolderView,
		},
		"separateCollection": {
			Name: "Separate Collection",
			View: separateCollectionView,
		},
		"mp3s": {
			Name: "Convert MP3s",
			View: convertMp3sView,
		},
		"convertSingleMp3": {
			Name: "Convert Single",
			View: convertSingleMp3View,
		},
		"convertFolderMp3": {
			Name: "Convert Folder",
			View: convertFolderMp3View,
		},
		"convertCollectionMp3": {
			Name: "Convert Collection",
			View: convertCollectionMp3View,
		},
		"tags": {
			Name: "Process Tags",
			View: tagsView,
		},
		"rereadTags": {
			Name: "Reread Tags",
			View: rereadTagsView,
		},
		"cleanTags": {
			Name: "Clean Tags",
			View: cleanTagsView,
		},
		"conversion": {
			Name: "Conversion",
			View: conversionView,
		},
		"playlistMatching": {
			Name: "Playlist Matching",
			View: playlistMatchingView,
		},
	}

	OperationIndex = map[string][]string{
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
)
