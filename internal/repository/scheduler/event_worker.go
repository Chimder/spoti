package scheduler

import (
	"context"
	"fmt"
	"sync"

	"time"

	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type EventWorker struct {
	db        *pgxpool.Pool
	maxWorker int
	ctx       context.Context
	clickDb   driver.Conn
	cancel    context.CancelFunc
}

type TracksData struct {
	TrackId    uuid.UUID `db:"track_id"`
	AlbumId    uuid.UUID `db:"album_id"`
	DurationMs int       `db:"duration_ms"`
	ArtistId   uuid.UUID `db:"artist_id"`
}

func NewEventWorker(ctx context.Context, db *pgxpool.Pool, clickDb driver.Conn) *EventWorker {
	ctx, cancel := context.WithCancel(ctx)
	return &EventWorker{
		db:        db,
		clickDb:   clickDb,
		maxWorker: 100,
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (ew *EventWorker) Start() {
	log.Info().Msg("Starting task manager")
	go ew.EventWorker()
}

func (ew *EventWorker) EventWorker() {
	for {
		// userIds, err := getUsersFromPostgres(ew.ctx, ew.db.Pool)
		tracksData, err := getTracksFromPostgres(ew.ctx, ew.db)
		if err != nil {
			log.Error().Err(err).Msg("err fetch postgres id for event")
			select {
			case <-ew.ctx.Done():
				return
			case <-time.After(2 * time.Minute):
				continue
			}
		}
		fmt.Printf("fetch %d users id", len(tracksData))

		stats, err := getStatsFromClickhouse(ew.ctx, ew.clickDb, tracksData)
		if err != nil {
			log.Error().Err(err).Msg("err fetch clickhouse stats")
			select {
			case <-ew.ctx.Done():
				return
			case <-time.After(2 * time.Minute):
				continue
			}
		}

		fmt.Printf("Get all stats %d\n", len(stats))
		trackDurationMap := make(map[string]int, len(tracksData))
		for _, t := range tracksData {
			trackDurationMap[t.TrackId.String()] = t.DurationMs
		}

		sem := make(chan struct{}, 300)
		var wg sync.WaitGroup
		for i, s := range stats {
			sem <- struct{}{}
			wg.Add(1)

			go func(i int, trackStats TrackStats) {
				defer func() {
					<-sem
					wg.Done()
				}()

				durationMs, ok := trackDurationMap[trackStats.TrackID]
				if !ok || durationMs == 0 {
					return
				}

				playCount := int64(trackStats.TotalListenedMs) / int64(durationMs)

				err := updateRecording(ew.ctx, ew.db, trackStats.TrackID, playCount)
				if err != nil {
					log.Error().Err(err).Msg("Err update recording")
					return
				}
			}(i, s)
		}
		wg.Wait()

		select {
		case <-ew.ctx.Done():
			// close(sem)
			log.Info().Msg("Stopping processing")
			return
		case <-time.After(15 * time.Minute):
		}
	}
}

type TrackStats struct {
	TrackID         string `ch:"track_id"`
	TotalListenedMs uint64 `ch:"total_listened_ms"`
}

func getStatsFromClickhouse(ctx context.Context, db driver.Conn, tracksData []TracksData) ([]TrackStats, error) {
	trackIds := make([]string, len(tracksData))
	for i, t := range tracksData {
		trackIds[i] = t.TrackId.String()
	}

	var allStats []TrackStats
	batchSize := 500
	total := len(trackIds)

	for i := 0; i < total; i += batchSize {
		end := i + batchSize
		if end > total {
			end = total
		}
		batch := trackIds[i:end]

		log.Info().Msgf("CH query batch %d-%d of %d", i, end, total)

		rows, err := db.Query(ctx, `
            SELECT
                track_id,
                sum(duration_ms) AS total_listened_ms
            FROM spotify.listening_events
            WHERE track_id IN (?)
            GROUP BY track_id
        `, batch)
		if err != nil {
			return nil, fmt.Errorf("batch %d-%d query err: %w", i, end, err)
		}

		batchCount := 0
		for rows.Next() {
			var s TrackStats
			if err := rows.ScanStruct(&s); err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan err: %w", err)
			}
			log.Info().Msgf("track_id=%s total_listened_ms=%d", s.TrackID, s.TotalListenedMs)
			allStats = append(allStats, s)
			batchCount++
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows err: %w", err)
		}

		log.Info().Msgf("batch %d-%d got %d rows from CH", i, end, batchCount)
	}

	log.Info().Msgf("total CH stats: %d", len(allStats))
	return allStats, nil
}

func updateRecording(ctx context.Context, pool *pgxpool.Pool, trackId string, playCount int64) error {
	id, err := uuid.Parse(trackId)
	if err != nil {
		return fmt.Errorf("err parse track id: %w", err)
	}

	query := `
        UPDATE recordings r
        SET play_count = $1
        FROM tracks t
        WHERE t.recording_id = r.id
          AND t.id = $2
    `

	tag, err := pool.Exec(ctx, query, playCount, id)
	if err != nil {
		return err
	}

	log.Info().Str("track_id", trackId).Int64("play_count", playCount).Int64("rows_affected", tag.RowsAffected()).Msg("recording updated")
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

func (ew *EventWorker) Stop() {
	ew.cancel()
	time.Sleep(500 * time.Millisecond)
	log.Info().Msg("Stop EventWorker.")
}
