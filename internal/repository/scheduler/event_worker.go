package scheduler

import (
	"context"
	"fmt"

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

type TrackDatas struct {
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
	countEvents := 0
	for {
		trackDatas, err := getTracksFromPostgres(ew.ctx, ew.db)
		if err != nil {
			log.Error().Err(err).Msg("err fetch postgres tracks")
			select {
			case <-ew.ctx.Done():
				return
			case <-time.After(10 * time.Minute):
				continue
			}
		}
		fmt.Printf("fetch %d tracks\n", len(trackDatas))

		boundary := time.Now().UTC()
		stats, err := getStatsFromClickhouse(ew.ctx, ew.clickDb, trackDatas, boundary)
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

		err = updateRecordingsBulk(ew.ctx, ew.db, stats, trackDatas)
		if err != nil {
			log.Error().Err(err).Msg("err update recor batch")
		}

		processedIDs := make([]string, len(stats))
		for i, s := range stats {
			processedIDs[i] = s.TrackID
		}

		err = deleteProcessedEvents(ew.ctx, ew.clickDb, processedIDs, boundary)
		if err != nil {
			log.Error().Err(err).Msg("Err delete CH events")
		}

		countEvents++
		fmt.Printf("End %d event \n", countEvents)

		select {
		case <-ew.ctx.Done():
			log.Info().Msg("Stop event worker")
			return
		case <-time.After(10 * time.Minute):
		}
	}
}

type TrackStats struct {
	TrackID         string `ch:"track_id"`
	TotalListenedMs uint64 `ch:"total_listened_ms"`
}

func getStatsFromClickhouse(ctx context.Context, db driver.Conn, tracksData []TrackDatas, boundary time.Time) ([]TrackStats, error) {
	trackIds := make([]string, len(tracksData))
	for i, t := range tracksData {
		trackIds[i] = t.TrackId.String()
	}

	var allStats []TrackStats
	batchSize := 500
	total := len(trackIds)

	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)
		batch := trackIds[i:end]

		rows, err := db.Query(ctx, `
            SELECT
                track_id,
                sum(duration_ms) AS total_listened_ms
            FROM spotify.listening_events
            WHERE track_id IN (?)
              AND created_at <= ?
            GROUP BY track_id
        `, batch, boundary)
		if err != nil {
			return nil, fmt.Errorf("batch %d-%d query err: %w", i, end, err)
		}

		for rows.Next() {
			var s TrackStats
			if err := rows.ScanStruct(&s); err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan err: %w", err)
			}
			allStats = append(allStats, s)
		}
		rows.Close()

		if err := rows.Err(); err != nil {
			return nil, fmt.Errorf("rows err: %w", err)
		}
	}

	return allStats, nil
}

func updateRecordingsBulk(ctx context.Context, pool *pgxpool.Pool, stats []TrackStats, trackDatas []TrackDatas) error {
	durationMap := make(map[string]int, len(trackDatas))
	for _, t := range trackDatas {
		durationMap[t.TrackId.String()] = t.DurationMs
	}

	trackIDs := make([]uuid.UUID, 0, len(stats))
	playCounts := make([]int64, 0, len(stats))

	for _, s := range stats {
		durationMs, ok := durationMap[s.TrackID]
		if !ok || durationMs == 0 {
			continue
		}

		pc := int64(s.TotalListenedMs) / int64(durationMs)
		if pc == 0 {
			continue
		}

		id, err := uuid.Parse(s.TrackID)
		if err != nil {
			log.Error().Err(err).Str("track_id", s.TrackID).Msg("bad track id")
			continue
		}

		trackIDs = append(trackIDs, id)
		playCounts = append(playCounts, pc)
	}

	if len(trackIDs) == 0 {
		return nil
	}

	tag, err := pool.Exec(ctx, `
		UPDATE recordings r
		SET play_count = r.play_count + v.play_count
		FROM unnest($1::uuid[], $2::bigint[]) AS v(track_id, play_count),
		tracks t
		WHERE t.id = v.track_id
		  AND t.recording_id = r.id
	`, trackIDs, playCounts)
	if err != nil {
		return fmt.Errorf("err bulk update: %w", err)
	}

	log.Info().Int64("rows", tag.RowsAffected()).Msg("recordings upd")
	return nil
}

func getTracksFromPostgres(ctx context.Context, pool *pgxpool.Pool) ([]TrackDatas, error) {
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

	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[TrackDatas])
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

func deleteProcessedEvents(ctx context.Context, db driver.Conn, trackIDs []string, boundary time.Time) error {
	batchSize := 500
	total := len(trackIDs)

	for i := 0; i < total; i += batchSize {
		end := min(i+batchSize, total)
		batch := trackIDs[i:end]

		err := db.Exec(ctx, `
            ALTER TABLE spotify.listening_events
            DELETE WHERE track_id IN (?) AND created_at <= ?
        `, batch, boundary)
		if err != nil {
			return fmt.Errorf("err delete events: %w", err)
		}

	}

	return nil
}

func (ew *EventWorker) Stop() {
	ew.cancel()
	time.Sleep(500 * time.Millisecond)
	log.Info().Msg("Stop EventWorker.")
}
