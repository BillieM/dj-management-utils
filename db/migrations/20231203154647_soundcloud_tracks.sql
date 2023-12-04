-- +goose Up
-- +goose StatementBegin
CREATE TABLE soundcloud_tracks (
    id                    INTEGER  PRIMARY KEY AUTOINCREMENT,
    created_at            DATETIME,
    updated_at            DATETIME,
    external_id           INTEGER,
    name                  TEXT,
    permalink             TEXT,
    purchase_title        TEXT,
    purchase_url          TEXT,
    has_downloads_left    BOOLEAN,
    genre                 TEXT,
    artwork_url           TEXT,
    tag_list              TEXT,
    publisher_artist      TEXT,
    sound_cloud_user      TEXT,
    local_path            TEXT,
    local_path_broken     BOOLEAN,
    removed_from_playlist BOOLEAN
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE soundcloud_tracks;
-- +goose StatementEnd
