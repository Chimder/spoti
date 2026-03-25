package grpc

import (
	"context"

	"github.com/Chimder/spoti/internal/domain/user"
	userv1 "github.com/Chimder/spoti/internal/gen/user/v1"
	"github.com/Chimder/spoti/internal/service"
	"github.com/google/uuid"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type UserHandler struct {
	srv *service.UserService
}

func NewUserHandler(s *service.UserService) *UserHandler {
	return &UserHandler{srv: s}
}

func (h *UserHandler) CreateUser(ctx context.Context, req *userv1.CreateUserRequest) (*userv1.CreateUserResponse, error) {
	if err := req.ValidateAll(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	id, err := h.srv.CreateUser(ctx, user.CreateUserReq{Name: req.Name, Email: req.Email,
		Image: req.Image, Followers: uint32(req.Followers), PremiumStatus: req.PremiumStatus})
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	return &userv1.CreateUserResponse{Id: id.String()}, nil
}

func (h *UserHandler) GetUserByID(ctx context.Context, req *userv1.GetUserByIDRequest) (*userv1.GetUserByIDResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	id, err := uuid.Parse(req.Id)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid uuid")
	}

	u, err := h.srv.GetUserById(ctx, id)
	if err != nil {
		return nil, err
	}

	return &userv1.GetUserByIDResponse{Id: u.Id.String(), Name: u.Name, Email: u.Email, Image: u.Image,
		Followers: u.Followers, PremiumStatus: u.PremiumStatus}, nil
}

func (h *UserHandler) FollowPlaylist(ctx context.Context, req *userv1.FollowPlaylistRequest) (*emptypb.Empty, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userId")
	}

	playlistID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid playlistId")
	}

	err = h.srv.FollowUserToPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to follow playlist")
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) UnfollowPlaylist(ctx context.Context, req *userv1.UnfollowPlaylistRequest) (*emptypb.Empty, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userId")
	}
	playlistID, err := uuid.Parse(req.PlaylistId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid playlistId")
	}

	err = h.srv.UnfollowUserFromPlaylist(ctx, userID, playlistID)
	if err != nil {
		return nil, status.Error(codes.NotFound, "failed to unfollow playlist")
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) FollowArtist(ctx context.Context, req *userv1.FollowArtistRequest) (*emptypb.Empty, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userId")
	}
	artistID, err := uuid.Parse(req.ArtistId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid artistId")
	}

	if err := h.srv.FollowUserToArtist(ctx, userID, artistID); err != nil {
		return nil, status.Error(codes.NotFound, "failed to follow artist")
	}

	return &emptypb.Empty{}, nil
}

func (h *UserHandler) UnfollowArtist(ctx context.Context, req *userv1.UnfollowArtistRequest) (*emptypb.Empty, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid userId")
	}
	artistID, err := uuid.Parse(req.ArtistId)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid artistId")
	}

	if err := h.srv.UnfollowUserFromArtist(ctx, userID, artistID); err != nil {
		return nil, status.Error(codes.NotFound, "failed to unfollow artist")
	}

	return &emptypb.Empty{}, nil
}
