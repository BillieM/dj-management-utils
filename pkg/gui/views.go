package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/billiem/seren-management/pkg/gui/iwidget"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/operations"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
guiView is a struct that contains the name of the view and the function that returns the view

Views define the content for the main content area of the application
*/
type guiView struct {
	name   string
	render func() fyne.CanvasObject
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
		"sync": {
			name:   "Playlist Matching",
			render: e.syncView,
		},
		"syncSoundCloud": {
			name:   "SoundCloud",
			render: e.syncSoundCloudView,
		},
		"syncSpotify": {
			name:   "Spotify",
			render: e.syncSpotifyView,
		},
	}
}

func (e *guiEnv) getViewIndex() map[string][]string {
	return map[string][]string{
		"": {"home", "stems", "mp3s", "tags", "conversion", "sync"},
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
		"sync": {
			"syncSoundCloud",
			"syncSpotify",
		},
	}
}

/*
setMainContent sets the main content of the window to the provided content

Called on tab change on the main menu
*/
func (e *guiEnv) setMainContent(contentStack *fyne.Container, view guiView) {

	labelContainer := container.NewVBox(widget.NewLabel(view.name), widget.NewSeparator())

	contentContainer := container.NewBorder(labelContainer, nil, nil, nil, view.render())

	contentStack.Objects = []fyne.CanvasObject{contentContainer}
	contentStack.Refresh()
}

func (e *guiEnv) homeView() fyne.CanvasObject {
	return widget.NewLabel("Welcome to Seren Library Management!")
}

func (e *guiEnv) stemsView() fyne.CanvasObject {
	content := widget.NewLabel("Contains a selection of utilities for separating stems from audio files.")

	return container.NewVBox(content)
}

func (e *guiEnv) separateSingleStemView() fyne.CanvasObject {
	ok, canvas := e.checkConfig([]func() (bool, string){})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.SeparateSingleStemOpts{}

	startFunc := func() {
		trackOperation := e.prepareTrackOperation()
		processContainerOuter.Add(trackOperation.runningOperation)
		go trackOperation.opEnv.SeparateSingleStem(
			trackOperation.ctx,
			opts,
		)
	}

	startButton := widget.NewButton("Separate stem", startFunc)
	startButton.Disable()

	trackPathCanvas := e.openFileCanvas("Track Path", &opts.InFilePath, []string{".wav", ".mp3"}, func() { enableBtnIfOptsOkay(opts, startButton) })
	stemTypeSelect := buildStemTypeSelect(&opts.Type, func() { enableBtnIfOptsOkay(opts, startButton) })

	return container.NewBorder(
		container.NewVBox(
			container.NewVBox(
				trackPathCanvas,
				stemTypeSelect,
			),
			startButton,
		), nil, nil, nil,
		processContainerOuter,
	)
}

func (e *guiEnv) separateFolderStemView() fyne.CanvasObject {
	ok, canvas := e.checkConfig([]func() (bool, string){})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.SeparateFolderStemOpts{}

	startFunc := func() {
		trackOperation := e.prepareTrackOperation()
		processContainerOuter.Add(trackOperation.runningOperation)
		go trackOperation.opEnv.SeparateFolderStem(
			trackOperation.ctx,
			opts,
		)
	}

	startButton := widget.NewButton("Separate folder", startFunc)
	startButton.Disable()

	trackPathCanvas := e.openDirCanvas("Folder Path", &opts.InDirPath, func() { enableBtnIfOptsOkay(opts, startButton) })
	stemTypeSelect := buildStemTypeSelect(&opts.Type, func() { enableBtnIfOptsOkay(opts, startButton) })

	return container.NewBorder(
		container.NewVBox(
			container.NewVBox(
				trackPathCanvas,
				stemTypeSelect,
			),
			startButton,
		), nil, nil, nil,
		processContainerOuter,
	)
}

func (e *guiEnv) separateCollectionStemView() fyne.CanvasObject {
	return widget.NewLabel("separateCollectionView")
}

/*
Convert Mp3s Section
*/

// convertMp3sView returns the view for the convert mp3s info section
func (e *guiEnv) convertMp3sView() fyne.CanvasObject {
	return widget.NewLabel("convertMp3sView")
}

// convertSingleMp3View returns the view for the convert single mp3 operation
func (e *guiEnv) convertSingleMp3View() fyne.CanvasObject {
	ok, canvas := e.checkConfig([]func() (bool, string){})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.ConvertSingleMp3Opts{}

	startFunc := func() {
		trackOperation := e.prepareTrackOperation()
		processContainerOuter.Add(trackOperation.runningOperation)
		go trackOperation.opEnv.ConvertSingleMp3(
			trackOperation.ctx,
			opts,
		)
	}

	startButton := widget.NewButton("Convert to mp3", startFunc)
	startButton.Disable()

	trackPathCanvas := e.openFileCanvas("Track Path", &opts.InFilePath, []string{".wav", ".flac"}, func() { startButton.Enable() })

	return container.NewBorder(
		container.NewVBox(
			container.NewVBox(
				trackPathCanvas,
			),
			startButton,
		), nil, nil, nil,
		processContainerOuter,
	)
}

// convertFolderMp3View returns the view for the convert folder mp3 operation
func (e *guiEnv) convertFolderMp3View() fyne.CanvasObject {
	ok, canvas := e.checkConfig([]func() (bool, string){})

	if !ok {
		return canvas
	}

	processContainerOuter := container.NewStack()

	opts := operations.ConvertFolderMp3Opts{}

	startFunc := func() {
		trackOperation := e.prepareTrackOperation()
		processContainerOuter.Add(trackOperation.runningOperation)
		go trackOperation.opEnv.ConvertFolderMp3(
			trackOperation.ctx,
			opts,
		)
	}

	startButton := widget.NewButton("Convert folder to mp3", startFunc)
	startButton.Disable()

	trackPathCanvas := e.openDirCanvas("Folder Path", &opts.InDirPath, func() { startButton.Enable() })

	return container.NewBorder(
		container.NewVBox(
			container.NewVBox(
				trackPathCanvas,
			),
			startButton,
		), nil, nil, nil,
		processContainerOuter,
	)
}

// convertCollectionMp3View returns the view for the convert collection mp3 operation
func (e *guiEnv) convertCollectionMp3View() fyne.CanvasObject {

	path := widget.NewLabel("collection path: " + e.Config.TraktorCollectionPath)
	btn := widget.NewButton("do collection things :D", func() {
		// collection.ReadCollection(e.Config.TraktorCollectionPath)
	})

	return container.NewBorder(
		nil, nil, nil, nil,
		container.NewVBox(
			path,
			btn,
		),
	)
}

func (e *guiEnv) tagsView() fyne.CanvasObject {
	return widget.NewLabel("tagsView")
}

func (e *guiEnv) rereadTagsView() fyne.CanvasObject {
	return widget.NewLabel("rereadTagsView")
}

func (e *guiEnv) cleanTagsView() fyne.CanvasObject {
	return widget.NewLabel("cleanTagsView")
}

func (e *guiEnv) conversionView() fyne.CanvasObject {
	return widget.NewLabel("conversionView")
}

func (e *guiEnv) syncView() fyne.CanvasObject {
	return widget.NewLabel("syncView")
}

func (e *guiEnv) syncSoundCloudView() fyne.CanvasObject {

	playlistBindVals := iwidget.PlaylistBindingList{
		Base:  e.getWidgetBase(),
		Items: []*iwidget.PlaylistBindingItem{},
	}

	// Create the list widget that will display the playlists, and bind it to the playlistBindVals
	// this will cause the list to be updated when the bind value changes (i.e. when the playlists are loaded)
	playlistsList := widget.NewListWithData(
		&playlistBindVals,
		func() fyne.CanvasObject {
			return iwidget.NewPlaylist(
				func(playlistData streaming.SoundCloudPlaylist) {
					if e.guiState.busy {
						e.showErrorDialog(helpers.ErrBusyPleaseFinishFirst)
						return
					}
					e.openPlaylistPopup(playlistData)
				},
			)
		},
		func(i binding.DataItem, o fyne.CanvasObject) {

			playlistWidget := o.(*iwidget.Playlist)
			playlistBindingItem := i.(*iwidget.PlaylistBindingItem)

			playlistWidget.UpdateFromData(playlistBindingItem)
		},
	)

	// Load the playlists into the view in the background
	// This should be quick as it only requires a database query
	// Also show a loading screen, which will be hidden when the playlists are loaded
	loading := iwidget.NewViewLoading("Loading SoundCloud playlists...")
	go func() {
		err := e.loadSoundCloudPlaylists(&playlistBindVals)

		if err != nil {
			e.showErrorDialog(err)
			return
		}

		playlistsList.Refresh()
		loading.Hide()
	}()

	// Function to call when the add playlist button is clicked
	addPlaylistCanvas := iwidget.NewAddPlaylist(
		e.getAddPlaylistCallback(
			&playlistBindVals,
			func() {
				playlistsList.Refresh()
				playlistsList.ScrollToBottom()
			},
		),
	)

	return container.NewStack(
		container.NewBorder(
			nil, addPlaylistCanvas, nil, nil, playlistsList,
		),
		loading,
	)
}

func (e *guiEnv) syncSpotifyView() fyne.CanvasObject {
	return widget.NewLabel("syncView")
}
