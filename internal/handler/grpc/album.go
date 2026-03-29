package grpc

import (
	"context"
	"fmt"

	"github.com/Chimder/spoti/internal/domain/album"
	albumv1 "github.com/Chimder/spoti/internal/gen/album/v1"
	"github.com/Chimder/spoti/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type AlbumHandler struct {
	albumv1.UnimplementedAlbumServiceServer
	srv *service.AlbumService
}

func NewAlbumHandler(srv *service.AlbumService) *AlbumHandler {
	return &AlbumHandler{srv: srv}
}

func (h *AlbumHandler) CreateAlbum(ctx context.Context, req *albumv1.CreateAlbumRequest) (*emptypb.Empty, error) {
	if err := h.srv.CreateAlbum(ctx, protoToDomainCreate(req)); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (h *AlbumHandler) GetUserSavedAlbums(ctx context.Context, req *albumv1.GetUserSavedAlbumsRequest) (*albumv1.GetUserSavedAlbumsResponse, error) {
	albums, err := h.srv.GetUserSavedAlbums(ctx, req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "err to get saved albums")
	}
	return &albumv1.GetUserSavedAlbumsResponse{Albums: albumsToProto(albums)}, nil
}

func (h *AlbumHandler) SaveAlbumsForCurrentUser(ctx context.Context, req *albumv1.SaveAlbumsForCurrentUserRequest) (*emptypb.Empty, error) {
	if len(req.AlbumIds) == 0 {
		return nil, fmt.Errorf("invalid album_ids")
	}

	if err := h.srv.SaveAlbumsForCurrentUser(ctx, req.AlbumIds, req.UserId); err != nil {
		return nil, status.Error(codes.Internal, "err to saved albums")
	}

	return &emptypb.Empty{}, nil
}

func (h *AlbumHandler) RemoveAlbumsFromCurrentUser(ctx context.Context, req *albumv1.RemoveAlbumsFromCurrentUserRequest) (*emptypb.Empty, error) {
	if len(req.AlbumIds) == 0 {
		return nil, fmt.Errorf("invalid album_ids")
	}

	if err := h.srv.RemoveAlbumsFromCurrentUser(ctx, req.AlbumIds, req.UserId); err != nil {
		return nil, status.Error(codes.Internal, "err to remove albums")
	}

	return &emptypb.Empty{}, nil
}

func (h *AlbumHandler) CheckUsersSavedAlbums(ctx context.Context, req *albumv1.CheckUsersSavedAlbumsRequest) (*albumv1.CheckUsersSavedAlbumsResponse, error) {
	if len(req.AlbumIds) == 0 {
		return nil, fmt.Errorf("invalid album_ids")
	}

	result, err := h.srv.CheckUsersSavedAlbums(ctx, req.AlbumIds, req.UserId)
	if err != nil {
		return nil, status.Error(codes.Internal, "err to check saved albums")
	}

	return &albumv1.CheckUsersSavedAlbumsResponse{Result: result}, nil
}

func (h *AlbumHandler) GetNewReleases(ctx context.Context, req *albumv1.GetNewReleasesRequest) (*albumv1.GetNewReleasesResponse, error) {
	limit := int(req.Limit)
	if limit == 0 {
		limit = 10
	}

	albums, err := h.srv.GetNewReleases(ctx, limit)
	if err != nil {
		return nil, status.Error(codes.NotFound, "err to get new releases")
	}

	return &albumv1.GetNewReleasesResponse{Albums: albumsToProto(albums)}, nil
}

func protoToDomainCreate(req *albumv1.CreateAlbumRequest) album.CreateAlbumReq {
	return album.CreateAlbumReq{
		AlbumType:   req.AlbumType,
		TotalTracks: int(req.TotalTracks),
		Image:       req.Image,
		AlbumName:   req.AlbumName,
		Uri:         req.Uri,
		Copyrights:  req.Copyrights,
		AlbumLabel:  req.AlbumLabel,
		Popularity:  int(req.Popularity),
		ReleaseDate: req.ReleaseDate.AsTime(),
	}
}

func albumToProto(a album.Album) *albumv1.Album {
	return &albumv1.Album{
		Id:          a.ID.String(),
		AlbumType:   a.AlbumType,
		TotalTracks: int32(a.TotalTracks),
		Image:       a.Image,
		Name:        a.Name,
		Uri:         a.URI,
		Copyrights:  a.Copyrights,
		Label:       a.Label,
		Popularity:  int32(a.Popularity),
		ReleaseDate: timestamppb.New(a.ReleaseDate),
		CreatedAt:   timestamppb.New(a.CreatedAt),
	}
}
func albumsToProto(a []album.Album) []*albumv1.Album {
	albums := make([]*albumv1.Album, len(a))
	for i, v := range a {
		albums[i] = albumToProto(v)
	}
	return albums
}
