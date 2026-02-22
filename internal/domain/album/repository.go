package album

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type AlbumRepository interface {
	CreateAlbum(ctx context.Context, a CreateAlbumReq) (uuid.UUID, error)
	GetAlbum(ctx context.Context, albumID string) (json.RawMessage, error)
	GetAlbumsByIds(ctx context.Context, albumIDs []string) (json.RawMessage, error)
	GetAlbumsTracks(ctx context.Context, albumID string) (json.RawMessage, error)
	GetUserSavedAlbums(ctx context.Context, userId string) ([]Album, error)
	SaveAlbumsForCurrentUser(ctx context.Context, albumIds []string, userId string) error
	RemoveAlbumsFromCurrentUser(ctx context.Context, albumIds []string, userId string) error
	CheckUsersSavedAlbums(ctx context.Context, albumIDs []string, userID string) ([]bool, error)
	GetNewReleases(ctx context.Context, limit int) ([]Album, error)
}
