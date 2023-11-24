package streaming

import (
	"errors"
	"io"
	"net/http"
)

/*
cool things:

auto tagging suggestion API
https://api-v2.soundcloud.com/tags/suggested/soundcloud:playlists:1720019865?client_id=1ZRkRXa5klyxfeCePlMbkWl1xNzz1Bu3&limit=10&offset=0&linked_partitioning=1&app_version=1700828706&app_locale=en
	- requires client_id

tracks api

https://api-v2.soundcloud.com/tracks?ids=1032527881%2C1164496324%2C1380759907%2C1399327981%2C1426593295%2C1436426464%2C1476413599%2C1491581392%2C1552610359%2C1598446053%2C241418628%2C253195240%2C762760012%2C950173294%2C978990424&client_id=1ZRkRXa5klyxfeCePlMbkWl1xNzz1Bu3&%5Bobject%20Object%5D=&app_version=1700828706&app_locale=en
	- requires client_id
*/

type SoundCloud struct {
	ClientID string
}

func (s SoundCloud) GetSoundCloudPlaylist(playlistUrl string) (SoundCloudPlaylist, error) {

	url := "https://soundcloud.com/serrene/sets/not-chill-but-u-know/s-HfjKTSg2san?si=d0884248a6024ab7a5567e86a732fb0f"

	resp, err := http.Get(url)

	if err != nil {
		return SoundCloudPlaylist{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return SoundCloudPlaylist{}, err
	}

	hydratableStr, err := extractSCHydrationString(string(body))

	if err != nil {
		return SoundCloudPlaylist{}, err
	}

	h := Hydratable{}
	err = h.UnmarshalJSON([]byte(hydratableStr))

	if err != nil {
		return SoundCloudPlaylist{}, err
	}

	err = h.Playlist.CompleteTracks()

	return SoundCloudPlaylist{}, nil
}

/*
completePlaylistTracks adds missing data to tracks in a SoundCloudPlaylist struct

This is needed as soundcloud only returns IDs for any tracks beyond the first 5
*/
func (playlist *SoundCloudPlaylist) completePlaylistTracks() error {

	for _, track := range playlist.Tracks {


	return nil
}

func (track TrackElement) check() (bool, error) {
	if track.ID == "" {
		return false, helpers.ErrTrackMissingID	
	}
}

func getTracks(ids []string) ([]TrackElement, error) {

	return nil, nil
}
