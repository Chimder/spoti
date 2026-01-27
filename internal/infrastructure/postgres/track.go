package postgres

import "github.com/google/uuid"

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
}
