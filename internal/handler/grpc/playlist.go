package grpc

import (
	"context"

	"github.com/Chimder/spoti/internal/domain/playlist"
	playlistv1 "github.com/Chimder/spoti/internal/gen/playlist/v1"
	"github.com/Chimder/spoti/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type PlaylistHandler struct {
	playlistv1.UnimplementedPlaylistServiceServer
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
	// playlistID := req.PlaylistId
	// limit := req.Limit
	// offset := req.Offset

	// if limit == 0 {
	// 	limit = 50
	// }

	// if offset == 0 {
	// 	offset = 0
	// }

	// playlist, err := h.srv.GetPlaylistById(ctx, playlistID, int(limit), int(offset))
	// if err != nil {
	// 	return nil, status.Error(codes.NotFound, "not found playlist by id")
	// }

	// return &playlistv1.GetPlaylistByIdResponse{Playlist: playlist, }, status.Error(codes.Unimplemented, "method GetPlaylistById not implemented")
	return nil, status.Error(codes.Unimplemented, "method GetPlaylistById not implemented")
}
func (h *PlaylistHandler) AddToPlaylist(ctx context.Context, req *playlistv1.AddToPlaylistRequest) (*emptypb.Empty, error) {
	err := h.srv.AddToPlaylist(ctx, req.PlaylistId, req.TrackId)
	if err != nil {
		return nil, status.Error(codes.Internal, "err to add playlist")
	}
	return nil, nil
}

func (h *PlaylistHandler) UpdatePlaylist(ctx context.Context, req *playlistv1.UpdatePlaylistRequest) (*playlistv1.UpdatePlaylistResponse, error) {
	if err := h.srv.UpdatePlaylist(ctx, req.PlaylistId, playlist.UpdatePlaylistReq{Name: req.Name,
		Description: req.Description, Public: req.IsPublic}); err != nil {

		return nil, status.Error(codes.NotFound, "err update playlist")
	}

	return &playlistv1.UpdatePlaylistResponse{}, nil
}

func (h *PlaylistHandler) DeleteFromPlaylist(ctx context.Context, req *playlistv1.DeleteFromPlaylistRequest) (*playlistv1.DeleteFromPlaylistResponse, error) {
	if err := h.srv.DeleteFromPlaylist(ctx, req.PlaylistId, req.TrackId); err != nil {
		return nil, status.Error(codes.NotFound, "err delete playlist")
	}

	return &playlistv1.DeleteFromPlaylistResponse{}, nil
}

func (h *PlaylistHandler) GetAllUserPlaylists(ctx context.Context, req *playlistv1.GetAllUserPlaylistsRequest) (*playlistv1.GetAllUserPlaylistsResponse, error) {
	data, err := h.srv.GetAllUserPlaylists(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.NotFound, "err get all playlists")
	}
	return &playlistv1.GetAllUserPlaylistsResponse{Playlists: playlistsToProto(data)}, status.Error(codes.Unimplemented, "method GetAllUserPlaylists not implemented")
}

func playlistToProto(p playlist.Playlist) *playlistv1.Playlist {
	return &playlistv1.Playlist{
		Id:          p.Id.String(),
		OwnerId:     p.Owner.String(),
		Name:        p.Name,
		Description: p.Description,
		DiscNumber:  int32(p.DiscNumber),
		Image:       p.Img,
		IsPublic:    p.Public,
		Total:       uint32(p.Total),
		CreatedAt:   timestamppb.New(p.CreatedAt),
	}
}
func playlistsToProto(playlists []playlist.Playlist) []*playlistv1.Playlist {
	res := make([]*playlistv1.Playlist, len(playlists))
	for i, p := range playlists {
		res[i] = playlistToProto(p)
	}
	return res
}
