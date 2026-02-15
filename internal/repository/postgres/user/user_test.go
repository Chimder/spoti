package userrepo_test

import (
	"context"
	artistrepo "spoti/internal/repository/postgres/artist"
	playlistrepo "spoti/internal/repository/postgres/playlist"
	"spoti/internal/repository/postgres/testhelpers"
	userrepo "spoti/internal/repository/postgres/user"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepo_CreateUser(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	repo := userrepo.NewUserRepo(db)
	ctx := context.Background()

	t.Run("success create user", func(t *testing.T) {
		userId := testhelpers.CreateTestUser(t, repo)

		assert.NotEqual(t, uuid.Nil, userId)
	})

	t.Run("error duplicate email", func(t *testing.T) {
		fakeUser := testhelpers.GetFakeUser()

		_, err := repo.CreateUser(ctx, fakeUser)
		require.NoError(t, err)

		_, err = repo.CreateUser(ctx, fakeUser)
		assert.Error(t, err)
	})
}

func TestUserRepo_GetUserById(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	repo := userrepo.NewUserRepo(db)
	ctx := context.Background()

	t.Run("success existing user", func(t *testing.T) {
		userId := testhelpers.CreateTestUser(t, repo)

		user, err := repo.GetUserById(ctx, userId)

		require.NoError(t, err)
		assert.Equal(t, userId, user.Id)
	})

	t.Run("error not found", func(t *testing.T) {
		_, err := repo.GetUserById(ctx, uuid.New())
		assert.ErrorIs(t, err, pgx.ErrNoRows)
	})

	t.Run("error - nil uuid", func(t *testing.T) {
		_, err := repo.GetUserById(ctx, uuid.Nil)
		assert.Error(t, err)
	})
}

func TestUserRepo_FollowUnfollowPlaylist(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	userRepo := userrepo.NewUserRepo(db)
	playlistRepo := playlistrepo.NewPlaylistRepo(db)

	ctx := context.Background()

	userId := testhelpers.CreateTestUser(t, userRepo)
	playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, userId)

	t.Run("follow success", func(t *testing.T) {
		err := userRepo.FollowUserToPlaylist(ctx, userId.String(), playlistId.String())
		require.NoError(t, err)
	})

	t.Run("follow repeat", func(t *testing.T) {
		err := userRepo.FollowUserToPlaylist(ctx, userId.String(), playlistId.String())
		require.NoError(t, err)

		err = userRepo.FollowUserToPlaylist(ctx, userId.String(), playlistId.String())
		require.NoError(t, err)
	})

	t.Run("unfollow success", func(t *testing.T) {
		err := userRepo.UnfollowUserFromPlaylist(ctx, userId.String(), playlistId.String())
		require.NoError(t, err)
	})

	t.Run("unfollow repeat", func(t *testing.T) {
		err := userRepo.UnfollowUserFromPlaylist(ctx, userId.String(), playlistId.String())
		require.NoError(t, err)
	})

	t.Run("invalid user id", func(t *testing.T) {
		err := userRepo.FollowUserToPlaylist(ctx, "invalid", playlistId.String())
		assert.Error(t, err)
	})
}

func TestUserRepo_FollowUnfollowArtist(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	userRepo := userrepo.NewUserRepo(db)
	artistRepo := artistrepo.NewArtistRepo(db)
	ctx := context.Background()

	userId := testhelpers.CreateTestUser(t, userRepo)
	artistId := testhelpers.CreateTestArtist(t, artistRepo)

	t.Run("follow success", func(t *testing.T) {
		err := userRepo.FollowUserToArtist(ctx, userId.String(), artistId.String())
		require.NoError(t, err)
	})

	t.Run("follow repeat", func(t *testing.T) {
		err := userRepo.FollowUserToArtist(ctx, userId.String(), artistId.String())
		require.NoError(t, err)
	})

	t.Run("unfollow success", func(t *testing.T) {
		err := userRepo.UnfollowUserFromArtist(ctx, userId.String(), artistId.String())
		require.NoError(t, err)
	})

	t.Run("invalid artist id", func(t *testing.T) {
		err := userRepo.FollowUserToArtist(ctx, userId.String(), "")
		assert.Error(t, err)
	})
}
