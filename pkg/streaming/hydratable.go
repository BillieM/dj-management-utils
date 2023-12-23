package streaming

import (
	"encoding/json"

	"github.com/billiem/seren-management/pkg/helpers"
)

func (h *Hydratable) UnmarshalJSON(p []byte) error {
	var tmp []any
	if err := json.Unmarshal(p, &tmp); err != nil {
		return err
	}

	for _, v := range tmp {
		m := v.(map[string]any)
		hydratable, hydratableExists := m["hydratable"]
		data, dataExists := m["data"]

		if !hydratableExists || !dataExists {
			continue
		}

		jsonBytes, err := json.Marshal(data)

		if err != nil {
			return err
		}

		switch hydratable {
		case "user":
			var user HydratableSoundCloudUser

			err := json.Unmarshal(jsonBytes, &user)
			if err != nil {
				return err
			}
			h.User = user
		case "playlist":
			var playlist HydratableSoundCloudPlaylist
			err := json.Unmarshal(jsonBytes, &playlist)
			if err != nil {
				return err
			}
			h.Playlist = playlist
		}

	}

	return nil
}

/*
extractSCHydrationString extracts data from the returned HTML of a SoundCloud request

Little bit of a hacky way to do it but it works for now given the API limitations
*/
func extractSCHydrationString(body string) (string, error) {
	res := helpers.RegexAllSubmatches(body, `(?m)window.__sc_hydration = (\[.*\])`)

	if len(res) == 0 {
		return "", helpers.ErrExtractingHydrationString
	}

	if len(res[0]) < 2 {
		return "", helpers.ErrExtractingHydrationString
	}

	return res[0][1], nil
}

type Hydratable struct {
	User     HydratableSoundCloudUser
	Playlist HydratableSoundCloudPlaylist
}

type HydratableSoundCloudUser struct {
	AvatarURL            string                          `json:"avatar_url"`
	City                 string                          `json:"city"`
	CommentsCount        int64                           `json:"comments_count"`
	CountryCode          string                          `json:"country_code"`
	CreatedAt            string                          `json:"created_at"`
	CreatorSubscriptions []HydratableCreatorSubscription `json:"creator_subscriptions"`
	CreatorSubscription  HydratableCreatorSubscription   `json:"creator_subscription"`
	Description          string                          `json:"description"`
	FollowersCount       int64                           `json:"followers_count"`
	FollowingsCount      int64                           `json:"followings_count"`
	FirstName            string                          `json:"first_name"`
	FullName             string                          `json:"full_name"`
	GroupsCount          int64                           `json:"groups_count"`
	ID                   int64                           `json:"id"`
	Kind                 string                          `json:"kind"`
	LastModified         string                          `json:"last_modified"`
	LastName             string                          `json:"last_name"`
	LikesCount           int64                           `json:"likes_count"`
	PlaylistLikesCount   int64                           `json:"playlist_likes_count"`
	Permalink            string                          `json:"permalink"`
	PermalinkURL         string                          `json:"permalink_url"`
	PlaylistCount        int64                           `json:"playlist_count"`
	RepostsCount         int64                           `json:"reposts_count"`
	TrackCount           int64                           `json:"track_count"`
	URI                  string                          `json:"uri"`
	Urn                  string                          `json:"urn"`
	Username             string                          `json:"username"`
	Verified             bool                            `json:"verified"`
	Visuals              interface{}                     `json:"visuals"`
	Badges               HydratableBadges                `json:"badges"`
	StationUrn           string                          `json:"station_urn"`
	StationPermalink     string                          `json:"station_permalink"`
	URL                  string                          `json:"url"`
}

type HydratableBadges struct {
	Pro          bool `json:"pro"`
	ProUnlimited bool `json:"pro_unlimited"`
	Verified     bool `json:"verified"`
}

type HydratableCreatorSubscription struct {
	Product HydratableProduct `json:"product"`
}

type HydratableSoundCloudPlaylist struct {
	ArtworkURL     interface{}                      `json:"artwork_url"`
	CreatedAt      string                           `json:"created_at"`
	Description    string                           `json:"description"`
	Duration       int64                            `json:"duration"`
	EmbeddableBy   string                           `json:"embeddable_by"`
	Genre          string                           `json:"genre"`
	ID             int64                            `json:"id"`
	Kind           string                           `json:"kind"`
	LabelName      string                           `json:"label_name"`
	LastModified   string                           `json:"last_modified"`
	License        string                           `json:"license"`
	LikesCount     int64                            `json:"likes_count"`
	ManagedByFeeds bool                             `json:"managed_by_feeds"`
	Permalink      string                           `json:"permalink"`
	PermalinkURL   string                           `json:"permalink_url"`
	Public         bool                             `json:"public"`
	PurchaseTitle  interface{}                      `json:"purchase_title"`
	PurchaseURL    interface{}                      `json:"purchase_url"`
	ReleaseDate    interface{}                      `json:"release_date"`
	RepostsCount   int64                            `json:"reposts_count"`
	SecretToken    string                           `json:"secret_token"`
	Sharing        string                           `json:"sharing"`
	TagList        string                           `json:"tag_list"`
	Title          string                           `json:"title"`
	URI            string                           `json:"uri"`
	UserID         int64                            `json:"user_id"`
	SetType        string                           `json:"set_type"`
	IsAlbum        bool                             `json:"is_album"`
	PublishedAt    string                           `json:"published_at"`
	DisplayDate    string                           `json:"display_date"`
	User           HydratableSoundCloudPlaylistUser `json:"user"`
	Tracks         []HydratableSoundCloudTrack      `json:"tracks"`
	TrackCount     int64                            `json:"track_count"`
	URL            string                           `json:"url"`
}

type HydratableSoundCloudTrack struct {
	ArtworkURL         string                      `json:"artwork_url,omitempty"`
	Caption            interface{}                 `json:"caption"`
	Commentable        bool                        `json:"commentable,omitempty"`
	CommentCount       int64                       `json:"comment_count,omitempty"`
	CreatedAt          string                      `json:"created_at,omitempty"`
	Description        string                      `json:"description"`
	Downloadable       bool                        `json:"downloadable,omitempty"`
	DownloadCount      int64                       `json:"download_count,omitempty"`
	Duration           int64                       `json:"duration,omitempty"`
	FullDuration       int64                       `json:"full_duration,omitempty"`
	EmbeddableBy       string                      `json:"embeddable_by,omitempty"`
	Genre              string                      `json:"genre,omitempty"`
	HasDownloadsLeft   bool                        `json:"has_downloads_left,omitempty"`
	ID                 int64                       `json:"id"`
	Kind               HydratableKind              `json:"kind"`
	LabelName          string                      `json:"label_name"`
	LastModified       string                      `json:"last_modified,omitempty"`
	License            string                      `json:"license,omitempty"`
	LikesCount         int64                       `json:"likes_count,omitempty"`
	Permalink          string                      `json:"permalink,omitempty"`
	PermalinkURL       string                      `json:"permalink_url,omitempty"`
	PlaybackCount      int64                       `json:"playback_count,omitempty"`
	Public             bool                        `json:"public,omitempty"`
	PublisherMetadata  HydratablePublisherMetadata `json:"publisher_metadata,omitempty"`
	PurchaseTitle      string                      `json:"purchase_title"`
	PurchaseURL        string                      `json:"purchase_url"`
	ReleaseDate        string                      `json:"release_date"`
	RepostsCount       int64                       `json:"reposts_count,omitempty"`
	SecretToken        interface{}                 `json:"secret_token"`
	Sharing            string                      `json:"sharing,omitempty"`
	State              string                      `json:"state,omitempty"`
	Streamable         bool                        `json:"streamable,omitempty"`
	TagList            string                      `json:"tag_list,omitempty"`
	Title              string                      `json:"title,omitempty"`
	TrackFormat        string                      `json:"track_format,omitempty"`
	URI                string                      `json:"uri,omitempty"`
	Urn                string                      `json:"urn,omitempty"`
	UserID             int64                       `json:"user_id,omitempty"`
	Visuals            interface{}                 `json:"visuals"`
	WaveformURL        string                      `json:"waveform_url,omitempty"`
	DisplayDate        string                      `json:"display_date,omitempty"`
	Media              HydratableMedia             `json:"media,omitempty"`
	StationUrn         string                      `json:"station_urn,omitempty"`
	StationPermalink   string                      `json:"station_permalink,omitempty"`
	TrackAuthorization string                      `json:"track_authorization,omitempty"`
	MonetizationModel  HydratableMonetizationModel `json:"monetization_model"`
	Policy             HydratablePolicy            `json:"policy"`
	User               HydratableTrackUser         `json:"user,omitempty"`
}

type HydratableMedia struct {
	Transcodings []HydratableTranscoding `json:"transcodings"`
}

type HydratableTranscoding struct {
	URL      string            `json:"url"`
	Preset   HydratablePreset  `json:"preset"`
	Duration int64             `json:"duration"`
	Snipped  bool              `json:"snipped"`
	Format   HydratableFormat  `json:"format"`
	Quality  HydratableQuality `json:"quality"`
}

type HydratableFormat struct {
	Protocol HydratableProtocol `json:"protocol"`
	MIMEType HydratableMIMEType `json:"mime_type"`
}

type HydratablePublisherMetadata struct {
	ID              int64   `json:"id"`
	Urn             string  `json:"urn"`
	Artist          string  `json:"artist"`
	ContainsMusic   bool    `json:"contains_music"`
	WriterComposer  *string `json:"writer_composer,omitempty"`
	ReleaseTitle    string  `json:"release_title"`
	AlbumTitle      *string `json:"album_title,omitempty"`
	UpcOrEan        *string `json:"upc_or_ean,omitempty"`
	Isrc            *string `json:"isrc,omitempty"`
	Explicit        *bool   `json:"explicit,omitempty"`
	PLine           *string `json:"p_line,omitempty"`
	PLineForDisplay *string `json:"p_line_for_display,omitempty"`
	CLine           *string `json:"c_line,omitempty"`
	CLineForDisplay *string `json:"c_line_for_display,omitempty"`
	Publisher       *string `json:"publisher,omitempty"`
}

type HydratableTrackUser struct {
	AvatarURL        string           `json:"avatar_url"`
	FirstName        string           `json:"first_name"`
	FollowersCount   int64            `json:"followers_count"`
	FullName         string           `json:"full_name"`
	ID               int64            `json:"id"`
	Kind             string           `json:"kind"`
	LastModified     string           `json:"last_modified"`
	LastName         string           `json:"last_name"`
	Permalink        string           `json:"permalink"`
	PermalinkURL     string           `json:"permalink_url"`
	URI              string           `json:"uri"`
	Urn              string           `json:"urn"`
	Username         string           `json:"username"`
	Verified         bool             `json:"verified"`
	City             string           `json:"city"`
	CountryCode      *string          `json:"country_code"`
	Badges           HydratableBadges `json:"badges"`
	StationUrn       string           `json:"station_urn"`
	StationPermalink string           `json:"station_permalink"`
}

type HydratableSoundCloudPlaylistUser struct {
	AvatarURL            string                          `json:"avatar_url"`
	City                 string                          `json:"city"`
	CommentsCount        int64                           `json:"comments_count"`
	CountryCode          string                          `json:"country_code"`
	CreatedAt            string                          `json:"created_at"`
	CreatorSubscriptions []HydratableCreatorSubscription `json:"creator_subscriptions"`
	CreatorSubscription  HydratableCreatorSubscription   `json:"creator_subscription"`
	Description          string                          `json:"description"`
	FollowersCount       int64                           `json:"followers_count"`
	FollowingsCount      int64                           `json:"followings_count"`
	FirstName            string                          `json:"first_name"`
	FullName             string                          `json:"full_name"`
	GroupsCount          int64                           `json:"groups_count"`
	ID                   int64                           `json:"id"`
	Kind                 string                          `json:"kind"`
	LastModified         string                          `json:"last_modified"`
	LastName             string                          `json:"last_name"`
	LikesCount           int64                           `json:"likes_count"`
	PlaylistLikesCount   int64                           `json:"playlist_likes_count"`
	Permalink            string                          `json:"permalink"`
	PermalinkURL         string                          `json:"permalink_url"`
	PlaylistCount        int64                           `json:"playlist_count"`
	RepostsCount         interface{}                     `json:"reposts_count"`
	TrackCount           int64                           `json:"track_count"`
	URI                  string                          `json:"uri"`
	Urn                  string                          `json:"urn"`
	Username             string                          `json:"username"`
	Verified             bool                            `json:"verified"`
	Visuals              interface{}                     `json:"visuals"`
	Badges               HydratableBadges                `json:"badges"`
	StationUrn           string                          `json:"station_urn"`
	StationPermalink     string                          `json:"station_permalink"`
}

type HydratableProduct struct {
	ID string `json:"id"`
}

type HydratableKind string

const (
	Track HydratableKind = "track"
)

type HydratableMIMEType string

const (
	AudioMPEG          HydratableMIMEType = "audio/mpeg"
	AudioOggCodecsOpus HydratableMIMEType = "audio/ogg; codecs=\"opus\""
)

type HydratableProtocol string

const (
	HLS         HydratableProtocol = "hls"
	Progressive HydratableProtocol = "progressive"
)

type HydratablePreset string

const (
	Mp31_0  HydratablePreset = "mp3_1_0"
	Opus0_0 HydratablePreset = "opus_0_0"
)

type HydratableQuality string

const (
	Sq HydratableQuality = "sq"
)

type HydratableMonetizationModel string

const (
	AdSupported   HydratableMonetizationModel = "AD_SUPPORTED"
	Blackbox      HydratableMonetizationModel = "BLACKBOX"
	NotApplicable HydratableMonetizationModel = "NOT_APPLICABLE"
	SubHighTier   HydratableMonetizationModel = "SUB_HIGH_TIER"
)

type HydratablePolicy string

const (
	Allow    HydratablePolicy = "ALLOW"
	Monetize HydratablePolicy = "MONETIZE"
	Snip     HydratablePolicy = "SNIP"
)
