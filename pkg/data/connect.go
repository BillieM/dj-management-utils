package data

import (
	"database/sql"

	"github.com/charmbracelet/log"

	"github.com/billiem/seren-management/pkg/helpers"
	_ "github.com/mattn/go-sqlite3"
	sqldblogger "github.com/simukti/sqldb-logger"
)

type SerenDB struct {
	*sql.DB
	*Queries
}

func Connect(c helpers.Config, l log.Logger) (*SerenDB, error) {

	dsn := "file:seren.db?cache=shared&mode=rwc"

	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		return nil, err
	}

	db = sqldblogger.OpenDriver(
		dsn,
		db.Driver(),
		&helpers.CharmLogAdapter{
			Logger: l,
		},
	)

	queries := New(db)

	return &SerenDB{
		DB:      db,
		Queries: queries,
	}, nil
}
