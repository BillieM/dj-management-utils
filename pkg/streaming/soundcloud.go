package streaming

import (
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/billiem/seren-management/pkg/helpers"
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
	playlistUrl string
}

type GetSoundCloudPlaylistOpts struct {
	playlistUrl string
}

type soundcloudRequestParams struct {
	resource string
	params   url.Values
}

type soundcloudResponse struct {
}

func (o GetSoundCloudPlaylistOpts) Build(cfg helpers.Config) StreamingPlatform {
	return &SoundCloud{
		requestUrl: o.requestUrl,
	}
}

func (s SoundCloud) makeRequest(p soundcloudRequestParams) error {
	baseURL := "https://api-v2.soundcloud.com"

	u, _ := url.ParseRequestURI(baseURL)
	u.Path = p.resource
	u.RawQuery = p.params.Encode()
	urlStr := fmt.Sprintf("%v", u) // "http://example.com/path?param1=value1&param2=value2"

	fmt.Println(urlStr)

	resp, err := http.Get(urlStr)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	// var data YourStruct // Replace YourStruct with the appropriate struct type for unmarshaling
	// err = json.Unmarshal(body, &data)
	// if err != nil {
	// 	return err
	// }

	// Process the unmarshaled data

	return nil
}

func (s SoundCloud) GetPlaylist() error {

	url := "https://soundcloud.com/serrene/sets/not-chill-but-u-know/s-HfjKTSg2san?si=d0884248a6024ab7a5567e86a732fb0f"

	resp, err := http.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return err
	}

	fmt.Println(string(body))

	return nil
}

/*
SoundCloudPlaylist implements the StreamingPlaylist interface for SoundCloud

It holds the data for a SoundCloud playlist
*/
type SoundCloudPlaylist struct {
	playlistName string
	tracks       []SoundCloudTrack
}

type SoundCloudTrack struct {
	TrackID   int
	TrackName string
}
