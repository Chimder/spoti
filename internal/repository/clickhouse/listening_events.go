package clickhouse

import (
	"context"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type ListeningEventRepo struct {
	db driver.Conn
}

func NewListeningEventRepo(db driver.Conn) *ListeningEventRepo {
	return &ListeningEventRepo{
		db: db,
	}
}

type AddEventReq struct {
	UserId     uuid.UUID
	TrackId    uuid.UUID
	AlbumId    uuid.UUID
	DurationMs int
	ArtistId   uuid.UUID
	isSkipped  bool
}

func (le *ListeningEventRepo) AddListeningEvent(ctx context.Context, e AddEventReq) error {

	err := le.db.Exec(ctx, `
				INSERT INTO spotify.listening_events (
					user_id, track_id, artist_id, album_id, duration_ms, is_skipped
				) VALUES (?, ?, ?, ?, ?, ?)
			`, e.UserId, e.TrackId, e.ArtistId, e.AlbumId, e.DurationMs, e.isSkipped)
	if err != nil {
		log.Error().Err(err).Msg("err add listenign event to id")
		return err
	}

	return err
}
