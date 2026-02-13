package postgres

import (
	"context"
	"encoding/json"
	"spoti/internal/domain/playlist"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type PlaylistRepo struct {
	db *pgxpool.Pool
}

func NewPlaylistRepo(db *pgxpool.Pool) *PlaylistRepo {
	return &PlaylistRepo{
		db: db,
	}
}

func (pl *PlaylistRepo) CreatePlaylist(ctx context.Context, p playlist.CreatePlaylistReq) (uuid.UUID, error) {
	query := `
INSERT INTO playlists (owner_id, playlist_name, description, image, is_public)
VALUES ($1, $2, $3, $4, $5)
RETURNING id
`
	var id uuid.UUID
	err := pl.db.QueryRow(ctx, query, p.OwnerId, p.PlaylistName, p.Description, p.Image, p.IsPublic).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("err create playlist")
		return uuid.UUID{}, err
	}

	return id, err
}

func (pl *PlaylistRepo) GetPlaylistById(ctx context.Context, playlistId string, limit, offset int) (playlist.PlaylistJson, error) {
	if limit <= 0 || limit > 20 {
		limit = 20
	}
	if offset < 0 {
		offset = 0
	}
	query := `
	WITH playlist_cte AS (
		SELECT *
		FROM playlists pl
		WHERE pl.id = $1
	),
	artist_tracks_cte AS (
		SELECT
			at.track_id,
			jsonb_agg(
				jsonb_build_object(
					'id', ar.id,
					'name', ar.artist_name,
					'uri', ar.uri
				)
			) as artists
		FROM artist_tracks at
		JOIN artists ar ON ar.id = at.artist_id
		GROUP BY at.track_id
	),
	playlist_items_agg AS (
		SELECT
			jsonb_agg(
				jsonb_build_object(
					'added_at', pt.added_at,
					'is_local', tr.islocal,
					'track', to_jsonb(tr.*) || jsonb_build_object(
						'artists', COALESCE(atc.artists, '[]'::jsonb)
					)
				) ORDER BY pt.track_position
			) as items
		FROM playlist_tracks pt
		JOIN tracks tr ON tr.id = pt.track_id
		LEFT JOIN artist_tracks_cte atc ON atc.track_id = tr.id
		WHERE pt.playlist_id = $1
		LIMIT $2 OFFSET $3
	)
	SELECT
		to_jsonb(pi.*) ||
		jsonb_build_object(
			'tracks', jsonb_build_object(
				'total', pi.total,
				'limit', $2,
				'offset', $3,
				'items', COALESCE((SELECT items FROM playlist_items_agg), '[]'::jsonb)
			)
		) as playlist
	FROM playlist_cte pi;
	`

	var data json.RawMessage
	err := pl.db.QueryRow(ctx, query, playlistId, limit, offset).Scan(&data)
	if err != nil {
		return playlist.PlaylistJson{}, err
	}

	var playlistData playlist.PlaylistJson
	if err := json.Unmarshal(data, &playlistData); err != nil {
		return playlist.PlaylistJson{}, err
	}

	return playlistData, nil
}

func (pl *PlaylistRepo) AddToPlaylist(ctx context.Context, playlist_id, track_id string) error {
	query := `
INSERT INTO playlist_tracks (playlist_id, track_id, track_position)
SELECT $1, $2, COALESCE(MAX(track_position), 0) + 1
FROM playlist_tracks
WHERE playlist_id = $1;
`
	_, err := pl.db.Exec(ctx, query, playlist_id, track_id)
	if err != nil {
		log.Error().Err(err).Msg("err add to playlist")
		return err
	}

	return err
}

type UpdatePlaylistReq struct {
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Public      *bool   `json:"public"`
}

func (pl *PlaylistRepo) UpdatePlaylist(ctx context.Context, playlistId string, req UpdatePlaylistReq) error {
	query := `
        UPDATE playlists
        SET
            playlist_name = COALESCE($2, playlist_name),
            description = COALESCE($3, description),
            is_public = COALESCE($4, is_public)
        WHERE id = $1
    `

	_, err := pl.db.Exec(ctx, query, playlistId, req.Name, req.Description, req.Public)
	if err != nil {
		log.Error().Err(err).Msg("err updatePlaylist db")
		return err
	}

	return err
}
func (pl *PlaylistRepo) DeleteFromPlaylist(ctx context.Context, playlistId, trackId string) error {
	query := `
WITH removed AS (
    DELETE FROM playlist_tracks
    WHERE playlist_id = $1
      AND track_id = $2
    RETURNING track_position
)
UPDATE playlist_tracks
SET track_position = track_position - 1
WHERE playlist_id = $1
  AND track_position > (SELECT track_position FROM removed);
`
	_, err := pl.db.Exec(ctx, query, playlistId, trackId)
	if err != nil {
		log.Error().Err(err).Msg("err delete from playlist")
		return err
	}

	return err
}

func (pl *PlaylistRepo) GetAllUserPlaylists(ctx context.Context, userId string) ([]playlist.Playlist, error) {
	query := `SELECT * FROM playlists WHERE owner_id = $1;`

	rows, err := pl.db.Query(ctx, query, userId)
	if err != nil {
		log.Error().Err(err).Msg("err get all playlists by id db")
		return nil, err
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[PlaylistRow])
	if err != nil {
		log.Error().Err(err).Msg("err collect all playlists by id")
		return nil, err
	}

	return PlayListsToDomain(data), nil
}
