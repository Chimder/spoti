package postgres

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type AlbumRepo struct {
	db *pgxpool.Pool
}

func NewAlbumRepo(db *pgxpool.Pool) *AlbumRepo {
	return &AlbumRepo{
		db: db,
	}
}

func (t *AlbumRepo) GetAlbum(ctx context.Context, albumID string) (json.RawMessage, error) {
	query := `
		WITH album AS (
    SELECT
        a.id,
        a.album_type,
        a.total_tracks,
        a.album_name,
        a.release_date,
        a.uri
    FROM albums a
    WHERE a.id = $1
),
album_artists AS (
    SELECT
        aa.album_id,
        jsonb_agg(
            jsonb_build_object(
                'id', ar.id,
                'name', ar.artist_name,
                'uri', ar.uri
            )
        ) AS artists
    FROM album_artists aa
    JOIN artists ar ON ar.id = aa.artist_id
    WHERE aa.album_id = $1
    GROUP BY aa.album_id
),
track_artists AS (
    SELECT
        t.id AS track_id,
        jsonb_agg(
            jsonb_build_object(
                'id', ar.id,
                'name', ar.artist_name,
                'uri', ar.uri
            )
        ) AS artists
    FROM tracks t
    JOIN artist_tracks at ON at.track_id = t.id
    JOIN artists ar ON ar.id = at.artist_id
    WHERE t.album_id = $1
    GROUP BY t.id
),
tracks AS (
    SELECT
        jsonb_agg(
            jsonb_build_object(
                'id', t.id,
                'name', t.track_name,
                'track_number', t.track_number,
                'disc_number', t.disc_number,
                'duration_ms', r.duration_ms,
                'explicit', t.explicit,
                'uri', t.uri,
                'artists', COALESCE(ta.artists, '[]')
            )
            ORDER BY t.disc_number, t.track_number
        ) AS items
    FROM tracks t
    JOIN recordings r ON r.id = t.recording_id
    LEFT JOIN track_artists ta ON ta.track_id = t.id
    WHERE t.album_id = $1
)
SELECT jsonb_build_object(
    'album_type', a.album_type,
    'total_tracks', a.total_tracks,
    'id', a.id,
    'name', a.album_name,
    'release_date', a.release_date,
    'uri', a.uri,
    'artists', COALESCE(aa.artists, '[]'),
    'tracks', jsonb_build_object(
        'items', COALESCE(t.items, '[]')
    )
)
FROM album a
LEFT JOIN album_artists aa ON aa.album_id = a.id
LEFT JOIN tracks t ON TRUE;
	`

	var data json.RawMessage
	err := t.db.QueryRow(ctx, query, albumID).Scan(&data)
	if err != nil {
		log.Error().Err(err).Msg("Get Album from db")
		return nil, err
	}

	return data, nil
}

func (al *AlbumRepo) GetAlbumsByIds(ctx context.Context, albumIDs []string) (json.RawMessage, error) {
	query := `WITH album AS (
    SELECT
        a.id,
        a.album_type,
        a.total_tracks,
        a.album_name,
        a.release_date,
        a.uri
    FROM albums a
    WHERE a.id = ANY($1::uuid[])
),

album_artists AS (
    SELECT
        aa.album_id,
        jsonb_agg(
            jsonb_build_object(
                'id', ar.id,
                'name', ar.artist_name,
                'uri', ar.uri
            )
            ORDER BY ar.artist_name
        ) AS artists
    FROM album_artists aa
    JOIN artists ar ON ar.id = aa.artist_id
    WHERE aa.album_id = ANY($1::uuid[])
    GROUP BY aa.album_id
),

track_artists AS (
    SELECT
        t.album_id,
        t.id AS track_id,
        jsonb_agg(
            jsonb_build_object(
                'id', ar.id,
                'name', ar.artist_name,
                'uri', ar.uri
            )
            ORDER BY ar.artist_name
        ) AS artists
    FROM tracks t
    JOIN artist_tracks at ON at.track_id = t.id
    JOIN artists ar ON ar.id = at.artist_id
    WHERE t.album_id = ANY($1::uuid[])
    GROUP BY t.album_id, t.id
),

tracks AS (
    SELECT
        t.album_id,
        jsonb_agg(
            jsonb_build_object(
                'id', t.id,
                'name', t.track_name,
                'track_number', t.track_number,
                'disc_number', t.disc_number,
                'duration_ms', r.duration_ms,
                'explicit', t.explicit,
                'uri', t.uri,
                'artists', COALESCE(ta.artists, '[]'::jsonb)
            )
            ORDER BY t.disc_number, t.track_number
        ) AS items
    FROM tracks t
    JOIN recordings r ON r.id = t.recording_id
    LEFT JOIN track_artists ta ON ta.track_id = t.id
    WHERE t.album_id = ANY($1::uuid[])
    GROUP BY t.album_id
)

SELECT jsonb_agg(
    jsonb_build_object(
        'album_type', a.album_type,
        'total_tracks', a.total_tracks,
        'id', a.id,
        'name', a.album_name,
        'release_date', a.release_date,
        'uri', a.uri,
        'artists', COALESCE(aa.artists, '[]'::jsonb),
        'tracks', jsonb_build_object(
            'items', COALESCE(t.items, '[]'::jsonb)
        )
    )
) AS albums
FROM album a
LEFT JOIN album_artists aa ON aa.album_id = a.id
LEFT JOIN tracks t ON t.album_id = a.id;
`
	var data json.RawMessage
	ids, err := al.parseUUIDs(albumIDs)
	if err != nil {
		log.Error().Err(err).Msg("err parse to uuid")
		return nil, err
	}

	err = al.db.QueryRow(ctx, query, ids).Scan(&data)
	if err != nil {
		log.Error().Err(err).Msg("Get Album from db")
		return nil, err
	}

	return data, nil
}

func (al *AlbumRepo) parseUUIDs(ids []string) ([]uuid.UUID, error) {
	res := make([]uuid.UUID, 0, len(ids))

	for _, id := range ids {
		uid, err := uuid.Parse(id)
		if err != nil {
			return nil, fmt.Errorf("invalid uuid %q: %w", id, err)
		}
		res = append(res, uid)
	}

	return res, nil
}

func (al *AlbumRepo) GetAlbumsTracks(ctx context.Context, albumID string) (json.RawMessage, error) {
	query := `SELECT jsonb_build_object(
    'album', to_jsonb(a) || jsonb_build_object(
    'tracks', (
            SELECT jsonb_agg(
                to_jsonb(t) || jsonb_build_object(
                'recording', to_jsonb(r),
                'artists', (
                        SELECT jsonb_agg(
                            jsonb_build_object(
                                'id', ar.id,
                                'name', ar.artist_name,
                                'uri', ar.uri
                            )
                            ORDER BY ar.artist_name
                        )
                        FROM artist_tracks at
                        JOIN artists ar ON ar.id = at.artist_id
                        WHERE at.track_id = t.id
                    )
                )
                ORDER BY t.disc_number, t.track_number
            )
            FROM tracks t
            JOIN recordings r ON r.id = t.recording_id
            WHERE t.album_id = a.id
        )
    )
) AS album
FROM albums a
WHERE a.id = $1;`

	var data json.RawMessage
	err := al.db.QueryRow(ctx, query, albumID).Scan(&data)
	if err != nil {
		log.Error().Err(err).Msg("Get Album from db")
		return nil, err
	}

	return data, nil
}
