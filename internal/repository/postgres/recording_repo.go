package postgres

import (
	"context"
	"spoti/internal/domain/recording"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type RecordingRepo struct {
	db *pgxpool.Pool
}

func NewRecordingRepo(db *pgxpool.Pool) *RecordingRepo {
	return &RecordingRepo{
		db: db,
	}
}
func (rc *RecordingRepo) CreateRecording(ctx context.Context, r recording.CreateRecordingReq) (uuid.UUID, error) {
	query := `
			INSERT INTO recordings (isrc, duration_ms, audio_uri)
			VALUES ($1, $2, $3)
			RETURNING id
		`

	var id uuid.UUID
	err := rc.db.QueryRow(ctx, query, r.ISRC, r.DurationMs, r.AudioUri).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("err create recording")
		return uuid.UUID{}, err
	}

	return id, err
}

func (rc *RecordingRepo) GetRecordingById(ctx context.Context, recordingId string) (recording.Recording, error) {
	query := `SELECT * FROM recordings WHERE id = $1`

	rows, err := rc.db.Query(ctx, query, recordingId)
	if err != nil {
		log.Error().Err(err).Msg("err get track by id")
		return recording.Recording{}, err
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[RecordingDB])
	if err != nil {
		log.Error().Err(err).Msg("err collect track by id")
		return recording.Recording{}, err
	}

	return data.ToDomain(), nil
}
