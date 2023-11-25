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
	Tracks     []*SoundCloudTrack `gorm:"many2many:playlist_tracks;"`

	NumTracks int `gorm:"-:all"` // not stored in DB, returned when getting playlist without tracks
}

type SoundCloudTrack struct {
	gorm.Model
	ExternalID int64 `gorm:"uniqueIndex"`
	Name       string
	Permalink  string
	Playlists  []*SoundCloudPlaylist `gorm:"many2many:playlist_tracks;"`
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
	result := s.Where("external_id = ?", searchUrl).Limit(1).Find(&playlist)

	if result.RowsAffected == 0 { // nil
		return playlist, result.Error
	}

	return playlist, result.Error

}
