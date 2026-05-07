package recordingrepo

import (
	"context"

	"github.com/Chimder/spoti/internal/domain/recording"
	"github.com/Chimder/spoti/internal/repository/postgres/pgiface"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)
type RecordingRepository interface {
	CreateRecording(ctx context.Context, r recording.CreateRecordingReq) (uuid.UUID, error)
	GetRecordingById(ctx context.Context, recordingId string) (recording.Recording, error)
}
type RecordingRepo struct {
	db pgiface.Querier
}

func NewRecordingRepo(db pgiface.Querier) *RecordingRepo {
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
		log.Error().Err(err).Msg("err get recording by id")
		return recording.Recording{}, err
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[RecordingDB])
	if err != nil {
		log.Error().Err(err).Msg("err collect recording by id")
		return recording.Recording{}, err
	}

	return data.ToDomain(), nil
}
