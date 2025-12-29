package user

import (
	"time"

	"github.com/google/uuid"
)

type ListeningEvent struct {
	Id         uuid.UUID
	UserId     uuid.UUID
	TrackId    uuid.UUID
	StartedAt  time.Time
	DurationMs uint32
	Skipped    bool
}
