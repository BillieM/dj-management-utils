package streaming

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/Southclaws/fault"
	"github.com/Southclaws/fault/fctx"
	"github.com/Southclaws/fault/fmsg"
	"github.com/billiem/seren-management/pkg/data"
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

type SoundCloudPlaylist struct {
	ExternalID   int64
	Name         string
	SearchUrl    string
	PermalinkUrl string

	Tracks []SoundCloudTrack

	NumTracks int
}

func (p SoundCloudPlaylist) String() string {
	return fmt.Sprintf(
		"%v: %s",
		p.ExternalID,
		p.Name,
	)
}

func (p *SoundCloudPlaylist) loadFromHydratable(hp HydratableSoundCloudPlaylist) {

	tracks := []SoundCloudTrack{}

	for _, track := range hp.Tracks {
		var newTrack SoundCloudTrack
		newTrack.loadFromHydratable(track)
		tracks = append(tracks, newTrack)
	}

	p.ExternalID = hp.ID
	p.Name = hp.Title
	p.PermalinkUrl = hp.PermalinkURL
	p.Tracks = tracks
}

/*
LoadFromDB loads a SoundCloudPlaylist from a data.SoundcloudPlaylist and a slice of data.SoundcloudTrack,
this exists in order to properly map fields between the database and the struct
*/
func (p *SoundCloudPlaylist) LoadFromDB(dp data.SoundcloudPlaylist, tracks []data.SoundcloudTrack) {

	for _, track := range tracks {
		var newTrack SoundCloudTrack
		newTrack.LoadFromDB(track)
		p.Tracks = append(p.Tracks, newTrack)
	}

	p.ExternalID = dp.ExternalID.Int64
	p.Name = dp.Name.String
	p.SearchUrl = dp.SearchUrl.String
	p.PermalinkUrl = dp.PermalinkUrl.String
}

func (p *SoundCloudPlaylist) ToDB() (data.SoundcloudPlaylist, []data.SoundcloudTrack) {
	dataP := data.SoundcloudPlaylist{
		ExternalID:   sql.NullInt64{Valid: true, Int64: p.ExternalID},
		Name:         sql.NullString{Valid: true, String: p.Name},
		SearchUrl:    sql.NullString{Valid: true, String: p.SearchUrl},
		PermalinkUrl: sql.NullString{Valid: true, String: p.PermalinkUrl},
	}

	dataTracks := []data.SoundcloudTrack{}

	for _, track := range p.Tracks {
		dataTracks = append(dataTracks, track.ToDB())
	}

	return dataP, dataTracks
}

type SoundCloudTrack struct {
	ExternalID          int64  `gorm:"uniqueIndex"`
	Name                string // change to title
	PermalinkUrl        string
	PurchaseTitle       string
	PurchaseURL         string
	HasDownloadsLeft    bool
	Genre               string
	ArtworkURL          string
	TagList             string
	PublisherArtist     string // users/ artists use relationships
	SoundCloudUser      string
	LocalPath           string
	LocalPathBroken     bool
	RemovedFromPlaylist bool

	Playlists []SoundCloudPlaylist `gorm:"many2many:playlist_tracks;"`
}

func (t SoundCloudTrack) String() string {
	return fmt.Sprintf(
		"%v: %s",
		t.ExternalID,
		t.Name,
	)
}

func (t *SoundCloudTrack) loadFromHydratable(ht HydratableSoundCloudTrack) {

	t.ExternalID = ht.ID
	t.Name = ht.Title
	t.PermalinkUrl = ht.PermalinkURL
	t.PurchaseTitle = ht.PurchaseTitle
	t.PurchaseURL = ht.PurchaseURL
	t.HasDownloadsLeft = ht.HasDownloadsLeft
	t.Genre = ht.Genre
	t.ArtworkURL = ht.ArtworkURL
	t.TagList = ht.TagList
	t.PublisherArtist = ht.PublisherMetadata.Artist
	t.SoundCloudUser = ht.User.Username
}

/*
LoadFromDB loads a SoundCloudTrack from a data.SoundcloudTrack,
this exists in order to properly map fields between the database and the struct
*/
func (t *SoundCloudTrack) LoadFromDB(dt data.SoundcloudTrack) {
	t.ExternalID = dt.ExternalID.Int64
	t.Name = dt.Name.String
	t.PermalinkUrl = dt.PermalinkUrl.String
	t.PurchaseTitle = dt.PurchaseTitle.String
	t.PurchaseURL = dt.PurchaseUrl.String
	t.HasDownloadsLeft = dt.HasDownloadsLeft.Bool
	t.Genre = dt.Genre.String
	t.ArtworkURL = dt.ArtworkUrl.String
	t.TagList = dt.TagList.String
	t.PublisherArtist = dt.PublisherArtist.String
	t.SoundCloudUser = dt.SoundCloudUser.String
	t.LocalPath = dt.LocalPath.String
	t.LocalPathBroken = dt.LocalPathBroken.Bool
	t.RemovedFromPlaylist = dt.RemovedFromPlaylist.Bool
}

func (t *SoundCloudTrack) ToDB() data.SoundcloudTrack {
	return data.SoundcloudTrack{
		ExternalID:          sql.NullInt64{Valid: true, Int64: t.ExternalID},
		Name:                sql.NullString{Valid: true, String: t.Name},
		PermalinkUrl:        sql.NullString{Valid: true, String: t.PermalinkUrl},
		PurchaseTitle:       sql.NullString{Valid: true, String: t.PurchaseTitle},
		PurchaseUrl:         sql.NullString{Valid: true, String: t.PurchaseURL},
		HasDownloadsLeft:    sql.NullBool{Valid: true, Bool: t.HasDownloadsLeft},
		Genre:               sql.NullString{Valid: true, String: t.Genre},
		ArtworkUrl:          sql.NullString{Valid: true, String: t.ArtworkURL},
		TagList:             sql.NullString{Valid: true, String: t.TagList},
		PublisherArtist:     sql.NullString{Valid: true, String: t.PublisherArtist},
		SoundCloudUser:      sql.NullString{Valid: true, String: t.SoundCloudUser},
		LocalPath:           sql.NullString{Valid: true, String: t.LocalPath},
		LocalPathBroken:     sql.NullBool{Valid: true, Bool: t.LocalPathBroken},
		RemovedFromPlaylist: sql.NullBool{Valid: true, Bool: t.RemovedFromPlaylist},
	}
}

func (s SoundCloud) GetSoundCloudPlaylist(ctx context.Context, playlistUrl string) (SoundCloudPlaylist, error) {

	ctx = fctx.WithMeta(
		ctx,
		"playlist_url", playlistUrl,
		"client_id", s.ClientID,
	)

	resp, err := http.Get(playlistUrl)

	if err != nil {
		return SoundCloudPlaylist{}, fault.Wrap(
			err,
			fctx.With(ctx),
			fmsg.With("Error making request to get SoundCloud playlist"),
		)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return SoundCloudPlaylist{}, fault.Wrap(
			err,
			fctx.With(ctx),
			fmsg.With("Error reading response body"),
		)
	}

	hydratableStr, err := extractSCHydrationString(string(body))

	if err != nil {
		return SoundCloudPlaylist{}, fault.Wrap(
			err,
			fctx.With(ctx),
			fmsg.With("Error extracting hydration string"),
		)
	}

	h := Hydratable{}
	err = h.UnmarshalJSON([]byte(hydratableStr))

	if err != nil {
		return SoundCloudPlaylist{}, fault.Wrap(
			err,
			fctx.With(ctx),
			fmsg.With("Error unmarshalling hydratable string"),
		)
	}

	if h.Playlist.ID == 0 {
		return SoundCloudPlaylist{}, fault.Wrap(
			fault.New("Missing playlist ID"),
			fctx.With(ctx),
		)
	}

	err = s.completeTracks(ctx, &h.Playlist)

	if err != nil {
		return SoundCloudPlaylist{}, fault.Wrap(
			err,
			fctx.With(ctx),
			fmsg.With("Error completing tracks"),
		)
	}

	p := SoundCloudPlaylist{}
	p.loadFromHydratable(h.Playlist)

	return p, nil
}

/*
completePlaylistTracks adds missing data to tracks in a SoundCloudPlaylist struct

This is needed as soundcloud only returns IDs for any tracks beyond the first 5
*/
func (s SoundCloud) completeTracks(ctx context.Context, p *HydratableSoundCloudPlaylist) error {

	okayTracks := []HydratableSoundCloudTrack{}
	trackIdsToRequest := []int64{}

	for _, track := range p.Tracks {
		needToRequestTrack, err := track.check()

		if err != nil {
			return fault.Wrap(
				err,
				fmsg.With(fmt.Sprintf("Error checking track %d", track.ID)),
			)
		}

		if !needToRequestTrack {
			trackIdsToRequest = append(trackIdsToRequest, track.ID)
			continue
		}

		okayTracks = append(okayTracks, track)
	}

	remainingTracks, err := s.getRemainingTracks(ctx, trackIdsToRequest)

	if err != nil {
		return fault.Wrap(
			err,
			fmsg.With("Error getting remaining tracks"),
		)
	}

	p.Tracks = append(okayTracks, remainingTracks...)

	return nil
}

func (track HydratableSoundCloudTrack) check() (bool, error) {
	if track.ID == 0 {
		return false, helpers.ErrTrackMissingID
	}

	// Track is missing title, so we need to request it
	if track.Title == "" {
		return false, nil
	}

	return true, nil
}

func (s SoundCloud) getRemainingTracks(ctx context.Context, ids []int64) ([]HydratableSoundCloudTrack, error) {

	ctx, cancel := context.WithCancelCause(ctx)

	trackIDChan := pipeline.Emit(ids...)

	// TODO: figure out how big this array can be
	tracksOut := pipeline.ProcessBatchConcurrently(ctx, 2, 50, time.Second*15, pipeline.NewProcessor(func(ctx context.Context, ids []int64) ([]HydratableSoundCloudTrack, error) {
		trackArr, err := s.makeSoundCloudTracksRequest(ids)
		if err != nil {
			return nil, fault.Wrap(
				err,
				fmsg.With("Error making SoundCloud tracks request"),
			)
		}
		return trackArr, nil
	}, func(ids []int64, err error) {
		if err != nil {
			cancel(err)
		}
	}), trackIDChan)

	outTracks := []HydratableSoundCloudTrack{}

	for track := range tracksOut {
		outTracks = append(outTracks, track)
	}

	if ctx.Err() != nil {
		return nil, fault.Wrap(
			context.Cause(ctx),
			fmsg.With("Error making batch requests for remaining tracks"),
		)
	}

	return outTracks, ctx.Err()
}

func (s SoundCloud) makeSoundCloudTracksRequest(ids []int64) ([]HydratableSoundCloudTrack, error) {

	req, err := http.NewRequest("GET", "https://api-v2.soundcloud.com/tracks", nil)

	if err != nil {
		return nil, fault.Wrap(
			err,
			fmsg.With("Error creating request"),
		)
	}

	q := req.URL.Query()

	q.Add("client_id", s.ClientID)
	q.Add("app_locale", "en")
	q.Add("ids", helpers.Int64ArrayToJoinedString(ids))
	q.Add("app_version", "1700828706")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fault.Wrap(
			err,
			fmsg.With("Error making request"),
		)
	}

	if resp.StatusCode != 200 {
		return nil, fault.New(fmt.Sprintf(
			"Error making request to get SoundCloud tracks, status code %d", resp.StatusCode,
		))
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fault.Wrap(
			err,
			fmsg.With("Error reading response body"),
		)
	}

	var tracks []HydratableSoundCloudTrack
	err = json.Unmarshal(body, &tracks)
	if err != nil {
		return nil, fault.Wrap(
			err,
			fmsg.With("Error unmarshalling response body to tracks"),
		)
	}

	return tracks, nil
}

/*
DownloadFile downloads a file from SoundCloud

returns the path to the downloaded file if successful, otherwise returns an error

The file extension is gathered from the Content-Disposition header
*/
func (s SoundCloud) DownloadFile(dirPath string, id int64) (string, error) {

	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://api-v2.soundcloud.com/tracks/%d/download", id),
		nil,
	)

	if err != nil {
		return "", err
	}

	q := req.URL.Query()

	q.Add("client_id", s.ClientID)
	q.Add("app_locale", "en")
	q.Add("app_version", "1700828706")

	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var uriMap map[string]string

	err = json.Unmarshal(body, &uriMap)

	if err != nil {
		return "", err
	}

	val, ok := uriMap["redirectUri"]

	if !ok {
		return "", helpers.ErrMissingRedirectURI
	}

	resp, err = http.Get(val)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	contentDisposition := resp.Header.Get("Content-Disposition")

	filename, err := helpers.GetFileNameFromContentDisposition(contentDisposition)

	if err != nil {
		return "", err
	}

	err = helpers.CreateDirIfNotExists(dirPath)

	if err != nil {
		return "", err
	}

	filePath := helpers.JoinFilepathToSlash(dirPath, filename)

	f, err := os.Create(filePath)

	if err != nil {
		return "", err
	}

	defer f.Close()

	_, err = io.Copy(f, resp.Body)

	if err != nil {
		return "", err
	}

	return filePath, nil
}

func GetAndSetSoundCloudClientID(cfg helpers.Config) (string, error) {
	if cfg.SoundCloudClientID != "" {
		return cfg.SoundCloudClientID, nil
	}

	clientID, err := GenerateSoundCloudClientID()

	if err != nil {
		return "", fault.Wrap(
			err,
			fmsg.With("Error generating SoundCloud client_id"),
		)
	}

	cfg.SoundCloudClientID = clientID

	return clientID, nil
}

func GenerateSoundCloudClientID() (string, error) {

	resp, err := http.Get("https://www.soundcloud.com")

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", err
	}

	assetMatches := helpers.RegexAllSubmatches(string(body), `src="(https[^"]+.js)"`)

	for _, match := range assetMatches {

		assetResp, err := http.Get(match[1])

		if err != nil {
			continue
		}

		defer assetResp.Body.Close()

		assetBody, err := io.ReadAll(assetResp.Body)

		if err != nil {
			continue
		}

		clientIDMatches := helpers.RegexAllSubmatches(string(assetBody), `client_id[=:]\s*"*([a-zA-Z0-9]{32})`)

		if len(clientIDMatches) > 0 {
			return clientIDMatches[0][1], nil
		}
	}

	return "", fault.New("Error generating SoundCloud client_id")
}
