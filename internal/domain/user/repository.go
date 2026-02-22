package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUserReq) (uuid.UUID, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (User, error)
	FollowUserToPlaylist(ctx context.Context, userId, playlistId string) error
	UnfollowUserFromPlaylist(ctx context.Context, userId, playlistId string) error
	FollowUserToArtist(ctx context.Context, userId, artistId string) error
	UnfollowUserFromArtist(ctx context.Context, userId, artistId string) error
}
