package postgres

import (
	"spoti/internal/domain/recording"
	"time"

	"github.com/google/uuid"
)

type RecordingDB struct {
	ID         uuid.UUID `db:"id"`
	ISRC       string    `db:"isrc"`
	DurationMs int64     `db:"duration_ms"`
	Popularity int       `db:"popularity"`
	PlayCount  int64     `db:"play_count"`
	AudioURI   string    `db:"audio_uri"`
	CreatedAt  time.Time `db:"created_at"`
}

func (r *RecordingDB) ToDomain() recording.Recording {
	return recording.Recording{
		ID:         r.ID,
		ISRC:       r.ISRC,
		DurationMs: r.DurationMs,
		Popularity: r.Popularity,
		PlayCount:  r.PlayCount,
		AudioURI:   r.AudioURI,
		CreatedAt:  r.CreatedAt,
	}
}

func RecordingsToDomain(rows []RecordingDB) []recording.Recording {
	result := make([]recording.Recording, len(rows))
	for i, row := range rows {
		result[i] = row.ToDomain()
	}
	return result
}
