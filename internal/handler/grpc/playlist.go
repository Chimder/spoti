package grpc

import (
	"context"

	"github.com/Chimder/spoti/internal/domain/playlist"
	playlistv1 "github.com/Chimder/spoti/internal/gen/playlist/v1"
	"github.com/Chimder/spoti/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type PlaylistHandler struct {
	srv *service.PlaylistService
}

func NewPlaylistHandler(srv *service.PlaylistService) *PlaylistHandler {
	return &PlaylistHandler{srv: srv}
}

func (h *PlaylistHandler) CreatePlaylist(ctx context.Context, req *playlistv1.CreatePlaylistRequest) (*playlistv1.CreatePlaylistResponse, error) {
	ownerId, err := uuid.Parse(req.OwnerId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid ownerId")
	}

	id, err := h.srv.CreatePlaylist(ctx, playlist.CreatePlaylistReq{OwnerId: ownerId, PlaylistName: req.PlaylistName,
		Description: req.Description, Image: req.Image, IsPublic: req.IsPublic})
	if err != nil {
		return nil, status.Error(codes.Internal, "err create playlist")
	}
	return &playlistv1.CreatePlaylistResponse{Id: id.String()}, nil
}

func (h *PlaylistHandler) GetPlaylistById(ctx context.Context, req *playlistv1.GetPlaylistByIdRequest) (*playlistv1.GetPlaylistByIdResponse, error) {
	playlistID := req.PlaylistId
	limit := req.Limit
	offset := req.Offset

	if limit == 0 {
		limit = 50
	}

	// if offset == 0 {
	// 	offset = 0
	// }

	playlist, err := h.srv.GetPlaylistById(ctx, playlistID, int(limit), int(offset))
	if err != nil {
		return nil, status.Error(codes.NotFound, "not found playlist by id")
	}

	// return &playlistv1.GetPlaylistByIdResponse{Playlist: playlist, }, status.Error(codes.Unimplemented, "method GetPlaylistById not implemented")
}
func (h *PlaylistHandler) AddToPlaylist(ctx context.Context, req *playlistv1.AddToPlaylistRequest) (*playlistv1.AddToPlaylistResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method AddToPlaylist not implemented")
}
func (h *PlaylistHandler) UpdatePlaylist(ctx context.Context, req *playlistv1.UpdatePlaylistRequest) (*playlistv1.UpdatePlaylistResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method UpdatePlaylist not implemented")
}
func (h *PlaylistHandler) DeleteFromPlaylist(ctx context.Context, req *playlistv1.DeleteFromPlaylistRequest) (*playlistv1.DeleteFromPlaylistResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method DeleteFromPlaylist not implemented")
}
func (h *PlaylistHandler) GetAllUserPlaylists(ctx context.Context, req *playlistv1.GetAllUserPlaylistsRequest) (*playlistv1.GetAllUserPlaylistsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "method GetAllUserPlaylists not implemented")
}
