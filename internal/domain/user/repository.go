package user

import (
	"context"

	"github.com/google/uuid"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user CreateUserReq) (uuid.UUID, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (User, error)
	FollowUserToPlaylist(ctx context.Context, userId, playlistId uuid.UUID) error
	UnfollowUserFromPlaylist(ctx context.Context, userId, playlistId uuid.UUID) error
	FollowUserToArtist(ctx context.Context, userId, artistId uuid.UUID) error
	UnfollowUserFromArtist(ctx context.Context, userId, artistId uuid.UUID) error
}
