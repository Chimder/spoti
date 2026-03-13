package service

import (
	"context"
	"spoti/internal/domain/playlist"
	meilisearchrepo "spoti/internal/repository/meilisearch"
	"spoti/internal/repository/postgres"
	rediscache "spoti/internal/repository/redis"
	"time"
)

const playlistTTL = 60 * time.Minute

func userPlaylistsKey(id string) string { return "user:" + id + ":playlists" }

type PlaylistService struct {
	repo  *postgres.Repository
	cache rediscache.Cache
	meili *meilisearchrepo.MeiliRepository
}

func NewPlaylistService(repo *postgres.Repository, cache rediscache.Cache,
	meili *meilisearchrepo.MeiliRepository) *PlaylistService {
	return &PlaylistService{repo: repo, cache: cache, meili: meili}
}

func (ps *PlaylistService) CreatePlaylist(ctx context.Context, p playlist.CreatePlaylistReq) error {
	id, err := ps.repo.Playlist.CreatePlaylist(ctx, p)
	if err != nil {
		return err
	}
	err = ps.meili.Add(ctx, meilisearchrepo.Document{ID: id.String(), Type: "playlist", Name: p.PlaylistName})
	return err
}

func (ps *PlaylistService) GetPlaylistById(ctx context.Context, playlistID string, limit, offset int) (playlist.PlaylistJson, error) {
	return ps.repo.Playlist.GetPlaylistById(ctx, playlistID, limit, offset)
}

func (ps *PlaylistService) AddToPlaylist(ctx context.Context, playlist_id, track_id string) error {
	return ps.repo.Playlist.AddToPlaylist(ctx, playlist_id, track_id)
}

func (ps *PlaylistService) UpdatePlaylist(ctx context.Context, playlistId string, req playlist.UpdatePlaylistReq) error {
	return ps.repo.Playlist.UpdatePlaylist(ctx, playlistId, req)
}

func (ps *PlaylistService) DeleteFromPlaylist(ctx context.Context, playlistId, trackId string) error {
	return ps.repo.Playlist.DeleteFromPlaylist(ctx, playlistId, trackId)
}

func (ps *PlaylistService) GetAllUserPlaylists(ctx context.Context, userID string) ([]playlist.Playlist, error) {
	var playlists []playlist.Playlist
	if err := ps.cache.Get(ctx, userPlaylistsKey(userID), &playlists); err == nil {
		return playlists, nil
	}

	playlists, err := ps.repo.Playlist.GetAllUserPlaylists(ctx, userID)
	if err != nil {
		return nil, err
	}

	_ = ps.cache.Set(ctx, userPlaylistsKey(userID), playlists, playlistTTL)
	return playlists, nil
}
