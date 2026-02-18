package testhelpers

import (
	"context"
	"fmt"
	"math/rand"
	"spoti/internal/domain/album"
	"spoti/internal/domain/artist"
	"spoti/internal/domain/playlist"
	"spoti/internal/domain/recording"
	"spoti/internal/domain/track"
	"spoti/internal/domain/user"
	albumrepo "spoti/internal/repository/postgres/album"
	artistrepo "spoti/internal/repository/postgres/artist"
	playlistrepo "spoti/internal/repository/postgres/playlist"
	recordingrepo "spoti/internal/repository/postgres/recording"
	trackrepo "spoti/internal/repository/postgres/track"
	userrepo "spoti/internal/repository/postgres/user"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func CreateTestUser(t *testing.T, repo *userrepo.UserRepo) uuid.UUID {
	t.Helper()

	id, err := repo.CreateUser(context.Background(), GetFakeUser())
	require.NoError(t, err)

	return id
}
func CreateTestArtist(t *testing.T, repo *artistrepo.ArtistRepo) uuid.UUID {
	t.Helper()

	id, err := repo.CreateArtist(context.Background(), GetFakeArtist())
	require.NoError(t, err)

	return id
}
func CreateTestAlbum(t *testing.T, repo *albumrepo.AlbumRepo) uuid.UUID {
	t.Helper()

	id, err := repo.CreateAlbum(context.Background(), GetFakeAlbums())
	require.NoError(t, err)

	return id
}
func CreateTestPlaylist(t *testing.T, repo *playlistrepo.PlaylistRepo, owner uuid.UUID) uuid.UUID {
	t.Helper()

	id, err := repo.CreatePlaylist(context.Background(), GetFakePlaylist(owner))
	require.NoError(t, err)

	return id
}

func CreateTestRecording(t *testing.T, repo *recordingrepo.RecordingRepo) uuid.UUID {
	t.Helper()

	id, err := repo.CreateRecording(context.Background(), GetFakeRecording())
	require.NoError(t, err)

	return id
}

func CreateTestTrack(t *testing.T, repo *trackrepo.TrackRepo, albumId, recordingId uuid.UUID, trackNum, discNum int) uuid.UUID {
	t.Helper()

	id, err := repo.CreateTrack(context.Background(), GetFakeTrack(albumId, recordingId, trackNum, discNum))
	require.NoError(t, err)

	return id
}

func GetFakeUser() user.CreateUserReq {
	return user.CreateUserReq{
		Name:          gofakeit.Username(),
		Email:         gofakeit.Email(),
		Image:         fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		Followers:     uint32(gofakeit.Number(0, 1_000_000)),
		PremiumStatus: gofakeit.Bool(),
	}
}

func GetFakeArtist() artist.CreateArtistReq {
	genres := []string{"rock", "pop", "jazz", "hip-hop", "electronic", "classical", "indie", "metal", "folk", "r&b"}
	numGenres := gofakeit.Number(1, 4)
	newGenres := make([]string, numGenres)
	for i := range numGenres {
		newGenres[i] = genres[rand.Intn(len(genres))]
	}

	artistID := uuid.New().String()
	return artist.CreateArtistReq{
		Url:        fmt.Sprintf("https://open.spotify.com/artist/%s", artistID),
		Uri:        fmt.Sprintf("spotify:artist:%s", artistID),
		ArtistName: gofakeit.Name(),
		Image:      fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		Followers:  gofakeit.Number(100, 5000000),
		Popularity: gofakeit.Number(0, 100),
		Genres:     newGenres,
	}
}

func GetFakeAlbums() album.CreateAlbumReq {
	albumTypes := []string{"album", "single", "compilation"}
	albumType := albumTypes[rand.Intn(len(albumTypes))]
	var totalTracks int
	switch albumType {
	case "single":
		totalTracks = gofakeit.Number(1, 2)
	case "album":
		totalTracks = gofakeit.Number(2, 15)
	case "compilation":
		totalTracks = gofakeit.Number(15, 40)
	}
	albumID := uuid.New().String()
	albumName := fmt.Sprintf("%s - %s (%s)",
		gofakeit.BuzzWord(),
		gofakeit.Noun(),
		uuid.New().String()[:8])
	releaseDate := gofakeit.DateRange(
		time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Now(),
	)
	return album.CreateAlbumReq{
		AlbumType:   albumType,
		TotalTracks: totalTracks,
		Image:       fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		AlbumName:   albumName,
		Uri:         fmt.Sprintf("spotify:album:%s", albumID),
		Copyrights:  fmt.Sprintf("© %d %s", releaseDate.Year(), gofakeit.Company()),
		AlbumLabel:  gofakeit.Company(),
		Popularity:  gofakeit.Number(0, 100),
		ReleaseDate: releaseDate,
	}
}

func GetFakePlaylist(OwnerId uuid.UUID) playlist.CreatePlaylistReq {
	return playlist.CreatePlaylistReq{
		OwnerId:      OwnerId,
		PlaylistName: gofakeit.Sentence(2),
		Description:  gofakeit.Sentence(10),
		Image:        fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		IsPublic:     gofakeit.Bool(),
	}
}

func GetFakeRecording() recording.CreateRecordingReq {
	isrc := fmt.Sprintf("US%s%02d%05d", gofakeit.LetterN(3), gofakeit.Number(0, 99), gofakeit.Number(0, 99999))
	recordingID := uuid.New().String()
	return recording.CreateRecordingReq{
		ISRC:       isrc,
		DurationMs: int64(gofakeit.Number(30000, 900000)),
		AudioUri:   fmt.Sprintf("https://audio.cdn.example.com/%s.mp3", recordingID),
	}
}

func GetFakeTrack(albumId, recordingId uuid.UUID, trackNum, discNum int) track.CreateTrackReq {
	return track.CreateTrackReq{
		AlbumId:     albumId,
		RecordingId: recordingId,
		Name:        fmt.Sprintf("%s %s", gofakeit.Adjective(), gofakeit.Noun()),
		Number:      int16(trackNum),
		DiscNumber:  int16(discNum),
		Explicit:    gofakeit.Bool(),
		IsPlayable:  true,
		Type:        "track",
		URI:         fmt.Sprintf("spotify:track:%s", uuid.New()),
		IsLocal:     false,
	}
}
