package userrepo_test

import (
	"context"
	"os"
	"testing"

	"github.com/Chimder/spoti/internal/handler/http/middleware"
	artistrepo "github.com/Chimder/spoti/internal/repository/postgres/artist"
	playlistrepo "github.com/Chimder/spoti/internal/repository/postgres/playlist"
	"github.com/Chimder/spoti/internal/repository/postgres/testhelpers"
	userrepo "github.com/Chimder/spoti/internal/repository/postgres/user"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *pgxpool.Pool
var cleanup func()

func TestMain(m *testing.M) {
	testDB, cleanup = testhelpers.SetupContainerDB()
	code := m.Run()
	cleanup()
	os.Exit(code)
}
func TestUserRepo_CreateUser(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	repo := userrepo.NewUserRepo(testDB)
	ctx := context.Background()

	t.Run("success create user", func(t *testing.T) {
		userId := testhelpers.CreateTestUser(t, repo)

		assert.NotEqual(t, uuid.Nil, userId)
	})

	t.Run("error duplicate email", func(t *testing.T) {
		fakeUser := testhelpers.FakeUser()
		hashPass, err := middleware.GeneratePass(fakeUser.Password)
		require.NotEmpty(t, hashPass)
		require.NoError(t, err)

		_, err = repo.CreateUser(ctx, fakeUser, hashPass)
		require.NoError(t, err)

		_, err = repo.CreateUser(ctx, fakeUser, hashPass)
		assert.Error(t, err)
	})
}

func TestUserRepo_GetUserById(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	repo := userrepo.NewUserRepo(testDB)
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
	defer testhelpers.TruncateAll(t, testDB)
	userRepo := userrepo.NewUserRepo(testDB)
	playlistRepo := playlistrepo.NewPlaylistRepo(testDB)

	ctx := context.Background()

	userId := testhelpers.CreateTestUser(t, userRepo)

	playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, userId)

	t.Run("follow success", func(t *testing.T) {
		err := userRepo.FollowUserToPlaylist(ctx, userId, playlistId)
		require.NoError(t, err)
	})

	t.Run("follow repeat", func(t *testing.T) {
		err := userRepo.FollowUserToPlaylist(ctx, userId, playlistId)
		require.NoError(t, err)

		err = userRepo.FollowUserToPlaylist(ctx, userId, playlistId)
		require.NoError(t, err)
	})

	t.Run("unfollow success", func(t *testing.T) {
		err := userRepo.UnfollowUserFromPlaylist(ctx, userId, playlistId)
		require.NoError(t, err)
	})

	t.Run("unfollow repeat", func(t *testing.T) {
		err := userRepo.UnfollowUserFromPlaylist(ctx, userId, playlistId)
		require.NoError(t, err)
	})

	t.Run("Bad user id", func(t *testing.T) {
		badId := uuid.New()
		err := userRepo.FollowUserToPlaylist(ctx, badId, playlistId)
		assert.Error(t, err)
	})
}

func TestUserRepo_FollowUnfollowArtist(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	userRepo := userrepo.NewUserRepo(testDB)
	artistRepo := artistrepo.NewArtistRepo(testDB)
	ctx := context.Background()

	userId := testhelpers.CreateTestUser(t, userRepo)

	artistId := testhelpers.CreateTestArtist(t, artistRepo)

	t.Run("follow success", func(t *testing.T) {
		err := userRepo.FollowUserToArtist(ctx, userId, artistId)
		require.NoError(t, err)
	})

	t.Run("follow repeat", func(t *testing.T) {
		err := userRepo.FollowUserToArtist(ctx, userId, artistId)
		require.NoError(t, err)
	})

	t.Run("unfollow success", func(t *testing.T) {
		err := userRepo.UnfollowUserFromArtist(ctx, userId, artistId)
		require.NoError(t, err)
	})

	t.Run("invalid artist id", func(t *testing.T) {
		badId := uuid.New()
		err := userRepo.FollowUserToArtist(ctx, userId, badId)
		assert.Error(t, err)
	})
}
