package operations

import (
	"fmt"

	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
)

/*
TODO: handle all these with step handlers
*/

func (e *OpEnv) CheckLocalPaths() {

	var brokenPathChanged []database.SoundCloudTrack

	tracks, err := e.SerenDB.GetSoundCloudTracksWithLocalPaths()
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range tracks {
		fileExists := helpers.DoesFileExist(track.LocalPath)

		if fileExists == track.LocalPathBroken {
			track.LocalPathBroken = !fileExists

			brokenPathChanged = append(brokenPathChanged, track)
		}
	}

	if len(brokenPathChanged) == 0 {
		fmt.Println("No tracks have changed broken status")
		return
	}

	fmt.Printf("Found %d tracks with changed broken status\n", len(brokenPathChanged))

	e.SerenDB.SaveSoundCloudTracks(brokenPathChanged)
}

func (e *OpEnv) IndexCollections() {

}

func (e *OpEnv) IndexLocalFolders() {

}
