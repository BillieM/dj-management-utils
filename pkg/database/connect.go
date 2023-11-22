package database

import "gorm.io/gorm"

type SerenDB struct {
	*gorm.DB
}

func Connect() (*SerenDB, error) {
	return &SerenDB{}, nil
}
