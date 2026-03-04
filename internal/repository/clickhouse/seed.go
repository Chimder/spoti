package clickhouse

import (
	"context"
	"fmt"
	"log"
	"math/rand"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ListeningEvent struct {
	userID     uuid.UUID
	trackID    uuid.UUID
	artistID   uuid.UUID
	albumID    uuid.UUID
	durationMS uint32
}

var (
	ListeningEventsCount = 50000
	DaysToGenerate       = 90
)

func SeedListeningEvents(ctx context.Context, pool *pgxpool.Pool) error {

	conn, err := connToClick()
	if err != nil {
		return fmt.Errorf("err connection to clickhouse: %w", err)
	}
	defer conn.Close()

	batch, err := conn.PrepareBatch(ctx, `
		INSERT INTO spotify.listening_events (user_id, track_id, artist_id, album_id, duration_ms, is_skipped)
	`)
	if err != nil {
		return fmt.Errorf("err to prepareBatch: %w", err)
	}

	tracks, err := getTracksFromPostgres(ctx, pool)
	if err != nil {
		return fmt.Errorf("err to getTracksFromPostgres: %w", err)
	}

	ids, err := getUsersFromPostgres(ctx, pool)
	if err != nil {
		return fmt.Errorf("err to getTracksFromPostgres: %w", err)
	}

	eventsCount := 0

	for _, t := range tracks {

		eventsRand := rand.Intn(500)
		for range eventsRand {
			userID := ids[rand.Intn(len(ids))]

			listenPercent := 0.3 + rand.Float64()*0.7
			listenDuration := uint32(float64(t.DurationMs) * listenPercent)
			isSkipped := listenPercent < 0.8

			err := batch.Append(
				userID.String(),
				t.TrackId.String(),
				t.ArtistId.String(),
				t.AlbumId.String(),
				listenDuration,
				isSkipped,
			)
			if err != nil {
				return fmt.Errorf("err to append clickhouse: %w", err)
			}
			eventsCount++
		}
	}

	fmt.Printf("Generated event clickhouse %d\n", eventsCount)

	if err := batch.Send(); err != nil {
		return fmt.Errorf("err to send batch: %w", err)
	}

	return nil
}

func getTracksFromPostgres(ctx context.Context, pool *pgxpool.Pool) ([]TracksData, error) {
	query := `
		SELECT
			t.id as track_id,
			t.album_id,
			r.duration_ms,
			(SELECT artist_id FROM artist_tracks WHERE track_id = t.id LIMIT 1) as artist_id
		FROM tracks t
		JOIN recordings r ON t.recording_id = r.id
		WHERE r.duration_ms > 0
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[TracksData])
	if err != nil {
		return nil, err
	}

	return data, nil
}

func getUsersFromPostgres(ctx context.Context, pool *pgxpool.Pool) ([]uuid.UUID, error) {
	rows, err := pool.Query(ctx, `SELECT id FROM users`)
	if err != nil {
		log.Printf("Err ch query users id")
		return nil, err
	}

	ids, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		log.Printf("Err ch collect users id")
		return nil, err
	}

	return ids, nil
}

type TracksData struct {
	TrackId    uuid.UUID `db:"track_id"`
	AlbumId    uuid.UUID `db:"album_id"`
	DurationMs int       `db:"duration_ms"`
	ArtistId   uuid.UUID `db:"artist_id"`
}

func connToClick() (driver.Conn, error) {
	return clickhouse.Open(&clickhouse.Options{
		Addr: []string{"localhost:9000"},
		Auth: clickhouse.Auth{
			Database: "default",
			Username: "default",
			Password: "",
		},
		Compression: &clickhouse.Compression{
			Method: clickhouse.CompressionLZ4,
		},
		ClientInfo: clickhouse.ClientInfo{
			Products: []struct {
				Name    string
				Version string
			}{
				{Name: "go-client-spoti", Version: "0.1"},
			},
		},
		Debugf: func(format string, v ...any) {
			fmt.Printf("[ClickHouse DEBUG] "+format+"\n", v...)
		},
	})
}
