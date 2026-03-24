package service

import (
	"context"
	"time"

	"github.com/Chimder/spoti/internal/domain/artist"
	meilisearchrepo "github.com/Chimder/spoti/internal/repository/meilisearch"
	"github.com/Chimder/spoti/internal/repository/postgres"
	rediscache "github.com/Chimder/spoti/internal/repository/redis"
	"github.com/google/uuid"
)

const artistTTL = 60 * time.Minute

func artistKey(id string) string       { return "artist:" + id }
func artistAlbumsKey(id string) string { return "artist:" + id + ":albums" }

type ArtistService struct {
	repo  *postgres.Repository
	cache rediscache.Cache
	meili *meilisearchrepo.MeiliRepository
}

func NewArtistService(r *postgres.Repository, cache rediscache.Cache,
	meili *meilisearchrepo.MeiliRepository) *ArtistService {
	return &ArtistService{repo: r, cache: cache, meili: meili}
}
func (as *ArtistService) CreateArtist(ctx context.Context, a artist.CreateArtistReq) (uuid.UUID, error) {
	id, err := as.repo.Artist.CreateArtist(ctx, a)
	if err != nil {
		return uuid.Nil, err
	}
	err = as.meili.Add(ctx, meilisearchrepo.Document{ID: id.String(), Type: "artist", Name: a.ArtistName})
	return id, err
}

func (as *ArtistService) GetArtist(ctx context.Context, artistID string) (artist.Artist, error) {
	var a artist.Artist
	if err := as.cache.Get(ctx, artistKey(artistID), &a); err == nil {
		return a, nil
	}

	a, err := as.repo.Artist.GetArtist(ctx, artistID)
	if err != nil {
		return artist.Artist{}, err
	}

	_ = as.cache.Set(ctx, artistKey(artistID), a, artistTTL)
	return a, nil
}

func (as *ArtistService) GetArtistsByIDs(ctx context.Context, artistIds []string) ([]artist.Artist, error) {
	return as.repo.Artist.GetArtistsByIDs(ctx, artistIds)
}

func (as *ArtistService) GetArtistAlbums(ctx context.Context, artistID string) ([]artist.Artist, error) {
	var albums []artist.Artist
	if err := as.cache.Get(ctx, artistAlbumsKey(artistID), &albums); err == nil {
		return albums, nil
	}

	albums, err := as.repo.Artist.GetArtistAlbums(ctx, artistID)
	if err != nil {
		return nil, err
	}

	_ = as.cache.Set(ctx, artistAlbumsKey(artistID), albums, artistTTL)
	return albums, nil
}
