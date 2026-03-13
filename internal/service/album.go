package service

import (
	"context"
	"encoding/json"
	"spoti/internal/domain/album"
	meilisearchrepo "spoti/internal/repository/meilisearch"
	"spoti/internal/repository/postgres"
	rediscache "spoti/internal/repository/redis"
	"time"
)

type AlbumService struct {
	repo  *postgres.Repository
	cache rediscache.Cache
	meili *meilisearchrepo.MeiliRepository
}

func NewAlbumService(r *postgres.Repository, cache rediscache.Cache,
	meili *meilisearchrepo.MeiliRepository) *AlbumService {
	return &AlbumService{repo: r, cache: cache, meili: meili}
}

func (as *AlbumService) CreateAlbum(ctx context.Context, a album.CreateAlbumReq) error {
	id, err := as.repo.Album.CreateAlbum(ctx, a)
	if err != nil {
		return err
	}

	err = as.meili.Add(ctx, meilisearchrepo.Document{ID: id.String(), Type: "album", Name: a.AlbumName})
	return err
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
	op := "GetUserSavedAlbums"
	var album []album.Album
	if err := as.cache.Get(ctx, userId+op, &album); err == nil {
		return album, nil
	}

	albumData, err := as.repo.Album.GetUserSavedAlbums(ctx, userId)
	if err != nil {
		return nil, err
	}

	_ = as.cache.Set(ctx, userId+op, albumData, 60*time.Minute)

	return albumData, nil
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
