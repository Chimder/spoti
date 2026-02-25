package service

import (
	"context"
	"encoding/json"
	"spoti/internal/domain/album"
	"spoti/internal/repository/postgres"
)

type AlbumService struct {
	repo *postgres.Repository
}

func NewAlbumService(r *postgres.Repository) *AlbumService {
	return &AlbumService{
		repo: r,
	}
}

func (as *AlbumService) CreateAlbum(ctx context.Context, a album.CreateAlbumReq) error {
	_, err := as.repo.Album.CreateAlbum(ctx, a)
	if err != nil {
		return err
	}
	return nil
}

func (as *AlbumService) GetAlbumJson(ctx context.Context, albumID string) (json.RawMessage, error) {
	return as.repo.Album.GetAlbumJson(ctx, albumID)
}

func (as *AlbumService) GetAlbumsByIds(ctx context.Context, albumIDs []string) (json.RawMessage, error) {
	return as.repo.Album.GetAlbumsByIds(ctx, albumIDs)
}

func (as *AlbumService) GetAlbumsTracks(ctx context.Context, albumID string) (json.RawMessage, error) {
	return as.repo.Album.GetAlbumsTracks(ctx, albumID)
}

func (as *AlbumService) GetUserSavedAlbums(ctx context.Context, userId string) ([]album.Album, error) {
	return as.repo.Album.GetUserSavedAlbums(ctx, userId)
}

func (as *AlbumService) SaveAlbumsForCurrentUser(ctx context.Context, albumIds []string, userId string) error {
	return as.repo.Album.SaveAlbumsForCurrentUser(ctx, albumIds, userId)
}

func (as *AlbumService) RemoveAlbumsFromCurrentUser(ctx context.Context, albumIds []string, userId string) error {
	return as.repo.Album.RemoveAlbumsFromCurrentUser(ctx, albumIds, userId)
}

func (as *AlbumService) CheckUsersSavedAlbums(ctx context.Context, albumIDs []string, userID string) ([]bool, error) {
	return as.repo.Album.CheckUsersSavedAlbums(ctx, albumIDs, userID)
}

func (as *AlbumService) GetNewReleases(ctx context.Context, limit int) ([]album.Album, error) {
	return as.repo.Album.GetNewReleases(ctx, limit)
}
