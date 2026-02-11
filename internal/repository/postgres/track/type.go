package postgres

import (
	"spoti/internal/domain/track"

	"github.com/google/uuid"
)

type Track struct {
	ID          uuid.UUID `db:"id"`
	AlbumID     uuid.UUID `db:"album_id"`
	RecordingID uuid.UUID `db:"recording_id"`
	TrackName   string    `db:"track_name"`
	TrackNumber int16     `db:"track_number"`
	DiscNumber  int16     `db:"disc_number"`
	Explicit    bool      `db:"explicit"`
	IsPlayable  bool      `db:"is_playable"`
	TrackType   string    `db:"track_type"`
	URI         string    `db:"uri"`
	IsLocal     bool      `db:"islocal"`

	ISRC       string `db:"isrc"`
	DurationMs int64  `db:"duration_ms"`
	Popularity int    `db:"popularity"`
	PlayCount  int64  `db:"play_count"`
	AudioURI   string `db:"audio_uri"`
}

func (t *Track) ToDomain() track.Track {
	return track.Track{
		ID:          t.ID,
		AlbumID:     t.AlbumID,
		RecordingID: t.RecordingID,
		Name:        t.TrackName,
		Number:      t.TrackNumber,
		DiscNumber:  t.DiscNumber,
		Explicit:    t.Explicit,
		IsPlayable:  t.IsPlayable,
		Type:        t.TrackType,
		URI:         t.URI,
		IsLocal:     t.IsLocal,
		ISRC:        t.ISRC,
		DurationMs:  t.DurationMs,
		Popularity:  t.Popularity,
		PlayCount:   t.PlayCount,
		AudioURI:    t.AudioURI,
	}
}

func TracksToDomain(rows []Track) []track.Track {
	result := make([]track.Track, len(rows))
	for _, row := range rows {
		result = append(result, row.ToDomain())
	}
	return result
}
