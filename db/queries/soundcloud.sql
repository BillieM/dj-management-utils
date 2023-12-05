-- name: ListSoundCloudPlaylists :many
SELECT * 
FROM soundcloud_playlists;

-- name: GetNumSoundCloudPlaylistByURL :one
SELECT count(*)
FROM soundcloud_playlists
WHERE search_url = @url; 

-- name: GetNumSoundCloudPlaylistByExternalID :one
SELECT count(*)
FROM soundcloud_playlists
WHERE external_id = @external_id; 

-- name: ListSoundCloudTracksByPlaylistID :many
SELECT t.*
FROM soundcloud_tracks t
JOIN soundcloud_playlist_tracks pt 
    ON t.id = pt.soundcloud_track_id
WHERE soundcloud_playlist_id = @playlist_id;

-- name: ListSoundCloudTracksByPlaylistExternalID :many
SELECT t.*
FROM soundcloud_tracks t
JOIN soundcloud_playlist_tracks pt 
    ON t.id = pt.soundcloud_track_id
JOIN soundcloud_playlists p 
    ON pt.soundcloud_playlist_id = p.id
WHERE p.external_id = @playlist_external_id;

-- name: ListSoundCloudTracksHasLocalPath :many
SELECT t.*
FROM soundcloud_tracks t
WHERE local_path IS NOT NULL;

-- name: CountSoundCloudTracksByPlaylistID :one
SELECT COUNT(*)
FROM soundcloud_tracks t
JOIN soundcloud_playlist_tracks pt 
    ON t.id = pt.soundcloud_track_id
WHERE pt.soundcloud_playlist_id = @playlist_id;

-- name: CountSoundCloudTracksByExternalID :one
SELECT COUNT(*)
FROM soundcloud_tracks t
JOIN soundcloud_playlist_tracks pt 
    ON t.id = pt.soundcloud_track_id
WHERE pt.soundcloud_playlist_id = @playlist_external_id;

-- name: UpsertSoundCloudPlaylist :one
INSERT INTO soundcloud_playlists (
    created_at,
    updated_at,
    external_id,
    name,
    search_url,
    permalink
) VALUES (
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    sqlc.narg('external_id'),
    sqlc.narg('name'),
    sqlc.narg('search_url'),
    sqlc.narg('permalink')
) ON CONFLICT (external_id) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP

    name = coalesce(?2, name),
    search_url = coalesce(?3, search_url),
    permalink = coalesce(?4, permalink)

RETURNING *;

-- name: UpsertSoundCloudTrack :one
INSERT INTO soundcloud_tracks (
    created_at,
    updated_at,
    external_id,
    name,
    permalink,
    purchase_title,
    purchase_url,
    has_downloads_left,
    genre,
    artwork_url,
    tag_list,
    publisher_artist,
    sound_cloud_user,
    local_path,
    local_path_broken,
    removed_from_playlist
) VALUES (
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    sqlc.narg('external_id'),
    sqlc.narg('name'),
    sqlc.narg('permalink'),
    sqlc.narg('purchase_title'),
    sqlc.narg('purchase_url'),
    sqlc.narg('has_downloads_left'),
    sqlc.narg('genre'),
    sqlc.narg('artwork_url'),
    sqlc.narg('tag_list'),
    sqlc.narg('publisher_artist'),
    sqlc.narg('sound_cloud_user'),
    sqlc.narg('local_path'),
    sqlc.narg('local_path_broken'),
    sqlc.narg('removed_from_playlist')
) ON CONFLICT (external_id) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP,

    name = coalesce(?2, name),
    permalink = coalesce(?3, permalink),
    purchase_title = coalesce(?4, purchase_title),
    purchase_url = coalesce(?5, purchase_url),
    has_downloads_left = coalesce(?6, has_downloads_left),
    genre = coalesce(?7, genre),
    artwork_url = coalesce(?8, artwork_url),
    tag_list = coalesce(?9, tag_list),
    publisher_artist = coalesce(?10, publisher_artist),
    sound_cloud_user = coalesce(?11, sound_cloud_user),
    local_path = coalesce(?12, local_path),
    local_path_broken = coalesce(?13, local_path_broken),
    removed_from_playlist = coalesce(?14, removed_from_playlist)

RETURNING *;

-- name: UpsertSoundCloudPlaylistTrack :one
INSERT INTO soundcloud_playlist_tracks (
    soundcloud_playlist_id,
    soundcloud_track_id
) VALUES (
    sqlc.narg('soundcloud_playlist_id'),
    sqlc.narg('soundcloud_track_id')
) ON CONFLICT (soundcloud_playlist_id, soundcloud_track_id) DO NOTHING
RETURNING *;