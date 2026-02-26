package testhelpers

import (
	albumrepo "spoti/internal/repository/postgres/album"
	artistrepo "spoti/internal/repository/postgres/artist"
	playlistrepo "spoti/internal/repository/postgres/playlist"
	recordingrepo "spoti/internal/repository/postgres/recording"
	trackrepo "spoti/internal/repository/postgres/track"
	userrepo "spoti/internal/repository/postgres/user"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func CreateTestUser(t *testing.T, repo *userrepo.UserRepo) uuid.UUID {
	t.Helper()

	id, err := CreateUser(repo)
	require.NoError(t, err)
	return id
}

func CreateTestArtist(t *testing.T, repo *artistrepo.ArtistRepo) uuid.UUID {
	t.Helper()

	id, err := CreateArtist(repo)
	require.NoError(t, err)

	return id
}

func CreateTestAlbum(t *testing.T, repo *albumrepo.AlbumRepo) uuid.UUID {
	t.Helper()

	id, err := CreateAlbum(repo)
	require.NoError(t, err)

	return id
}

func CreateTestPlaylist(t *testing.T, repo *playlistrepo.PlaylistRepo, owner uuid.UUID) uuid.UUID {
	t.Helper()

	id, err := CreatePlaylist(repo, owner)
	require.NoError(t, err)

	return id
}

func CreateTestRecording(t *testing.T, repo *recordingrepo.RecordingRepo) uuid.UUID {
	t.Helper()

	id, err := CreateRecording(repo)
	require.NoError(t, err)

	return id
}

func CreateTestTrack(t *testing.T, repo *trackrepo.TrackRepo, albumId, recordingId uuid.UUID, trackNum, discNum int) uuid.UUID {
	t.Helper()

	id, err := CreateTrack(repo, albumId, recordingId, trackNum, discNum)
	require.NoError(t, err)

	return id
}
