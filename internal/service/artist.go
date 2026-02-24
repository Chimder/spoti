package service

import (
	"context"
	"spoti/internal/domain/artist"
	"spoti/internal/repository/postgres"
)

type ArtistService struct {
	repo *postgres.Repository
}

func NewArtistService(r *postgres.Repository) *ArtistService {
	return &ArtistService{repo: r}
}
func (as *ArtistService) CreateArtist(ctx context.Context, a artist.CreateArtistReq) error {
	_, err := as.repo.Artist.CreateArtist(ctx, a)
	if err != nil {
		return err
	}
	return nil
}

func (as *ArtistService) GetArtist(ctx context.Context, artistId string) (artist.Artist, error) {
	return as.repo.Artist.GetArtist(ctx, artistId)
}

func (as *ArtistService) GetArtistsByIDs(ctx context.Context, artistIds []string) ([]artist.Artist, error) {
	return as.repo.Artist.GetArtistsByIDs(ctx, artistIds)
}

func (as *ArtistService) GetArtistAlbums(ctx context.Context, artistId string) ([]artist.Artist, error) {
	return as.repo.Artist.GetArtistAlbums(ctx, artistId)
}
