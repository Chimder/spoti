package albumrepo_test

import (
	"context"
	"testing"

	albumrepo "github.com/Chimder/spoti/internal/repository/postgres/album"
	"github.com/Chimder/spoti/internal/repository/postgres/testhelpers"
	userrepo "github.com/Chimder/spoti/internal/repository/postgres/user"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAlbumRepo_CreateAlbum(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	repo := albumrepo.NewAlbumRepo(db)

	t.Run("Ok create album", func(t *testing.T) {
		id := testhelpers.CreateTestAlbum(t, repo)
		assert.NotEqual(t, uuid.Nil, id)
	})

	t.Run("duplicate album name", func(t *testing.T) {
		ctx := context.Background()
		fake := testhelpers.FakeAlbum()

		_, err := repo.CreateAlbum(ctx, fake)
		require.NoError(t, err)

		_, err = repo.CreateAlbum(ctx, fake)
		assert.Error(t, err)
	})
}

func TestAlbumRepo_SaveAndRemoveAlbumsForUser(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	albumRepo := albumrepo.NewAlbumRepo(db)
	userRepo := userrepo.NewUserRepo(db)
	ctx := context.Background()

	userId := testhelpers.CreateTestUser(t, userRepo)
	albumId := testhelpers.CreateTestAlbum(t, albumRepo)

	t.Run("save album success", func(t *testing.T) {
		err := albumRepo.SaveAlbumsForCurrentUser(ctx, []string{albumId.String()}, userId.String())
		require.NoError(t, err)
	})

	t.Run("save album idempotent", func(t *testing.T) {
		err := albumRepo.SaveAlbumsForCurrentUser(ctx, []string{albumId.String()}, userId.String())
		require.NoError(t, err)

		err = albumRepo.SaveAlbumsForCurrentUser(ctx, []string{albumId.String()}, userId.String())
		require.NoError(t, err)
	})

	t.Run("save empty list no error", func(t *testing.T) {
		err := albumRepo.SaveAlbumsForCurrentUser(ctx, []string{}, userId.String())
		require.NoError(t, err)
	})

	t.Run("remove album", func(t *testing.T) {
		err := albumRepo.SaveAlbumsForCurrentUser(ctx, []string{albumId.String()}, userId.String())
		require.NoError(t, err)

		err = albumRepo.RemoveAlbumsFromCurrentUser(ctx, []string{albumId.String()}, userId.String())
		require.NoError(t, err)
	})

	t.Run("remove non-saved album", func(t *testing.T) {
		err := albumRepo.RemoveAlbumsFromCurrentUser(ctx, []string{uuid.New().String()}, userId.String())
		require.NoError(t, err)
	})

	t.Run("remove empty list", func(t *testing.T) {
		err := albumRepo.RemoveAlbumsFromCurrentUser(ctx, []string{}, userId.String())
		require.NoError(t, err)
	})
}

func TestAlbumRepo_GetUserSavedAlbums(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	albumRepo := albumrepo.NewAlbumRepo(db)
	userRepo := userrepo.NewUserRepo(db)
	ctx := context.Background()

	userId := testhelpers.CreateTestUser(t, userRepo)

	t.Run("empty saved albums", func(t *testing.T) {
		albums, err := albumRepo.GetUserSavedAlbums(ctx, userId.String())

		require.NoError(t, err)
		assert.Empty(t, albums)
	})

	t.Run("returns saved albums", func(t *testing.T) {
		id1 := testhelpers.CreateTestAlbum(t, albumRepo)
		id2 := testhelpers.CreateTestAlbum(t, albumRepo)

		err := albumRepo.SaveAlbumsForCurrentUser(ctx, []string{id1.String(), id2.String()}, userId.String())
		require.NoError(t, err)

		albums, err := albumRepo.GetUserSavedAlbums(ctx, userId.String())

		require.NoError(t, err)
		assert.Len(t, albums, 2)
	})

	t.Run("Bad user id", func(t *testing.T) {
		_, err := albumRepo.GetUserSavedAlbums(ctx, "Bad")
		assert.Error(t, err)
	})
}

func TestAlbumRepo_CheckUsersSavedAlbums(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	albumRepo := albumrepo.NewAlbumRepo(db)
	userRepo := userrepo.NewUserRepo(db)
	ctx := context.Background()

	userId := testhelpers.CreateTestUser(t, userRepo)
	savedId := testhelpers.CreateTestAlbum(t, albumRepo)
	unsavedId := testhelpers.CreateTestAlbum(t, albumRepo)

	err := albumRepo.SaveAlbumsForCurrentUser(ctx, []string{savedId.String()}, userId.String())
	require.NoError(t, err)

	t.Run("correct saved", func(t *testing.T) {
		result, err := albumRepo.CheckUsersSavedAlbums(ctx,
			[]string{savedId.String(), unsavedId.String()},
			userId.String(),
		)

		require.NoError(t, err)
		require.Len(t, result, 2)
		assert.True(t, result[0])
		assert.False(t, result[1])
	})

	t.Run("all unsaved", func(t *testing.T) {
		result, err := albumRepo.CheckUsersSavedAlbums(ctx,
			[]string{unsavedId.String()},
			userId.String(),
		)

		require.NoError(t, err)
		assert.Equal(t, []bool{false}, result)
	})
}

func TestAlbumRepo_GetNewReleases(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	albumRepo := albumrepo.NewAlbumRepo(db)
	ctx := context.Background()

	t.Run("returns results", func(t *testing.T) {
		for range 5 {
			testhelpers.CreateTestAlbum(t, albumRepo)
		}

		result, err := albumRepo.GetNewReleases(ctx, 3)

		require.NoError(t, err)
		assert.LessOrEqual(t, len(result), 3)
	})

	t.Run("limit zero", func(t *testing.T) {
		result, err := albumRepo.GetNewReleases(ctx, 0)

		require.NoError(t, err)
		assert.Empty(t, result)
	})
}

// func TestAlbumRepo_AlbumArtistRelation(t *testing.T) {
// 	db := testhelpers.SetupContainerDB()
// 	albumRepo := albumrepo.NewAlbumRepo(db)
// 	artistRepo := artistrepo.NewArtistRepo(db)
// 	ctx := context.Background()

// 	albumId := testhelpers.CreateTestAlbum(t, albumRepo)
// 	artistId := testhelpers.CreateTestArtist(t, artistRepo)

// 	_, err := db.Exec(ctx,
// 		`INSERT INTO album_artists (album_id, artist_id) VALUES ($1, $2)`,
// 		albumId, artistId,
// 	)
// 	require.NoError(t, err)

// t.Run("artist albums not empty", func(t *testing.T) {
// 	_, err := artistRepo.GetArtistAlbums(ctx, artistId.String())
// 	_ = err
// })
// }
