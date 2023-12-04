-- +goose Up
-- +goose StatementBegin
CREATE TABLE soundcloud_playlist_tracks (
    soundcloud_track_id    INTEGER,
    soundcloud_playlist_id INTEGER,
    PRIMARY KEY (
        soundcloud_track_id,
        soundcloud_playlist_id
    ),
    CONSTRAINT fk_playlist_tracks_soundcloud_track FOREIGN KEY (
        soundcloud_track_id
    )
    REFERENCES soundcloud_tracks (id),
    CONSTRAINT fk_playlist_tracks_soundcloud_playlist FOREIGN KEY (
        soundcloud_playlist_id
    )
    REFERENCES soundcloud_playlists (id) 
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE soundcloud_playlist_tracks;
-- +goose StatementEnd
