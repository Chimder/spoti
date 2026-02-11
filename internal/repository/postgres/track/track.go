package postgres

import (
	"context"
	"spoti/internal/domain/track"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type TrackRepo struct {
	db *pgxpool.Pool
}

func NewTrackRepo(db *pgxpool.Pool) *TrackRepo {
	return &TrackRepo{
		db: db,
	}
}

func (tr *TrackRepo) GetTrackById(ctx context.Context, trackId string) (track.Track, error) {
	query := `
		SELECT t.*, r.isrc, r.duration_ms, r.popularity, r.play_count, r.audio_uri, r.preview_uri
		FROM tracks t
		JOIN recordings r ON t.recording_id = r.id
		WHERE t.id = $1
	`

	rows, err := tr.db.Query(ctx, query, trackId)
	if err != nil {
		log.Error().Err(err).Msg("err get track by id")
		return track.Track{}, err
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Track])
	if err != nil {
		log.Error().Err(err).Msg("err collect track by id")
		return track.Track{}, err
	}

	return data.ToDomain(), nil
}

func (tr *TrackRepo) GetTracksByIds(ctx context.Context, trackIds []string) ([]track.Track, error) {
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

	return TracksToDomain(data), nil
}

func (tr *TrackRepo) GetArtistTracks(ctx context.Context, artistId string) ([]track.Track, error) {
	query := `
		SELECT t.*, r.isrc, r.duration_ms, r.popularity, r.play_count, r.audio_uri, r.preview_uri
		FROM tracks t
		JOIN artist_tracks at ON at.track_id = t.id
		JOIN recordings r ON t.recording_id = r.id
		WHERE at.artist_id = $1
	`
	rows, err := tr.db.Query(ctx, query, artistId)
	if err != nil {
		log.Error().Err(err).Msg("err get tracks by artist id")
		return nil, err
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Track])
	if err != nil {
		log.Error().Err(err).Msg("err collect tracks by artist id")
		return nil, err
	}

	return TracksToDomain(data), nil
}
