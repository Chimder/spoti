package service

import (
	"context"

	"github.com/Chimder/spoti/internal/domain/user"
	meilisearchrepo "github.com/Chimder/spoti/internal/repository/meilisearch"
	"github.com/Chimder/spoti/internal/repository/postgres"
	rediscache "github.com/Chimder/spoti/internal/repository/redis"

	"github.com/google/uuid"
)

type UserService struct {
	repo  *postgres.Repository
	cache rediscache.Cache
	meili *meilisearchrepo.MeiliRepository
}

func NewUserService(repo *postgres.Repository, cache rediscache.Cache,
	meili *meilisearchrepo.MeiliRepository) *UserService {
	return &UserService{repo: repo, cache: cache}
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

func (us *UserService) GetUserByEmail(ctx context.Context, userEmail string) (user.User, error) {
	return us.repo.User.GetUserByEmail(ctx, userEmail)
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
