package clickhouse

import (
	"context"

	"github.com/Chimder/spoti/internal/domain/user"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
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

func (le *ListeningEventRepo) AddListeningEvent(ctx context.Context, e user.ListeningEventReq) error {
	err := le.db.Exec(ctx, `
				INSERT INTO spotify.listening_events (
					user_id, track_id, artist_id, album_id, duration_ms, is_skipped
				) VALUES (?, ?, ?, ?, ?, ?)
			`, e.UserId, e.TrackId, e.ArtistId, e.AlbumId, e.DurationMs, e.IsSkipped)
	if err != nil {
		log.Error().Err(err).Msg("err add listenign event to id")
		return err
	}

	return err
}

func (le *ListeningEventRepo) GetListeningEvent(ctx context.Context, e user.ListeningEventReq) error {
	err := le.db.Exec(ctx, `
				SELECT * FROM "spotify"."listening_events" WHERE user_id = ?
			`, e.UserId, e.TrackId, e.ArtistId, e.AlbumId, e.DurationMs, e.IsSkipped)

	if err != nil {
		log.Error().Err(err).Msg("err add listenign event to id")
		return err
	}

	return err
}
