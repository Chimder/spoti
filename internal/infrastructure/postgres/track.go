package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type Track struct {
	ID          uuid.UUID `db:"id"`
	AlbumID     uuid.UUID `db:"album_id"`
	RecordingID uuid.UUID `db:"recording_id"`
	TrackName   string    `db:"track_name"`
	TrackNumber int16     `db:"track_number"`
	DiscNumber  int16     `db:"disc_number"`
	Explicit    bool      `db:"explicit"`
	IsPlayable  bool      `db:"is_playable"`
	TrackType   string    `db:"track_type"`
	URI         string    `db:"uri"`
	IsLocal     bool      `db:"islocal"`

	ISRC       string `db:"isrc"`
	DurationMs int64  `db:"duration_ms"`
	Popularity int    `db:"popularity"`
	PlayCount  int64  `db:"play_count"`
	AudioURI   string `db:"audio_uri"`
}

type TrackRepo struct {
	db *pgxpool.Pool
}

func NewTrackRepo(db *pgxpool.Pool) *TrackRepo {
	return &TrackRepo{
		db: db,
	}
}

func (tr *TrackRepo) GetTrackById(ctx context.Context, trackId string) (Track, error) {
	query := `
		SELECT t.*, r.isrc, r.duration_ms, r.popularity, r.play_count, r.audio_uri, r.preview_uri
		FROM tracks t
		JOIN recordings r ON t.recording_id = r.id
		WHERE t.id = $1
	`

	rows, err := tr.db.Query(ctx, query, trackId)
	if err != nil {
		log.Error().Err(err).Msg("err get track by id")
		return Track{}, err
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Track])
	if err != nil {
		log.Error().Err(err).Msg("err collect track by id")
		return Track{}, err
	}

	return data, nil
}

func (tr *TrackRepo) GetTracksByIds(ctx context.Context, trackIds []string) ([]Track, error) {
	query := `
		SELECT t.*, r.isrc, r.duration_ms, r.popularity, r.play_count, r.audio_uri, r.preview_uri
		FROM tracks t
		JOIN recordings r ON t.recording_id = r.id
		WHERE t.id = ANY($1::uuid[])
	`

	rows, err := tr.db.Query(ctx, query, trackIds)
	if err != nil {
		log.Error().Err(err).Msg("err get tracks by ids")
		return nil, err
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Track])
	if err != nil {
		log.Error().Err(err).Msg("err collect tracks by ids")
		return nil, err
	}

	return data, nil
}
