package data

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type SerenDB struct {
	*sql.DB
	*Queries
}

func Connect() (*SerenDB, error) {
	db, err := sql.Open("sqlite3", "./seren.db")

	if err != nil {
		return nil, err
	}

	queries := New(db)

	return &SerenDB{
		DB:      db,
		Queries: queries,
	}, nil
}
