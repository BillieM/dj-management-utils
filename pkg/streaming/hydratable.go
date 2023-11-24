package streaming

import (
	"encoding/json"
	"fmt"

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
			var user SoundCloudUser

			err := json.Unmarshal(jsonBytes, &user)
			if err != nil {
				return err
			}
			h.User = user
		case "playlist":
			var playlist SoundCloudPlaylist
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
	User     SoundCloudUser
	Playlist SoundCloudPlaylist
}

type SoundCloudUser struct {
	AvatarURL            string                `json:"avatar_url"`
	City                 string                `json:"city"`
	CommentsCount        int64                 `json:"comments_count"`
	CountryCode          string                `json:"country_code"`
	CreatedAt            string                `json:"created_at"`
	CreatorSubscriptions []CreatorSubscription `json:"creator_subscriptions"`
	CreatorSubscription  CreatorSubscription   `json:"creator_subscription"`
	Description          string                `json:"description"`
	FollowersCount       int64                 `json:"followers_count"`
	FollowingsCount      int64                 `json:"followings_count"`
	FirstName            string                `json:"first_name"`
	FullName             string                `json:"full_name"`
	GroupsCount          int64                 `json:"groups_count"`
	ID                   int64                 `json:"id"`
	Kind                 string                `json:"kind"`
	LastModified         string                `json:"last_modified"`
	LastName             string                `json:"last_name"`
	LikesCount           int64                 `json:"likes_count"`
	PlaylistLikesCount   int64                 `json:"playlist_likes_count"`
	Permalink            string                `json:"permalink"`
	PermalinkURL         string                `json:"permalink_url"`
	PlaylistCount        int64                 `json:"playlist_count"`
	RepostsCount         int64                 `json:"reposts_count"`
	TrackCount           int64                 `json:"track_count"`
	URI                  string                `json:"uri"`
	Urn                  string                `json:"urn"`
	Username             string                `json:"username"`
	Verified             bool                  `json:"verified"`
	Visuals              interface{}           `json:"visuals"`
	Badges               Badges                `json:"badges"`
	StationUrn           string                `json:"station_urn"`
	StationPermalink     string                `json:"station_permalink"`
	URL                  string                `json:"url"`
}

type Badges struct {
	Pro          bool `json:"pro"`
	ProUnlimited bool `json:"pro_unlimited"`
	Verified     bool `json:"verified"`
}

type CreatorSubscription struct {
	Product Product `json:"product"`
}

type SoundCloudPlaylist struct {
	ArtworkURL     interface{}            `json:"artwork_url"`
	CreatedAt      string                 `json:"created_at"`
	Description    string                 `json:"description"`
	Duration       int64                  `json:"duration"`
	EmbeddableBy   string                 `json:"embeddable_by"`
	Genre          string                 `json:"genre"`
	ID             int64                  `json:"id"`
	Kind           string                 `json:"kind"`
	LabelName      string                 `json:"label_name"`
	LastModified   string                 `json:"last_modified"`
	License        string                 `json:"license"`
	LikesCount     int64                  `json:"likes_count"`
	ManagedByFeeds bool                   `json:"managed_by_feeds"`
	Permalink      string                 `json:"permalink"`
	PermalinkURL   string                 `json:"permalink_url"`
	Public         bool                   `json:"public"`
	PurchaseTitle  interface{}            `json:"purchase_title"`
	PurchaseURL    interface{}            `json:"purchase_url"`
	ReleaseDate    interface{}            `json:"release_date"`
	RepostsCount   int64                  `json:"reposts_count"`
	SecretToken    string                 `json:"secret_token"`
	Sharing        string                 `json:"sharing"`
	TagList        string                 `json:"tag_list"`
	Title          string                 `json:"title"`
	URI            string                 `json:"uri"`
	UserID         int64                  `json:"user_id"`
	SetType        string                 `json:"set_type"`
	IsAlbum        bool                   `json:"is_album"`
	PublishedAt    string                 `json:"published_at"`
	DisplayDate    string                 `json:"display_date"`
	User           SoundCloudPlaylistUser `json:"user"`
	Tracks         []TrackElement         `json:"tracks"`
	TrackCount     int64                  `json:"track_count"`
	URL            string                 `json:"url"`
}

type TrackElement struct {
	ArtworkURL         *string            `json:"artwork_url,omitempty"`
	Caption            interface{}        `json:"caption"`
	Commentable        *bool              `json:"commentable,omitempty"`
	CommentCount       *int64             `json:"comment_count,omitempty"`
	CreatedAt          *string            `json:"created_at,omitempty"`
	Description        *string            `json:"description"`
	Downloadable       *bool              `json:"downloadable,omitempty"`
	DownloadCount      *int64             `json:"download_count,omitempty"`
	Duration           *int64             `json:"duration,omitempty"`
	FullDuration       *int64             `json:"full_duration,omitempty"`
	EmbeddableBy       *string            `json:"embeddable_by,omitempty"`
	Genre              *string            `json:"genre,omitempty"`
	HasDownloadsLeft   *bool              `json:"has_downloads_left,omitempty"`
	ID                 int64              `json:"id"`
	Kind               Kind               `json:"kind"`
	LabelName          *string            `json:"label_name"`
	LastModified       *string            `json:"last_modified,omitempty"`
	License            *string            `json:"license,omitempty"`
	LikesCount         *int64             `json:"likes_count,omitempty"`
	Permalink          *string            `json:"permalink,omitempty"`
	PermalinkURL       *string            `json:"permalink_url,omitempty"`
	PlaybackCount      *int64             `json:"playback_count,omitempty"`
	Public             *bool              `json:"public,omitempty"`
	PublisherMetadata  *PublisherMetadata `json:"publisher_metadata,omitempty"`
	PurchaseTitle      *string            `json:"purchase_title"`
	PurchaseURL        *string            `json:"purchase_url"`
	ReleaseDate        *string            `json:"release_date"`
	RepostsCount       *int64             `json:"reposts_count,omitempty"`
	SecretToken        interface{}        `json:"secret_token"`
	Sharing            *string            `json:"sharing,omitempty"`
	State              *string            `json:"state,omitempty"`
	Streamable         *bool              `json:"streamable,omitempty"`
	TagList            *string            `json:"tag_list,omitempty"`
	Title              *string            `json:"title,omitempty"`
	TrackFormat        *string            `json:"track_format,omitempty"`
	URI                *string            `json:"uri,omitempty"`
	Urn                *string            `json:"urn,omitempty"`
	UserID             *int64             `json:"user_id,omitempty"`
	Visuals            interface{}        `json:"visuals"`
	WaveformURL        *string            `json:"waveform_url,omitempty"`
	DisplayDate        *string            `json:"display_date,omitempty"`
	Media              *Media             `json:"media,omitempty"`
	StationUrn         *string            `json:"station_urn,omitempty"`
	StationPermalink   *string            `json:"station_permalink,omitempty"`
	TrackAuthorization *string            `json:"track_authorization,omitempty"`
	MonetizationModel  MonetizationModel  `json:"monetization_model"`
	Policy             Policy             `json:"policy"`
	User               *TrackUser         `json:"user,omitempty"`
}

func (t TrackElement) String() string {
	return fmt.Sprintf("%v: %s - %s", t.ID, *t.Title, *t.PermalinkURL)
}

type Media struct {
	Transcodings []Transcoding `json:"transcodings"`
}

type Transcoding struct {
	URL      string  `json:"url"`
	Preset   Preset  `json:"preset"`
	Duration int64   `json:"duration"`
	Snipped  bool    `json:"snipped"`
	Format   Format  `json:"format"`
	Quality  Quality `json:"quality"`
}

type Format struct {
	Protocol Protocol `json:"protocol"`
	MIMEType MIMEType `json:"mime_type"`
}

type PublisherMetadata struct {
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

type TrackUser struct {
	AvatarURL        string  `json:"avatar_url"`
	FirstName        string  `json:"first_name"`
	FollowersCount   int64   `json:"followers_count"`
	FullName         string  `json:"full_name"`
	ID               int64   `json:"id"`
	Kind             string  `json:"kind"`
	LastModified     string  `json:"last_modified"`
	LastName         string  `json:"last_name"`
	Permalink        string  `json:"permalink"`
	PermalinkURL     string  `json:"permalink_url"`
	URI              string  `json:"uri"`
	Urn              string  `json:"urn"`
	Username         string  `json:"username"`
	Verified         bool    `json:"verified"`
	City             string  `json:"city"`
	CountryCode      *string `json:"country_code"`
	Badges           Badges  `json:"badges"`
	StationUrn       string  `json:"station_urn"`
	StationPermalink string  `json:"station_permalink"`
}

type SoundCloudPlaylistUser struct {
	AvatarURL            string                `json:"avatar_url"`
	City                 string                `json:"city"`
	CommentsCount        int64                 `json:"comments_count"`
	CountryCode          string                `json:"country_code"`
	CreatedAt            string                `json:"created_at"`
	CreatorSubscriptions []CreatorSubscription `json:"creator_subscriptions"`
	CreatorSubscription  CreatorSubscription   `json:"creator_subscription"`
	Description          string                `json:"description"`
	FollowersCount       int64                 `json:"followers_count"`
	FollowingsCount      int64                 `json:"followings_count"`
	FirstName            string                `json:"first_name"`
	FullName             string                `json:"full_name"`
	GroupsCount          int64                 `json:"groups_count"`
	ID                   int64                 `json:"id"`
	Kind                 string                `json:"kind"`
	LastModified         string                `json:"last_modified"`
	LastName             string                `json:"last_name"`
	LikesCount           int64                 `json:"likes_count"`
	PlaylistLikesCount   int64                 `json:"playlist_likes_count"`
	Permalink            string                `json:"permalink"`
	PermalinkURL         string                `json:"permalink_url"`
	PlaylistCount        int64                 `json:"playlist_count"`
	RepostsCount         interface{}           `json:"reposts_count"`
	TrackCount           int64                 `json:"track_count"`
	URI                  string                `json:"uri"`
	Urn                  string                `json:"urn"`
	Username             string                `json:"username"`
	Verified             bool                  `json:"verified"`
	Visuals              interface{}           `json:"visuals"`
	Badges               Badges                `json:"badges"`
	StationUrn           string                `json:"station_urn"`
	StationPermalink     string                `json:"station_permalink"`
}

type Product struct {
	ID string `json:"id"`
}

type Kind string

const (
	Track Kind = "track"
)

type MIMEType string

const (
	AudioMPEG          MIMEType = "audio/mpeg"
	AudioOggCodecsOpus MIMEType = "audio/ogg; codecs=\"opus\""
)

type Protocol string

const (
	HLS         Protocol = "hls"
	Progressive Protocol = "progressive"
)

type Preset string

const (
	Mp31_0  Preset = "mp3_1_0"
	Opus0_0 Preset = "opus_0_0"
)

type Quality string

const (
	Sq Quality = "sq"
)

type MonetizationModel string

const (
	AdSupported   MonetizationModel = "AD_SUPPORTED"
	Blackbox      MonetizationModel = "BLACKBOX"
	NotApplicable MonetizationModel = "NOT_APPLICABLE"
	SubHighTier   MonetizationModel = "SUB_HIGH_TIER"
)

type Policy string

const (
	Allow    Policy = "ALLOW"
	Monetize Policy = "MONETIZE"
	Snip     Policy = "SNIP"
)
