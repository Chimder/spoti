package playlistrepo_test

import (
	"context"
	albumrepo "spoti/internal/repository/postgres/album"
	playlistrepo "spoti/internal/repository/postgres/playlist"
	recordingrepo "spoti/internal/repository/postgres/recording"
	"spoti/internal/repository/postgres/testhelpers"
	trackrepo "spoti/internal/repository/postgres/track"
	userrepo "spoti/internal/repository/postgres/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPlaylistRepo_CreatePlaylist(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	userRepo := userrepo.NewUserRepo(db)
	playlistRepo := playlistrepo.NewPlaylistRepo(db)

	t.Run("Ok create playlist", func(t *testing.T) {
		ownerId := testhelpers.CreateTestUser(t, userRepo)

		id := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)
		assert.NotEqual(t, uuid.Nil, id)
	})

	t.Run("err owner id", func(t *testing.T) {
		ctx := context.Background()
		fake := testhelpers.GetFakePlaylist(uuid.New())

		_, err := playlistRepo.CreatePlaylist(ctx, fake)
		assert.Error(t, err)
	})
}

func TestPlaylistRepo_GetAllUserPlaylists(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	userRepo := userrepo.NewUserRepo(db)
	playlistRepo := playlistrepo.NewPlaylistRepo(db)
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
	db := testhelpers.SetupContainerDB()
	userRepo := userrepo.NewUserRepo(db)
	playlistRepo := playlistrepo.NewPlaylistRepo(db)
	albumRepo := albumrepo.NewAlbumRepo(db)
	recRepo := recordingrepo.NewRecordingRepo(db)
	trackRepo := trackrepo.NewTrackRepo(db)
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
		err := db.QueryRow(ctx,
			`SELECT track_position FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`,
			newPlaylistId, track3,
		).Scan(&pos)
		require.NoError(t, err)
		assert.Equal(t, 1, pos)
	})
}

func TestPlaylistRepo_UpdatePlaylist(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	userRepo := userrepo.NewUserRepo(db)
	playlistRepo := playlistrepo.NewPlaylistRepo(db)
	ctx := context.Background()

	ownerId := testhelpers.CreateTestUser(t, userRepo)
	playlistId := testhelpers.CreateTestPlaylist(t, playlistRepo, ownerId)

	t.Run("update name only", func(t *testing.T) {
		newName := "Updated Name"
		err := playlistRepo.UpdatePlaylist(ctx, playlistId.String(), playlistrepo.UpdatePlaylistReq{
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

		err := playlistRepo.UpdatePlaylist(ctx, playlistId.String(), playlistrepo.UpdatePlaylistReq{
			Name:        &name,
			Description: &desc,
			Public:      &public,
		})
		require.NoError(t, err)
	})

	t.Run("update with nil", func(t *testing.T) {
		err := playlistRepo.UpdatePlaylist(ctx, playlistId.String(), playlistrepo.UpdatePlaylistReq{})
		require.NoError(t, err)
	})

	t.Run("update non-exis playlist", func(t *testing.T) {
		name := "bad"
		err := playlistRepo.UpdatePlaylist(ctx, uuid.New().String(), playlistrepo.UpdatePlaylistReq{
			Name: &name,
		})
		require.NoError(t, err)
	})
}
