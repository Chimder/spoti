package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Playlist struct {
	ID           uuid.UUID `db:"id"`
	OwnerID      uuid.UUID `db:"owner_id"`
	PlaylistName string    `db:"playlist_name"`
	Description  string    `db:"description"`
	Image        string    `db:"image"`
	IsPublic     bool      `db:"is_public"`
	Total        int       `db:"total"`
	CreatedAt    time.Time `db:"created_at"`
}
type PlaylistRepo struct {
	db *pgxpool.Pool
}

func NewPlaylistRepo(db *pgxpool.Pool) *AlbumRepo {
	return &AlbumRepo{
		db: db,
	}
}

func (pl *PlaylistRepo) GetPlaylistById(ctx context.Context, playlistId string, limit, offset int) (json.RawMessage, error) {
	if limit <= 0 || limit > 100 {
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
			'collaborative', false,
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
		log.Error().Err(err).Msg("Get playlist from db")
		return nil, err
	}

	return data, nil
}

func (art *ArtistRepo) AddToPlaylist(ctx context.Context, playlist_id, track_id string) error {
	query := `
INSERT INTO playlist_tracks (playlist_id, track_id, track_position)
SELECT $1, $2, COALESCE(MAX(track_position), 0) + 1
FROM playlist_tracks
WHERE playlist_id = $1;
`
	_, err := art.db.Exec(ctx, query, playlist_id, track_id)
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
func (art *ArtistRepo) DeleteFromPlaylist(ctx context.Context, playlistId, trackId string) error {
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
	_, err := art.db.Exec(ctx, query, playlistId, trackId)
	if err != nil {
		log.Error().Err(err).Msg("err delete from playlist")
		return err
	}

	return err
}

func (art *ArtistRepo) GetAllUserPlaylists(ctx context.Context, userId string) ([]Playlist, error) {
	query := `SELECT * FROM playlists WHERE owner_id = $1;`

	rows, err := art.db.Query(ctx, query, userId)
	if err != nil {
		log.Error().Err(err).Msg("err get all playlists by id db")
		return nil, err
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Playlist])
	if err != nil {
		log.Error().Err(err).Msg("err collect all playlists by id")
		return nil, err
	}

	return data, nil
}

type CreatePlaylistReq struct {
	OwnerID      uuid.UUID
	PlaylistName string
	Description  string
	Image        string
	IsPublic     bool
}

func (art *ArtistRepo) CreatePlaylist(ctx context.Context, cp CreatePlaylistReq) error {
	query := `
	INSERT INTO playlists (owner_id, playlist_name, description, image, is_public)
	VALUES ($1, $2, $3, $4, $5)
	`

	_, err := art.db.Exec(ctx, query, cp.OwnerID, cp.PlaylistName, cp.Description, cp.Image, cp.IsPublic)
	if err != nil {
		log.Error().Err(err).Msg("err create playlist")
		return err
	}

	return err
}
