package service

import (
	"context"
	"spoti/internal/domain/playlist"
	"spoti/internal/repository/postgres"
)

type PlaylistService struct {
	repo *postgres.Repository
}

func NewPlaylistService(repo *postgres.Repository) *PlaylistService {
	return &PlaylistService{
		repo: repo,
	}
}

func (ps *PlaylistService) CreatePlaylist(ctx context.Context, p playlist.CreatePlaylistReq) error {
	_, err := ps.repo.Playlist.CreatePlaylist(ctx, p)
	if err != nil {
		return err
	}
	return nil
}

func (ps *PlaylistService) GetPlaylistById(ctx context.Context, playlistId string, limit, offset int) (playlist.PlaylistJson, error) {
	return ps.repo.Playlist.GetPlaylistById(ctx, playlistId, limit, offset)
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

func (ps *PlaylistService) GetAllUserPlaylists(ctx context.Context, userId string) ([]playlist.Playlist, error) {
	return ps.repo.Playlist.GetAllUserPlaylists(ctx, userId)
}
