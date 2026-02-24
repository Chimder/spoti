package service

import (
	"context"
	"spoti/internal/domain/user"
	"spoti/internal/repository/postgres"

	"github.com/google/uuid"
)

type UserService struct {
	repo *postgres.Repository
}

func NewUserService(repo *postgres.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}
func (us *UserService) CreateUser(ctx context.Context, user user.CreateUserReq) (uuid.UUID, error) {
	id, err := us.repo.User.CreateUser(ctx, user)
	if err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (us *UserService) GetUserById(ctx context.Context, userId uuid.UUID) (user.User, error) {
	return us.repo.User.GetUserById(ctx, userId)
}

func (us *UserService) FollowUserToPlaylist(ctx context.Context, userId, playlistId uuid.UUID) error {
	return us.repo.User.FollowUserToPlaylist(ctx, userId, playlistId)
}

func (us *UserService) UnfollowUserFromPlaylist(ctx context.Context, userId, playlistId uuid.UUID) error {
	return us.repo.User.UnfollowUserFromPlaylist(ctx, userId, playlistId)
}

func (us *UserService) FollowUserToArtist(ctx context.Context, userId, artistId uuid.UUID) error {
	return us.repo.User.FollowUserToArtist(ctx, userId, artistId)
}

func (us *UserService) UnfollowUserFromArtist(ctx context.Context, userId, artistId uuid.UUID) error {
	return us.repo.User.UnfollowUserFromArtist(ctx, userId, artistId)
}
