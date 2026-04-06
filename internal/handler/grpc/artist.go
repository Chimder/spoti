package grpc

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Chimder/spoti/internal/domain/artist"
	artistv1 "github.com/Chimder/spoti/internal/gen/artist/v1"
	"github.com/Chimder/spoti/internal/service"
)

type ArtistHandler struct {
	artistv1.UnimplementedArtistServiceServer

	srv *service.ArtistService
}

func NewArtistHandler(svc *service.ArtistService) *ArtistHandler {
	return &ArtistHandler{srv: svc}
}

func domainToProto(a artist.Artist) *artistv1.Artist {
	return &artistv1.Artist{
		Id:         a.Id.String(),
		Url:        a.Url,
		Followers:  a.Followers,
		Genres:     a.Genres,
		Image:      a.Image,
		Name:       a.Name,
		Popularity: uint32(a.Popularity),
		Uri:        a.URI,
		CreatedAt:  timestamppb.New(a.CreatedAt),
	}
}

func protoCreateToDomain(req *artistv1.CreateArtistRequest) artist.CreateArtistReq {
	return artist.CreateArtistReq{
		Url:        req.Url,
		Uri:        req.Uri,
		ArtistName: req.ArtistName,
		Image:      req.Image,
		Followers:  int(req.Followers),
		Popularity: int(req.Popularity),
		Genres:     req.Genres,
	}
}

func (h *ArtistHandler) GetArtist(ctx context.Context, req *artistv1.GetArtistRequest) (*artistv1.GetArtistResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	artist, err := h.srv.GetArtist(ctx, req.ArtistId)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &artistv1.GetArtistResponse{Artist: &artistv1.Artist{Id: artist.Id.String()}}, nil
}

func (h *ArtistHandler) GetArtistsByIds(ctx context.Context, req *artistv1.GetArtistsByIdsRequest) (*artistv1.GetArtistsByIdsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	artists, err := h.srv.GetArtistsByIDs(ctx, req.ArtistsIds)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	protoArtists := make([]*artistv1.Artist, len(artists))
	for i, a := range artists {
		protoArtists[i] = domainToProto(a)
	}

	return &artistv1.GetArtistsByIdsResponse{Artists: []*artistv1.Artist{}}, nil
}

func (h *ArtistHandler) GetArtistAlbums(ctx context.Context, req *artistv1.GetArtistAlbumsRequest) (*artistv1.GetArtistAlbumsResponse, error) {
	_, err := h.srv.GetArtistAlbums(ctx, req.String())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	return &artistv1.GetArtistAlbumsResponse{}, nil
}

func (h *ArtistHandler) CreateArtist(ctx context.Context, req *artistv1.CreateArtistRequest) (*artistv1.CreateArtistResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := h.srv.CreateArtist(ctx, protoCreateToDomain(req))
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &artistv1.CreateArtistResponse{Id: id.String()}, nil
}
