package grpc

import (
	"context"
	"fmt"

	"github.com/Chimder/spoti/internal/domain/track"
	trackv1 "github.com/Chimder/spoti/internal/gen/track/v1"
	"github.com/Chimder/spoti/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TrackHandler struct {
	trackv1.UnimplementedTrackServiceServer
	srv *service.TrackService
}

func NewTrackHandler(srv *service.TrackService) *TrackHandler {
	return &TrackHandler{srv: srv}
}

func (h *TrackHandler) CreateTrack(ctx context.Context, req *trackv1.CreateTrackRequest) (*trackv1.CreateTrackResponse, error) {
	albumID, err := uuid.Parse(req.AlbumId)
	if err != nil {
		return nil, fmt.Errorf("invalid album_id: %w", err)
	}
	recordingID, err := uuid.Parse(req.RecordingId)
	if err != nil {
		return nil, fmt.Errorf("invalid recording_id: %w", err)
	}

	id, err := h.srv.CreateTrack(ctx, track.CreateTrackReq{
		AlbumId:     albumID,
		RecordingId: recordingID,
		Name:        req.Name,
		Number:      int16(req.Number),
		DiscNumber:  int16(req.DiscNumber),
		Explicit:    req.Explicit,
		IsPlayable:  req.IsPlayable,
		Type:        req.Type,
		URI:         req.Uri,
		IsLocal:     req.IsLocal,
	})
	if err != nil {
		return nil, status.Error(codes.NotFound, "")
	}

	return &trackv1.CreateTrackResponse{Id: id.String()}, nil
}

func (h *TrackHandler) GetTrackById(ctx context.Context, req *trackv1.GetTrackByIdRequest) (*trackv1.GetTrackByIdResponse, error) {
	trackID, err := uuid.Parse(req.TrackId)
	if err != nil {
		return nil, fmt.Errorf("invalid Track id")
	}

	tr, err := h.srv.GetTrackById(ctx, trackID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found track")
	}
	return &trackv1.GetTrackByIdResponse{Track: trackToProto(tr)}, nil
}

func (h *TrackHandler) GetTracksByIds(ctx context.Context, req *trackv1.GetTracksByIdsRequest) (*trackv1.GetTracksByIdsResponse, error) {

	tracks, err := h.srv.GetTracksByIds(ctx, req.TrackIds)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found tracks")
	}
	return &trackv1.GetTracksByIdsResponse{Tracks: tracksToProto(tracks)}, nil
}

func (h *TrackHandler) GetArtistTracks(ctx context.Context, req *trackv1.GetArtistTracksRequest) (*trackv1.GetArtistTracksResponse, error) {
	tracks, err := h.srv.GetArtistTracks(ctx, req.ArtistId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found artist tracks")
	}

	return &trackv1.GetArtistTracksResponse{Tracks: tracksToProto(tracks)}, status.Error(codes.Unimplemented, "method GetArtistTracks not implemented")
}

func trackToProto(t track.Track) *trackv1.Track {
	return &trackv1.Track{
		Id:          t.Id.String(),
		AlbumId:     t.AlbumId.String(),
		RecordingId: t.RecordingId.String(),
		Name:        t.Name,
		Number:      int32(t.Number),
		DiscNumber:  int32(t.DiscNumber),
		Explicit:    t.Explicit,
		IsPlayable:  t.IsPlayable,
		Type:        t.Type,
		Uri:         t.URI,
		AudioUri:    t.AudioURI,
		IsLocal:     t.IsLocal,
		Isrc:        t.ISRC,
		DurationMs:  t.DurationMs,
		Popularity:  t.Popularity,
		PlayCount:   t.PlayCount,
	}
}
func tracksToProto(tracks []track.Track) []*trackv1.Track {
	res := make([]*trackv1.Track, len(tracks))

	for i, t := range tracks {
		res[i] = trackToProto(t)
	}
	return res
}
