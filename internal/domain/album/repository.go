package album

import (
	"context"

	"github.com/google/uuid"
)

type AlbumRepository interface {
	CreateAlbum(ctx context.Context, a CreateAlbumReq) (uuid.UUID, error)
	GetAlbum(ctx context.Context, albumID string) (Album, error)
	GetAlbumWithTracks(ctx context.Context, albumID string) (GetAlbumResponse, error)
	GetAlbumsByIds(ctx context.Context, albumIDs []string) (GetAlbumsByIdsResponse, error)
	GetUserSavedAlbums(ctx context.Context, userId string) ([]Album, error)
	SaveAlbumsForCurrentUser(ctx context.Context, albumIds []string, userId string) error
	RemoveAlbumsFromCurrentUser(ctx context.Context, albumIds []string, userId string) error
	CheckUsersSavedAlbums(ctx context.Context, albumIDs []string, userID string) ([]bool, error)
	GetNewReleases(ctx context.Context, limit int) ([]Album, error)
	AddArtistToAlbum(ctx context.Context, albumID, artistID uuid.UUID) error
}
