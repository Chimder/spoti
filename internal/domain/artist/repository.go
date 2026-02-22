package artist

import (
	"context"

	"github.com/google/uuid"
)

type ArtistRepository interface {
	CreateArtist(ctx context.Context, a CreateArtistReq) (uuid.UUID, error)
	GetArtist(ctx context.Context, artistId string) (Artist, error)
	GetArtistsByIDs(ctx context.Context, artistIds []string) ([]Artist, error)
	GetArtistAlbums(ctx context.Context, artistId string) ([]Artist, error)
}
