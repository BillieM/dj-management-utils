package database

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SoundCloudPlaylist struct {
	gorm.Model
	ExternalID int64 `gorm:"uniqueIndex"`
	Name       string
	SearchUrl  string
	Permalink  string

	Tracks []SoundCloudTrack `gorm:"many2many:playlist_tracks;"`

	NumTracks int `gorm:"-"` // not stored in db, calculated on select
}

type SoundCloudTrack struct {
	gorm.Model
	ExternalID          int64  `gorm:"uniqueIndex"`
	Name                string // change to title
	Permalink           string
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

// func (p *SoundCloudPlaylist) BeforeCreate(s *gorm.DB) (err error) {
// 	s.Statement.AddClause(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "external_id"}},
// 		DoNothing: true,
// 	})
// 	return nil
// }

// func (p *SoundCloudPlaylist) AfterSelect(s *gorm.DB) (err error) {
// 	var count int64
// 	s.Model(p).Association("Tracks").Count(&count)
// 	p.NumTracks = int(count)
// 	return nil
// }

// func (t *SoundCloudTrack) BeforeCreate(s *gorm.DB) (err error) {
// 	s.Statement.AddClause(clause.OnConflict{
// 		Columns:   []clause.Column{{Name: "external_id"}},
// 		DoUpdates:
// 	})
// 	return nil
// }

func (s *SerenDB) CreateSoundCloudPlaylist(p SoundCloudPlaylist) error {

	result := s.Clauses(
		clause.OnConflict{
			Columns:   []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{"updated_at"}),
		},
	).Create(&p)

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

/*
GetSoundCloudTracksByPlaylistID returns an array of SoundCloudTrack structs for a given playlist id
*/
func (s *SerenDB) GetSoundCloudTracksByPlaylistID(playlistId int64) ([]*SoundCloudTrack, error) {

	var playlist SoundCloudPlaylist
	var tracks []*SoundCloudTrack

	err := s.Where("external_id = ?", playlistId).Preload("Tracks").Find(&playlist).Error
	for _, t := range playlist.Tracks {
		v := t
		tracks = append(tracks, &v)
	}

	if err != nil {
		return tracks, err
	}

	return tracks, nil
}

func (s *SerenDB) GetSoundCloudTracksWithLocalPaths() ([]SoundCloudTrack, error) {

	var tracks []SoundCloudTrack

	result := s.Where("local_path IS NOT NULL").Find(&tracks)

	if result.Error != nil {
		return tracks, result.Error
	}

	return tracks, nil
}

/*
SaveSoundCloudTracks saves an array of SoundCloudTrack structs to the database
*/
func (s *SerenDB) SaveSoundCloudTracks(tracks []SoundCloudTrack) error {
	for _, track := range tracks {
		fmt.Println(track.Playlists)
	}

	result := s.Clauses(
		clause.OnConflict{
			Columns: []clause.Column{{Name: "external_id"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"updated_at",
				"local_path",
				"local_path_broken",
			}),
		},
	).Save(&tracks)

	if result.Error != nil {
		return result.Error
	}

	return nil
}
