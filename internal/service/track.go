package service

import (
	"context"
	"spoti/internal/domain/track"
	"spoti/internal/repository/postgres"
)

type TrackService struct {
	repo *postgres.Repository
}

func NewTrackService(repo *postgres.Repository) *TrackService {
	return &TrackService{
		repo: repo,
	}
}

func (ts *TrackService) CreateTrack(ctx context.Context, t track.CreateTrackReq) error {

	_, err := ts.repo.Track.CreateTrack(ctx, t)
	if err != nil {
		return err
	}
	return nil
}

func (ts *TrackService) GetTrackById(ctx context.Context, trackId string) (track.Track, error) {
	return ts.repo.Track.GetTrackById(ctx, trackId)
}

func (ts *TrackService) GetTracksByIds(ctx context.Context, trackIds []string) ([]track.Track, error) {
	return ts.repo.Track.GetTracksByIds(ctx, trackIds)
}

func (ts *TrackService) GetArtistTracks(ctx context.Context, artistId string) ([]track.Track, error) {
	return ts.repo.Track.GetArtistTracks(ctx, artistId)
}
