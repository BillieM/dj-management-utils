package streaming

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/billiem/seren-management/pkg/database"
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/deliveryhero/pipeline/v2"
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

func (s SoundCloud) GetSoundCloudPlaylist(ctx context.Context, playlistUrl string) (database.SoundCloudPlaylist, error) {

	resp, err := http.Get(playlistUrl)

	if err != nil {
		return database.SoundCloudPlaylist{}, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return database.SoundCloudPlaylist{}, err
	}

	hydratableStr, err := extractSCHydrationString(string(body))

	if err != nil {
		return database.SoundCloudPlaylist{}, err
	}

	h := Hydratable{}
	err = h.UnmarshalJSON([]byte(hydratableStr))

	if err != nil {
		return database.SoundCloudPlaylist{}, err
	}

	if h.Playlist.ID == 0 {
		return database.SoundCloudPlaylist{}, helpers.ErrRequestingPlaylist
	}

	err = s.completeTracks(ctx, &h.Playlist)

	if err != nil {
		return database.SoundCloudPlaylist{}, err
	}

	return h.Playlist.ToDB(), nil
}

/*
completePlaylistTracks adds missing data to tracks in a SoundCloudPlaylist struct

This is needed as soundcloud only returns IDs for any tracks beyond the first 5
*/
func (s SoundCloud) completeTracks(ctx context.Context, p *SoundCloudPlaylist) error {

	okayTracks := []TrackElement{}
	trackIdsToRequest := []int64{}

	for _, track := range p.Tracks {
		ok, err := track.check()

		if err != nil {
			return err
		}

		if !ok {
			trackIdsToRequest = append(trackIdsToRequest, track.ID)
			continue
		}

		okayTracks = append(okayTracks, track)
	}

	remainingTracks, err := s.getRemainingTracks(ctx, trackIdsToRequest)

	if err != nil {
		return err
	}

	p.Tracks = append(okayTracks, remainingTracks...)

	return nil
}

func (track TrackElement) check() (bool, error) {
	if track.ID == 0 {
		return false, helpers.ErrTrackMissingID
	}

	// Track is missing title, so we need to request it
	if track.Title == "" {
		return false, nil
	}

	return true, nil
}

func (s SoundCloud) getRemainingTracks(ctx context.Context, ids []int64) ([]TrackElement, error) {

	trackIDChan := pipeline.Emit(ids...)

	// TODO: figure out how big this array can be
	tracksOut := pipeline.ProcessBatchConcurrently(ctx, 2, 50, time.Second*15, pipeline.NewProcessor(func(ctx context.Context, ids []int64) ([]TrackElement, error) {
		trackArr, err := s.makeSoundCloudTracksRequest(ids)
		if err != nil {
			return nil, err
		}
		return trackArr, nil
	}, func(ids []int64, err error) {
		if err != nil {
			// TODO: make this not just panic...
			panic(err)
		}
	}), trackIDChan)

	outTracks := []TrackElement{}

	for track := range tracksOut {
		outTracks = append(outTracks, track)
	}

	return outTracks, nil
}

func (s SoundCloud) makeSoundCloudTracksRequest(ids []int64) ([]TrackElement, error) {

	req, err := http.NewRequest("GET", "https://api-v2.soundcloud.com/tracks", nil)

	if err != nil {
		return nil, err
	}

	q := req.URL.Query()

	q.Add("client_id", s.ClientID)
	q.Add("app_locale", "en")
	q.Add("ids", helpers.Int64ArrayToJoinedString(ids))
	q.Add("app_version", "1700828706")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var tracks []TrackElement
	err = json.Unmarshal(body, &tracks)
	if err != nil {
		return nil, err
	}

	return tracks, nil
}

func (s SoundCloud) DownloadFile(dirPath string, id int64) error {

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api-v2.soundcloud.com/tracks/%d/download", id),
		nil,
	)

	if err != nil {
		return err
	}

	q := req.URL.Query()

	q.Add("client_id", s.ClientID)
	q.Add("app_locale", "en")
	q.Add("app_version", "1700828706")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var uriMap map[string]string

	err = json.Unmarshal(body, &uriMap)

	if err != nil {
		return err
	}

	val, ok := uriMap["redirectUri"]

	if !ok {
		return helpers.ErrMissingRedirectURI
	}

	resp, err = http.Get(val)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	contentDisposition := resp.Header.Get("Content-Disposition")

	filename, err := helpers.GetFileNameFromContentDisposition(contentDisposition)

	if err != nil {
		return err
	}

	err = helpers.CreateDirIfNotExists(dirPath)

	if err != nil {
		return err
	}

	f, err := os.Create(helpers.JoinFilepathToSlash(dirPath, filename))

	if err != nil {
		return err
	}

	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	if err != nil {
		return err
	}

	return nil
}
