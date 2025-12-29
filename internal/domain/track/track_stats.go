package track

import (
	"time"

	"github.com/google/uuid"
)

type TrackStats struct {
	id              uuid.UUID
	trackId         uuid.UUID
	totalListens    uint64
	uniqueListeners uint64
	totalMinutes    uint64
	updatedAt       time.Time
}
