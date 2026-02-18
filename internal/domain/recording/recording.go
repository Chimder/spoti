package recording

import (
	"time"

	"github.com/google/uuid"
)

type Recording struct {
	Id         uuid.UUID
	ISRC       string
	DurationMs int64
	Popularity int
	PlayCount  int64
	AudioURI   string
	CreatedAt  time.Time
}
