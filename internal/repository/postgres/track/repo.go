package trackrepo

import (
	"context"
	"spoti/internal/domain/track"
	"spoti/internal/repository/postgres/pgiface"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type TrackRepo struct {
	db pgiface.Querier
}

func NewTrackRepo(db pgiface.Querier) *TrackRepo {
	return &TrackRepo{
		db: db,
	}
}

func (tr *TrackRepo) CreateTrack(ctx context.Context, t track.CreateTrackReq) (uuid.UUID, error) {
	query := `
	INSERT INTO tracks (
	album_id, recording_id, track_name, track_number, disc_number, explicit, is_playable, track_type, uri, islocal
	)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	RETURNING id
	`

	var id uuid.UUID
	err := tr.db.QueryRow(ctx, query, t.AlbumId, t.RecordingId, t.Name, t.Number,
		t.DiscNumber, t.Explicit, t.IsPlayable, t.Type, t.URI, t.IsLocal).Scan(&id)
	if err != nil {
		return uuid.UUID{}, err
	}

	return id, err
}

func (tr *TrackRepo) GetTrackById(ctx context.Context, trackId string) (track.Track, error) {
	query := `
		SELECT t.*, r.isrc, r.duration_ms, r.popularity, r.play_count, r.audio_uri
		FROM tracks t
		JOIN recordings r ON t.recording_id = r.id
		WHERE t.id = $1
	`
	// query := `
	// 	SELECT * FROM tracks WHERE id = $1
	// `

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
	// SELECT * FROM tracks WHERE id = ANY($1::uuid[])
	query := `
		SELECT t.*, r.isrc, r.duration_ms, r.popularity, r.play_count, r.audio_uri
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
		SELECT t.*, r.isrc, r.duration_ms, r.popularity, r.play_count, r.audio_uri
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
