package trackrepo_test

import (
	"context"
	albumrepo "github.com/Chimder/spoti/internal/repository/postgres/album"
	artistrepo "github.com/Chimder/spoti/internal/repository/postgres/artist"
	recordingrepo "github.com/Chimder/spoti/internal/repository/postgres/recording"
	"github.com/Chimder/spoti/internal/repository/postgres/testhelpers"
	trackrepo "github.com/Chimder/spoti/internal/repository/postgres/track"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setTrackDeps(t *testing.T, albumRepo *albumrepo.AlbumRepo, recRepo *recordingrepo.RecordingRepo) (uuid.UUID, uuid.UUID) {
	t.Helper()
	albumId := testhelpers.CreateTestAlbum(t, albumRepo)
	recId := testhelpers.CreateTestRecording(t, recRepo)
	return albumId, recId
}

func TestTrackRepo_CreateTrack(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	albumRepo := albumrepo.NewAlbumRepo(db)
	recRepo := recordingrepo.NewRecordingRepo(db)
	trackRepo := trackrepo.NewTrackRepo(db)

	t.Run("ok create track", func(t *testing.T) {
		albumId, recId := setTrackDeps(t, albumRepo, recRepo)

		id := testhelpers.CreateTestTrack(t, trackRepo, albumId, recId, 1, 1)

		assert.NotEqual(t, uuid.Nil, id)
	})

	t.Run("duplicate disc and track same disc", func(t *testing.T) {
		ctx := context.Background()
		albumId, recId := setTrackDeps(t, albumRepo, recRepo)

		_, err := trackRepo.CreateTrack(ctx, testhelpers.FakeTrack(albumId, recId, 1, 1))
		require.NoError(t, err)

		recId2, err := testhelpers.CreateRecording(recRepo)
		_, err = trackRepo.CreateTrack(ctx, testhelpers.FakeTrack(albumId, recId2, 1, 1))
		assert.Error(t, err)
	})

	t.Run("same track number on not same disc", func(t *testing.T) {
		ctx := context.Background()
		albumId, recId1 := setTrackDeps(t, albumRepo, recRepo)
		recId2 := testhelpers.CreateTestRecording(t, recRepo)

		_, err := trackRepo.CreateTrack(ctx, testhelpers.FakeTrack(albumId, recId1, 1, 1))
		require.NoError(t, err)

		_, err = trackRepo.CreateTrack(ctx, testhelpers.FakeTrack(albumId, recId2, 1, 2))
		require.NoError(t, err)
	})
}

func TestTrackRepo_GetTrackById(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	albumRepo := albumrepo.NewAlbumRepo(db)
	recRepo := recordingrepo.NewRecordingRepo(db)
	trackRepo := trackrepo.NewTrackRepo(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		albumId, recId := setTrackDeps(t, albumRepo, recRepo)
		id := testhelpers.CreateTestTrack(t, trackRepo, albumId, recId, 1, 1)

		got, err := trackRepo.GetTrackById(ctx, id)

		require.NoError(t, err)
		assert.Equal(t, id, got.Id)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := trackRepo.GetTrackById(ctx, uuid.New())
		assert.Error(t, err)
	})

	t.Run("bad uuid", func(t *testing.T) {
		_, err := trackRepo.GetTrackById(ctx, uuid.New())
		assert.Error(t, err)
	})
}

func TestTrackRepo_GetTracksByIds(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	albumRepo := albumrepo.NewAlbumRepo(db)
	recRepo := recordingrepo.NewRecordingRepo(db)
	trackRepo := trackrepo.NewTrackRepo(db)
	ctx := context.Background()

	t.Run("return all tracks", func(t *testing.T) {
		albumId, recId1 := setTrackDeps(t, albumRepo, recRepo)
		recId2, err := testhelpers.CreateRecording(recRepo)

		id1 := testhelpers.CreateTestTrack(t, trackRepo, albumId, recId1, 1, 1)
		id2 := testhelpers.CreateTestTrack(t, trackRepo, albumId, recId2, 2, 1)

		got, err := trackRepo.GetTracksByIds(ctx, []string{id1.String(), id2.String()})

		require.NoError(t, err)
		assert.Len(t, got, 2)
	})

	t.Run("miss id", func(t *testing.T) {
		albumId, recId := setTrackDeps(t, albumRepo, recRepo)
		id := testhelpers.CreateTestTrack(t, trackRepo, albumId, recId, 1, 1)

		got, err := trackRepo.GetTracksByIds(ctx, []string{id.String(), uuid.New().String()})

		require.NoError(t, err)
		assert.Len(t, got, 1)
	})

	t.Run("empty slice", func(t *testing.T) {
		got, err := trackRepo.GetTracksByIds(ctx, []string{})

		require.NoError(t, err)
		assert.Empty(t, got)
	})
}

func TestTrackRepo_GetArtistTracks(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	albumRepo := albumrepo.NewAlbumRepo(db)
	recRepo := recordingrepo.NewRecordingRepo(db)
	trackRepo := trackrepo.NewTrackRepo(db)
	artistRepo := artistrepo.NewArtistRepo(db)
	ctx := context.Background()

	t.Run("all artist tracks", func(t *testing.T) {
		artistId, err := testhelpers.CreateArtist(artistRepo)
		albumId, recId := setTrackDeps(t, albumRepo, recRepo)
		trackId, err := testhelpers.CreateTrack(trackRepo, albumId, recId, 1, 1)

		_, err = db.Exec(ctx,
			`INSERT INTO artist_tracks (artist_id, track_id) VALUES ($1, $2)`,
			artistId, trackId,
		)
		require.NoError(t, err)

		got, err := trackRepo.GetArtistTracks(ctx, artistId.String())

		require.NoError(t, err)
		require.Len(t, got, 1)
		assert.Equal(t, trackId, got[0].Id)
	})

	t.Run("artist with no tracks", func(t *testing.T) {
		artistId, err := testhelpers.CreateArtist(artistRepo)

		got, err := trackRepo.GetArtistTracks(ctx, artistId.String())

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("bad artist id", func(t *testing.T) {
		_, err := trackRepo.GetArtistTracks(ctx, "badId")
		assert.Error(t, err)
	})
}
