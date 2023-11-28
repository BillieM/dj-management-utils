package database

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SoundCloudPlaylist struct {
	gorm.Model
	ExternalID int64 `gorm:"uniqueIndex"`
	Name       string
	SearchUrl  string
	Permalink  string
	Tracks     []SoundCloudTrack `gorm:"many2many:playlist_tracks;"`
}

type SoundCloudTrack struct {
	gorm.Model
	ExternalID       int64  `gorm:"uniqueIndex"`
	Name             string // change to title
	Permalink        string
	PurchaseTitle    string
	PurchaseURL      string
	HasDownloadsLeft bool
	Genre            string
	ArtworkURL       string
	TagList          string
	PublisherArtist  string // users/ artists use relationships
	SoundCloudUser   string
	Playlists        []SoundCloudPlaylist `gorm:"many2many:playlist_tracks;"`
}

func (p *SoundCloudPlaylist) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoNothing: true,
	})
	return nil
}

func (t *SoundCloudTrack) BeforeCreate(tx *gorm.DB) (err error) {
	tx.Statement.AddClause(clause.OnConflict{
		Columns:   []clause.Column{{Name: "external_id"}},
		DoNothing: true,
	})
	return nil
}

func (s *SerenDB) CreateSoundCloudPlaylist(p SoundCloudPlaylist) error {

	result := s.Create(&p)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (s *SerenDB) GetSoundCloudPlaylists() ([]SoundCloudPlaylist, error) {

	var playlists []SoundCloudPlaylist

	s.Find(&playlists)

	return playlists, nil
}

func (s *SerenDB) GetSoundCloudPlaylistByURL(searchUrl string) (SoundCloudPlaylist, error) {

	var playlist SoundCloudPlaylist
	result := s.Where("search_url = ?", searchUrl).Limit(1).Find(&playlist)

	if result.RowsAffected == 0 { // nil
		return playlist, result.Error
	}

	return playlist, result.Error

}

func (s *SerenDB) GetSoundCloudPlaylistByExternalID(externalId int64) (SoundCloudPlaylist, error) {

	var playlist SoundCloudPlaylist
	result := s.Where("external_id = ?", externalId).Limit(1).Find(&playlist)

	if result.RowsAffected == 0 { // nil
		return playlist, result.Error
	}

	return playlist, result.Error

}

func (s *SerenDB) GetSoundCloudTracks(playlistId int64) ([]SoundCloudTrack, error) {

	var playlist SoundCloudPlaylist
	var tracks []SoundCloudTrack

	err := s.Where("external_id = ?", playlistId).Preload("Tracks").Find(&playlist).Error
	tracks = playlist.Tracks

	if err != nil {
		return tracks, err
	}

	return tracks, nil
}