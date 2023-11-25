package database

import (
	"github.com/billiem/seren-management/pkg/helpers"
	"github.com/billiem/seren-management/pkg/projectpath"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type SerenDB struct {
	*gorm.DB
}

func Connect() (*SerenDB, error) {

	dbPath := helpers.JoinFilepathToSlash(projectpath.Root, "seren.db")

	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		panic("failed to connect to database")
	}

	// Migrate the schema
	db.AutoMigrate(&SoundCloudPlaylist{})
	db.AutoMigrate(&SoundCloudTrack{})

	return &SerenDB{
		DB: db,
	}, nil
}
