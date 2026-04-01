package playlistrepo_test

import (
	"context"
	"os"
	"testing"

	"github.com/Chimder/spoti/internal/domain/playlist"
	albumrepo "github.com/Chimder/spoti/internal/repository/postgres/album"
	playlistrepo "github.com/Chimder/spoti/internal/repository/postgres/playlist"
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
func TestPlaylistRepo_CreatePlaylist(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	userRepo := userrepo.NewUserRepo(testDB)
	playlistRepo := playlistrepo.NewPlaylistRepo(testDB)

	t.Run("Ok create playlist", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)

		id := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)
		assert.NotEqual(t, uuid.Nil, id)
	})

	t.Run("err owner id", func(t *testing.T) {
		ctx := context.Background()
		fake := testhelpers.FakePlaylist(uuid.New())

		_, err := playlistRepo.CreatePlaylist(ctx, fake)
		assert.Error(t, err)
	})
}

func TestPlaylistRepo_GetPlaylistById(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)

	userRepo := userrepo.NewUserRepo(testDB)
	playlistRepo := playlistrepo.NewPlaylistRepo(testDB)
	recordingRepo := recordingrepo.NewRecordingRepo(testDB)
	trackRepo := trackrepo.NewTrackRepo(testDB)
	albumRepo := albumrepo.NewAlbumRepo(testDB)

	ctx := context.Background()

	t.Run("get playlist", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)
		playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

		resp, err := playlistRepo.GetPlaylistById(ctx, playlistId.String(), 10, 0)
		require.NoError(t, err)

		assert.Equal(t, playlistId, resp.Id)
		assert.Empty(t, resp.Tracks.Items)
	})

	t.Run("playlist with track", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)
		playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

		albumId := testhelpers.CreateTestAlbum(t, albumRepo)
		recordingId := testhelpers.CreateTestRecording(t, recordingRepo)
		trackId := testhelpers.CreateTestTrack(t, trackRepo, albumId, recordingId, 1, 1)

		err := playlistRepo.AddToPlaylist(ctx, playlistId.String(), trackId.String())
		require.NoError(t, err)

		resp, err := playlistRepo.GetPlaylistById(ctx, playlistId.String(), 10, 0)
		require.NoError(t, err)

		require.Len(t, resp.Tracks.Items, 1)

		item := resp.Tracks.Items[0]
		assert.Equal(t, trackId, item.Track.Id)
		assert.Equal(t, 1, item.Track.TrackNumber)
		assert.Equal(t, 1, item.Track.DiscNumber)
	})

	t.Run("playlist with many tracks", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)
		playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

		albumId := testhelpers.CreateTestAlbum(t, albumRepo)

		var trackIDs []uuid.UUID

		for i := range 3 {
			rec := testhelpers.CreateTestRecording(t, recordingRepo)
			tr := testhelpers.CreateTestTrack(t, trackRepo, albumId, rec, i+1, 1)

			err := playlistRepo.AddToPlaylist(ctx, playlistId.String(), tr.String())
			require.NoError(t, err)

			trackIDs = append(trackIDs, tr)
		}

		resp, err := playlistRepo.GetPlaylistById(ctx, playlistId.String(), 10, 0)
		require.NoError(t, err)

		require.Len(t, resp.Tracks.Items, 3)

		for i, item := range resp.Tracks.Items {
			assert.Equal(t, trackIDs[i], item.Track.Id)
		}
	})

	t.Run("limit and offset", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)
		playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

		albumId := testhelpers.CreateTestAlbum(t, albumRepo)

		for i := range 3 {
			rec := testhelpers.CreateTestRecording(t, recordingRepo)
			tr := testhelpers.CreateTestTrack(t, trackRepo, albumId, rec, i+1, 1)

			err := playlistRepo.AddToPlaylist(ctx, playlistId.String(), tr.String())
			require.NoError(t, err)
		}

		resp, err := playlistRepo.GetPlaylistById(ctx, playlistId.String(), 2, 1)
		require.NoError(t, err)

		assert.Len(t, resp.Tracks.Items, 2)
		assert.Equal(t, 2, resp.Tracks.Limit)
		assert.Equal(t, 1, resp.Tracks.Offset)
	})

	t.Run("correct limit", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)
		playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

		resp, err := playlistRepo.GetPlaylistById(ctx, playlistId.String(), 100, -10)
		require.NoError(t, err)

		assert.Equal(t, 20, resp.Tracks.Limit)
		assert.Equal(t, 0, resp.Tracks.Offset)
	})

	t.Run("playlist not found", func(t *testing.T) {
		_, err := playlistRepo.GetPlaylistById(ctx, uuid.New().String(), 10, 0)
		assert.Error(t, err)
	})
}
func TestPlaylistRepo_GetAllUserPlaylists(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	userRepo := userrepo.NewUserRepo(testDB)
	playlistRepo := playlistrepo.NewPlaylistRepo(testDB)
	ctx := context.Background()

	t.Run("all user playlists", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)

		testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)
		testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

		got, err := playlistRepo.GetAllUserPlaylists(ctx, ownerId.String())

		require.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("user with no playlists", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)

		got, err := playlistRepo.GetAllUserPlaylists(ctx, ownerId.String())

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("bad user id", func(t *testing.T) {
		_, err := playlistRepo.GetAllUserPlaylists(ctx, "badId")
		assert.Error(t, err)
	})
}

func TestPlaylistRepo_AddAndDeleteFromPlaylist(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	userRepo := userrepo.NewUserRepo(testDB)
	playlistRepo := playlistrepo.NewPlaylistRepo(testDB)
	albumRepo := albumrepo.NewAlbumRepo(testDB)
	recRepo := recordingrepo.NewRecordingRepo(testDB)
	trackRepo := trackrepo.NewTrackRepo(testDB)
	ctx := context.Background()

	ownerId := testhelpers.CreateTestUser(t, userRepo)
	playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)
	albumId := testhelpers.CreateTestAlbum(t, albumRepo)
	recId := testhelpers.CreateTestRecording(t, recRepo)
	trackId := testhelpers.CreateTestTrack(t, trackRepo, albumId, recId, 1, 1)

	t.Run("add track", func(t *testing.T) {
		err := playlistRepo.AddToPlaylist(ctx, playlistId.String(), trackId.String())
		require.NoError(t, err)
	})

	t.Run("add same track ", func(t *testing.T) {
		err := playlistRepo.AddToPlaylist(ctx, playlistId.String(), trackId.String())
		assert.Error(t, err)
	})

	t.Run("delete track", func(t *testing.T) {
		err := playlistRepo.DeleteFromPlaylist(ctx, playlistId.String(), trackId.String())
		require.NoError(t, err)
	})

	t.Run("delete non-ex track", func(t *testing.T) {
		err := playlistRepo.DeleteFromPlaylist(ctx, playlistId.String(), uuid.New().String())
		require.NoError(t, err)
	})

	t.Run("track position reorders", func(t *testing.T) {
		recId2 := testhelpers.CreateTestRecording(t, recRepo)
		recId3 := testhelpers.CreateTestRecording(t, recRepo)
		track2 := testhelpers.CreateTestTrack(t, trackRepo, albumId, recId2, 2, 1)
		track3 := testhelpers.CreateTestTrack(t, trackRepo, albumId, recId3, 3, 1)

		newPlaylistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

		require.NoError(t, playlistRepo.AddToPlaylist(ctx, newPlaylistId.String(), track2.String()))
		require.NoError(t, playlistRepo.AddToPlaylist(ctx, newPlaylistId.String(), track3.String()))
		require.NoError(t, playlistRepo.DeleteFromPlaylist(ctx, newPlaylistId.String(), track2.String()))

		var pos int
		err := testDB.QueryRow(ctx,
			`SELECT track_position FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`,
			newPlaylistId, track3,
		).Scan(&pos)
		require.NoError(t, err)
		assert.Equal(t, 1, pos)
	})
}

func TestPlaylistRepo_UpdatePlaylist(t *testing.T) {
	defer testhelpers.TruncateAll(t, testDB)
	userRepo := userrepo.NewUserRepo(testDB)
	playlistRepo := playlistrepo.NewPlaylistRepo(testDB)
	ctx := context.Background()

	ownerId := testhelpers.CreateTestUser(t, userRepo)
	playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

	t.Run("update name only", func(t *testing.T) {
		newName := "Updated Name"
		err := playlistRepo.UpdatePlaylist(ctx, playlistId.String(), playlist.UpdatePlaylistReq{
			Name: &newName,
		})
		require.NoError(t, err)

		playlists, err := playlistRepo.GetAllUserPlaylists(ctx, ownerId.String())
		require.NoError(t, err)
		require.Len(t, playlists, 1)
		assert.Equal(t, newName, playlists[0].Name)
	})

	t.Run("update all fields", func(t *testing.T) {
		name := "Update"
		desc := "New"
		public := true

		err := playlistRepo.UpdatePlaylist(ctx, playlistId.String(), playlist.UpdatePlaylistReq{
			Name:        &name,
			Description: &desc,
			Public:      &public,
		})
		require.NoError(t, err)
	})

	t.Run("update with nil", func(t *testing.T) {
		err := playlistRepo.UpdatePlaylist(ctx, playlistId.String(), playlist.UpdatePlaylistReq{})
		require.NoError(t, err)
	})

	t.Run("update non-exis playlist", func(t *testing.T) {
		name := "bad"
		err := playlistRepo.UpdatePlaylist(ctx, uuid.New().String(), playlist.UpdatePlaylistReq{
			Name: &name,
		})
		require.NoError(t, err)
	})
}
