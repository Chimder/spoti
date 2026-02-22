package track

import (
	"context"

	"github.com/google/uuid"
)

type TrackRepository interface {
	CreateTrack(ctx context.Context, t CreateTrackReq) (uuid.UUID, error)
	GetTrackById(ctx context.Context, trackId string) (Track, error)
	GetTracksByIds(ctx context.Context, trackIds []string) ([]Track, error)
	GetArtistTracks(ctx context.Context, artistId string) ([]Track, error)
}
