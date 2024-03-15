// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0
// source: soundcloud.sql

package data

import (
	"context"
	"database/sql"
)

const countSoundCloudTracksByExternalID = `-- name: CountSoundCloudTracksByExternalID :one
SELECT COUNT(*)
FROM soundcloud_tracks t
JOIN soundcloud_playlist_tracks pt 
    ON t.id = pt.soundcloud_track_id
WHERE pt.soundcloud_playlist_id = ?1
`

func (q *Queries) CountSoundCloudTracksByExternalID(ctx context.Context, playlistExternalID sql.NullInt64) (int64, error) {
	row := q.db.QueryRowContext(ctx, countSoundCloudTracksByExternalID, playlistExternalID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countSoundCloudTracksByPlaylistID = `-- name: CountSoundCloudTracksByPlaylistID :one
SELECT COUNT(*)
FROM soundcloud_tracks t
JOIN soundcloud_playlist_tracks pt 
    ON t.id = pt.soundcloud_track_id
WHERE pt.soundcloud_playlist_id = ?1
`

func (q *Queries) CountSoundCloudTracksByPlaylistID(ctx context.Context, playlistID sql.NullInt64) (int64, error) {
	row := q.db.QueryRowContext(ctx, countSoundCloudTracksByPlaylistID, playlistID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getNumSoundCloudPlaylistByExternalID = `-- name: GetNumSoundCloudPlaylistByExternalID :one
SELECT count(*)
FROM soundcloud_playlists
WHERE external_id = ?1
`

func (q *Queries) GetNumSoundCloudPlaylistByExternalID(ctx context.Context, externalID sql.NullInt64) (int64, error) {
	row := q.db.QueryRowContext(ctx, getNumSoundCloudPlaylistByExternalID, externalID)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getNumSoundCloudPlaylistByURL = `-- name: GetNumSoundCloudPlaylistByURL :one
SELECT count(*)
FROM soundcloud_playlists
WHERE search_url = ?1
`

func (q *Queries) GetNumSoundCloudPlaylistByURL(ctx context.Context, url sql.NullString) (int64, error) {
	row := q.db.QueryRowContext(ctx, getNumSoundCloudPlaylistByURL, url)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const listSoundCloudPlaylists = `-- name: ListSoundCloudPlaylists :many
SELECT id, created_at, updated_at, external_id, name, search_url, permalink_url 
FROM soundcloud_playlists
`

func (q *Queries) ListSoundCloudPlaylists(ctx context.Context) ([]SoundcloudPlaylist, error) {
	rows, err := q.db.QueryContext(ctx, listSoundCloudPlaylists)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SoundcloudPlaylist
	for rows.Next() {
		var i SoundcloudPlaylist
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExternalID,
			&i.Name,
			&i.SearchUrl,
			&i.PermalinkUrl,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSoundCloudTracksByPlaylistExternalID = `-- name: ListSoundCloudTracksByPlaylistExternalID :many
SELECT t.id, t.created_at, t.updated_at, t.external_id, t.name, t.permalink_url, t.purchase_title, t.purchase_url, t.has_downloads_left, t.genre, t.artwork_url, t.tag_list, t.publisher_artist, t.sound_cloud_user, t.local_path, t.local_path_broken, t.removed_from_playlist
FROM soundcloud_tracks t
JOIN soundcloud_playlist_tracks pt 
    ON t.id = pt.soundcloud_track_id
JOIN soundcloud_playlists p 
    ON pt.soundcloud_playlist_id = p.id
WHERE p.external_id = ?1
`

func (q *Queries) ListSoundCloudTracksByPlaylistExternalID(ctx context.Context, playlistExternalID sql.NullInt64) ([]SoundcloudTrack, error) {
	rows, err := q.db.QueryContext(ctx, listSoundCloudTracksByPlaylistExternalID, playlistExternalID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SoundcloudTrack
	for rows.Next() {
		var i SoundcloudTrack
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExternalID,
			&i.Name,
			&i.PermalinkUrl,
			&i.PurchaseTitle,
			&i.PurchaseUrl,
			&i.HasDownloadsLeft,
			&i.Genre,
			&i.ArtworkUrl,
			&i.TagList,
			&i.PublisherArtist,
			&i.SoundCloudUser,
			&i.LocalPath,
			&i.LocalPathBroken,
			&i.RemovedFromPlaylist,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSoundCloudTracksByPlaylistID = `-- name: ListSoundCloudTracksByPlaylistID :many
SELECT t.id, t.created_at, t.updated_at, t.external_id, t.name, t.permalink_url, t.purchase_title, t.purchase_url, t.has_downloads_left, t.genre, t.artwork_url, t.tag_list, t.publisher_artist, t.sound_cloud_user, t.local_path, t.local_path_broken, t.removed_from_playlist
FROM soundcloud_tracks t
JOIN soundcloud_playlist_tracks pt 
    ON t.id = pt.soundcloud_track_id
WHERE soundcloud_playlist_id = ?1
`

func (q *Queries) ListSoundCloudTracksByPlaylistID(ctx context.Context, playlistID sql.NullInt64) ([]SoundcloudTrack, error) {
	rows, err := q.db.QueryContext(ctx, listSoundCloudTracksByPlaylistID, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SoundcloudTrack
	for rows.Next() {
		var i SoundcloudTrack
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExternalID,
			&i.Name,
			&i.PermalinkUrl,
			&i.PurchaseTitle,
			&i.PurchaseUrl,
			&i.HasDownloadsLeft,
			&i.Genre,
			&i.ArtworkUrl,
			&i.TagList,
			&i.PublisherArtist,
			&i.SoundCloudUser,
			&i.LocalPath,
			&i.LocalPathBroken,
			&i.RemovedFromPlaylist,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listSoundCloudTracksHasLocalPath = `-- name: ListSoundCloudTracksHasLocalPath :many
SELECT t.id, t.created_at, t.updated_at, t.external_id, t.name, t.permalink_url, t.purchase_title, t.purchase_url, t.has_downloads_left, t.genre, t.artwork_url, t.tag_list, t.publisher_artist, t.sound_cloud_user, t.local_path, t.local_path_broken, t.removed_from_playlist
FROM soundcloud_tracks t
WHERE local_path IS NOT NULL
`

func (q *Queries) ListSoundCloudTracksHasLocalPath(ctx context.Context) ([]SoundcloudTrack, error) {
	rows, err := q.db.QueryContext(ctx, listSoundCloudTracksHasLocalPath)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SoundcloudTrack
	for rows.Next() {
		var i SoundcloudTrack
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.ExternalID,
			&i.Name,
			&i.PermalinkUrl,
			&i.PurchaseTitle,
			&i.PurchaseUrl,
			&i.HasDownloadsLeft,
			&i.Genre,
			&i.ArtworkUrl,
			&i.TagList,
			&i.PublisherArtist,
			&i.SoundCloudUser,
			&i.LocalPath,
			&i.LocalPathBroken,
			&i.RemovedFromPlaylist,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const upsertSoundCloudPlaylist = `-- name: UpsertSoundCloudPlaylist :one
INSERT INTO soundcloud_playlists (
    created_at,
    updated_at,
    external_id,
    name,
    search_url,
    permalink_url
) VALUES (
    CURRENT_TIMESTAMP,
    CURRENT_TIMESTAMP,
    ?1,
    ?2,
    ?3,
    ?4
) ON CONFLICT (external_id) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP,

    name = coalesce(?2, name),
    search_url = coalesce(?3, search_url),
    permalink_url = coalesce(?4, permalink_url)

RETURNING id, created_at, updated_at, external_id, name, search_url, permalink_url
`

type UpsertSoundCloudPlaylistParams struct {
	ExternalID   sql.NullInt64
	Name         sql.NullString
	SearchUrl    sql.NullString
	PermalinkUrl sql.NullString
}

func (q *Queries) UpsertSoundCloudPlaylist(ctx context.Context, arg UpsertSoundCloudPlaylistParams) (SoundcloudPlaylist, error) {
	row := q.db.QueryRowContext(ctx, upsertSoundCloudPlaylist,
		arg.ExternalID,
		arg.Name,
		arg.SearchUrl,
		arg.PermalinkUrl,
	)
	var i SoundcloudPlaylist
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExternalID,
		&i.Name,
		&i.SearchUrl,
		&i.PermalinkUrl,
	)
	return i, err
}

const upsertSoundCloudPlaylistTrack = `-- name: UpsertSoundCloudPlaylistTrack :one
INSERT INTO soundcloud_playlist_tracks (
    soundcloud_playlist_id,
    soundcloud_track_id
) VALUES (
    ?1,
    ?2
) ON CONFLICT (soundcloud_playlist_id, soundcloud_track_id) DO NOTHING
RETURNING soundcloud_track_id, soundcloud_playlist_id
`

type UpsertSoundCloudPlaylistTrackParams struct {
	SoundcloudPlaylistID sql.NullInt64
	SoundcloudTrackID    sql.NullInt64
}

func (q *Queries) UpsertSoundCloudPlaylistTrack(ctx context.Context, arg UpsertSoundCloudPlaylistTrackParams) (SoundcloudPlaylistTrack, error) {
	row := q.db.QueryRowContext(ctx, upsertSoundCloudPlaylistTrack, arg.SoundcloudPlaylistID, arg.SoundcloudTrackID)
	var i SoundcloudPlaylistTrack
	err := row.Scan(&i.SoundcloudTrackID, &i.SoundcloudPlaylistID)
	return i, err
}

const upsertSoundCloudTrack = `-- name: UpsertSoundCloudTrack :one
INSERT INTO soundcloud_tracks (
    created_at,
    updated_at,
    external_id,
    name,
    permalink_url,
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
    ?1,
    ?2,
    ?3,
    ?4,
    ?5,
    ?6,
    ?7,
    ?8,
    ?9,
    ?10,
    ?11,
    ?12,
    ?13,
    ?14
) ON CONFLICT (external_id) DO UPDATE SET
    updated_at = CURRENT_TIMESTAMP,

    name = coalesce(?2, name),
    permalink_url = coalesce(?3, permalink_url),
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

RETURNING id, created_at, updated_at, external_id, name, permalink_url, purchase_title, purchase_url, has_downloads_left, genre, artwork_url, tag_list, publisher_artist, sound_cloud_user, local_path, local_path_broken, removed_from_playlist
`

type UpsertSoundCloudTrackParams struct {
	ExternalID          sql.NullInt64
	Name                sql.NullString
	PermalinkUrl        sql.NullString
	PurchaseTitle       sql.NullString
	PurchaseUrl         sql.NullString
	HasDownloadsLeft    sql.NullBool
	Genre               sql.NullString
	ArtworkUrl          sql.NullString
	TagList             sql.NullString
	PublisherArtist     sql.NullString
	SoundCloudUser      sql.NullString
	LocalPath           sql.NullString
	LocalPathBroken     sql.NullBool
	RemovedFromPlaylist sql.NullBool
}

func (q *Queries) UpsertSoundCloudTrack(ctx context.Context, arg UpsertSoundCloudTrackParams) (SoundcloudTrack, error) {
	row := q.db.QueryRowContext(ctx, upsertSoundCloudTrack,
		arg.ExternalID,
		arg.Name,
		arg.PermalinkUrl,
		arg.PurchaseTitle,
		arg.PurchaseUrl,
		arg.HasDownloadsLeft,
		arg.Genre,
		arg.ArtworkUrl,
		arg.TagList,
		arg.PublisherArtist,
		arg.SoundCloudUser,
		arg.LocalPath,
		arg.LocalPathBroken,
		arg.RemovedFromPlaylist,
	)
	var i SoundcloudTrack
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.ExternalID,
		&i.Name,
		&i.PermalinkUrl,
		&i.PurchaseTitle,
		&i.PurchaseUrl,
		&i.HasDownloadsLeft,
		&i.Genre,
		&i.ArtworkUrl,
		&i.TagList,
		&i.PublisherArtist,
		&i.SoundCloudUser,
		&i.LocalPath,
		&i.LocalPathBroken,
		&i.RemovedFromPlaylist,
	)
	return i, err
}
