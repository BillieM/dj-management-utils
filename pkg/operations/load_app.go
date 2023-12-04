package operations

import (
	"context"
	"fmt"

	"github.com/billiem/seren-management/pkg/data"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/streaming"
)

/*
TODO: handle all these with step handlers
*/

func (e *OpEnv) CheckLocalPaths() {

	var brokenPathChanged []streaming.SoundCloudTrack

	tracks, err := e.SerenDB.ListSoundCloudTracksHasLocalPath(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, track := range tracks {

		t := streaming.SoundCloudTrack{}
		t.LoadFromDB(track)

		fileExists := helpers.DoesFileExist(t.LocalPath)

		if fileExists == t.LocalPathBroken {
			t.LocalPathBroken = !fileExists

			brokenPathChanged = append(brokenPathChanged, t)
		}
	}

	if len(brokenPathChanged) == 0 {
		fmt.Println("No tracks have changed broken status")
		return
	}

	fmt.Printf("Found %d tracks with changed broken status\n", len(brokenPathChanged))

	dataT := make([]data.SoundcloudTrack, len(brokenPathChanged))

	for i, t := range brokenPathChanged {
		dataT[i] = t.ToDB()
	}

	e.SerenDB.TxUpsertSoundCloudTracks(dataT)
}

func (e *OpEnv) IndexCollections() {

}

func (e *OpEnv) IndexLocalFolders() {

}
