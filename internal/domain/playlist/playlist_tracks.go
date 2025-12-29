package playlist

import (
	"spoti/internal/domain/track"

	"github.com/google/uuid"
)

type PlayListTrack struct {
	// id      uuid.UUID
	AddedAt    string
	AddedBy    uuid.UUID
	PlaylistId uuid.UUID
	Track track.Track
}
