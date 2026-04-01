package albumrepo_test

import (
	"context"
	"os"
	"testing"

	albumrepo "github.com/Chimder/spoti/internal/repository/postgres/album"
	artistrepo "github.com/Chimder/spoti/internal/repository/postgres/artist"
	recordingrepo "github.com/Chimder/spoti/internal/repository/postgres/recording"
	"github.com/Chimder/spoti/internal/repository/postgres/testhelpers"
	trackrepo "github.com/Chimder/spoti/internal/repository/postgres/track"
	userrepo "github.com/Chimder/spoti/internal/repository/postgres/user"
	"github.com/google/uuid"
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
func TestAlbumRepo_CreateAlbum(t *testing.T) {
	repo := albumrepo.NewAlbumRepo(testDB)

	t.Run("Ok create album", func(t *testing.T) {
		defer testhelpers.TruncateAll(t, testDB)
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
func TestAlbumRepo_GetAlbum(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	albumRepo := albumrepo.NewAlbumRepo(testDB)

	ctx := context.Background()
	t.Run("album", func(t *testing.T) {
		albumID := testhelpers.CreateTestAlbum(t, albumRepo)

		album, err := albumRepo.GetAlbum(ctx, albumID.String())
		require.NoError(t, err)
		assert.Equal(t, albumID, album.ID)
	})

	t.Run("album bad uuid", func(t *testing.T) {
		_, err := albumRepo.GetAlbum(ctx, uuid.New().String())
		assert.Error(t, err)
	})
}

func TestAlbumRepo_GetAlbumWithTracks(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)

	albumRepo := albumrepo.NewAlbumRepo(testDB)
	artistRepo := artistrepo.NewArtistRepo(testDB)
	recordingRepo := recordingrepo.NewRecordingRepo(testDB)
	trackRepo := trackrepo.NewTrackRepo(testDB)

	ctx := context.Background()

	t.Run("album without tracks", func(t *testing.T) {
		albumID := testhelpers.CreateTestAlbum(t, albumRepo)

		resp, err := albumRepo.GetAlbumWithTracks(ctx, albumID.String())
		require.NoError(t, err)

		assert.Equal(t, albumID, resp.ID)
		assert.Empty(t, resp.Tracks.Items)
	})

	t.Run("album with track", func(t *testing.T) {
		albumID := testhelpers.CreateTestAlbum(t, albumRepo)

		recordingID := testhelpers.CreateTestRecording(t, recordingRepo)
		trackID := testhelpers.CreateTestTrack(t, trackRepo, albumID, recordingID, 1, 1)

		resp, err := albumRepo.GetAlbumWithTracks(ctx, albumID.String())
		require.NoError(t, err)

		require.Len(t, resp.Tracks.Items, 1)

		track := resp.Tracks.Items[0]
		assert.Equal(t, trackID, track.ID)
		assert.Equal(t, 1, track.TrackNumber)
		assert.Equal(t, 1, track.DiscNumber)
		assert.Empty(t, track.Artists)
	})

	t.Run("album with artist track", func(t *testing.T) {
		albumID := testhelpers.CreateTestAlbum(t, albumRepo)

		recordingID := testhelpers.CreateTestRecording(t, recordingRepo)
		trackID := testhelpers.CreateTestTrack(t, trackRepo, albumID, recordingID, 1, 1)
		artistID := testhelpers.CreateTestArtist(t, artistRepo)

		err := trackRepo.AddArtistToTrack(ctx, trackID, artistID)
		require.NoError(t, err)

		resp, err := albumRepo.GetAlbumWithTracks(ctx, albumID.String())
		require.NoError(t, err)

		require.Len(t, resp.Tracks.Items, 1)

		track := resp.Tracks.Items[0]
		require.Len(t, track.Artists, 1)
		assert.Equal(t, artistID, track.Artists[0].ID)
	})

	t.Run("album bad uuid", func(t *testing.T) {
		_, err := albumRepo.GetAlbumWithTracks(ctx, uuid.New().String())
		assert.Error(t, err)
	})
}
func TestAlbumRepo_SaveAndRemoveAlbumsForUser(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	albumRepo := albumrepo.NewAlbumRepo(testDB)
	userRepo := userrepo.NewUserRepo(testDB)
	ctx := context.Background()

	userId := testhelpers.CreateTestUser(t, userRepo)
	albumId := testhelpers.CreateTestAlbum(t, albumRepo)

	t.Run("save album success", func(t *testing.T) {
		err := albumRepo.SaveAlbumsForCurrentUser(ctx, []string{albumId.String()}, userId.String())
		require.NoError(t, err)
	})

	t.Run("save same album ", func(t *testing.T) {
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
	defer testhelpers.TruncateAll(t, testDB)
	albumRepo := albumrepo.NewAlbumRepo(testDB)
	userRepo := userrepo.NewUserRepo(testDB)
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
	defer testhelpers.TruncateAll(t, testDB)
	albumRepo := albumrepo.NewAlbumRepo(testDB)
	userRepo := userrepo.NewUserRepo(testDB)
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
	defer testhelpers.TruncateAll(t, testDB)
	albumRepo := albumrepo.NewAlbumRepo(testDB)
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
