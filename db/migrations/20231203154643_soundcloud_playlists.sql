-- +goose Up
-- +goose StatementBegin
CREATE TABLE soundcloud_playlists (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    created_at  DATETIME,
    updated_at  DATETIME,
    external_id INTEGER UNIQUE,
    name        TEXT,
    search_url  TEXT,
    permalink   TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE soundcloud_playlists;
-- +goose StatementEnd
