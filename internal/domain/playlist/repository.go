package playlist

import (
	"context"

	"github.com/google/uuid"
)

type PlaylistRepository interface {
	CreatePlaylist(ctx context.Context, p CreatePlaylistReq) (uuid.UUID, error)
	GetPlaylistById(ctx context.Context, playlistId string, limit, offset int) (PlaylistJson, error)
	AddToPlaylist(ctx context.Context, playlist_id, track_id string) error
	UpdatePlaylist(ctx context.Context, playlistId string, req UpdatePlaylistReq) error
	DeleteFromPlaylist(ctx context.Context, playlistId, trackId string) error
	GetAllUserPlaylists(ctx context.Context, userId string) ([]Playlist, error)
}
