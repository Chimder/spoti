package postgres

import (
	"context"
	"encoding/json"
	"time"

	"github.com/google/uuid"
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
func (art *ArtistRepo) GetPlaylistById(ctx context.Context, playlistId string) (Playlist, error) {
	return Playlist{}, nil
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
