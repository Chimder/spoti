package service

import (
	"context"
	"spoti/internal/domain/track"
	"spoti/internal/repository/postgres"
	rediscache "spoti/internal/repository/redis"
	"time"

	"github.com/google/uuid"
)

const trackTTL = 60 * time.Minute

func trackKey(id uuid.UUID) string     { return "track:" + id.String() }
func artistTracksKey(id string) string { return "artist:" + id + ":tracks" }

type TrackService struct {
	repo  *postgres.Repository
	cache rediscache.Cache
}

func NewTrackService(repo *postgres.Repository, cache rediscache.Cache) *TrackService {
	return &TrackService{repo: repo, cache: cache}
}

func (ts *TrackService) CreateTrack(ctx context.Context, t track.CreateTrackReq) error {
	_, err := ts.repo.Track.CreateTrack(ctx, t)
	return err
}

func (ts *TrackService) GetTrackById(ctx context.Context, trackID uuid.UUID) (track.Track, error) {
	var t track.Track
	if err := ts.cache.Get(ctx, trackKey(trackID), &t); err == nil {
		return t, nil
	}

	t, err := ts.repo.Track.GetTrackById(ctx, trackID)
	if err != nil {
		return track.Track{}, err
	}

	_ = ts.cache.Set(ctx, trackKey(trackID), t, trackTTL)
	return t, nil
}

func (ts *TrackService) GetTracksByIds(ctx context.Context, trackIds []string) ([]track.Track, error) {
	return ts.repo.Track.GetTracksByIds(ctx, trackIds)
}

func (ts *TrackService) GetArtistTracks(ctx context.Context, artistID string) ([]track.Track, error) {
	var tracks []track.Track
	if err := ts.cache.Get(ctx, artistTracksKey(artistID), &tracks); err == nil {
		return tracks, nil
	}

	tracks, err := ts.repo.Track.GetArtistTracks(ctx, artistID)
	if err != nil {
		return nil, err
	}

	_ = ts.cache.Set(ctx, artistTracksKey(artistID), tracks, trackTTL)
	return tracks, nil
}
