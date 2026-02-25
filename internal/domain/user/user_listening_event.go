package user

import (
	"github.com/google/uuid"
)

type ListeningEventReq struct {
	UserId     uuid.UUID
	TrackId    uuid.UUID
	AlbumId    uuid.UUID
	ArtistId   uuid.UUID
	DurationMs int
	IsSkipped  bool
}
