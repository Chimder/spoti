package testhelpers

import (
	"context"
	"fmt"
	"math/rand"

	"time"

	"github.com/Chimder/spoti/internal/domain/album"
	"github.com/Chimder/spoti/internal/domain/artist"
	"github.com/Chimder/spoti/internal/domain/playlist"
	"github.com/Chimder/spoti/internal/domain/recording"
	"github.com/Chimder/spoti/internal/domain/track"
	"github.com/Chimder/spoti/internal/domain/user"
	"github.com/Chimder/spoti/internal/handler/http/middleware"
	albumrepo "github.com/Chimder/spoti/internal/repository/postgres/album"
	artistrepo "github.com/Chimder/spoti/internal/repository/postgres/artist"
	playlistrepo "github.com/Chimder/spoti/internal/repository/postgres/playlist"
	recordingrepo "github.com/Chimder/spoti/internal/repository/postgres/recording"
	trackrepo "github.com/Chimder/spoti/internal/repository/postgres/track"
	userrepo "github.com/Chimder/spoti/internal/repository/postgres/user"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

func CreateUser(repo userrepo.UserRepository) (uuid.UUID, error) {
	user := FakeUser()
	hashPass, err := middleware.GeneratePass(user.Password)
	if err != nil {
		return uuid.Nil, err
	}
	return repo.CreateUser(context.Background(), user, hashPass)
}

func FakeUser() user.CreateUserReq {
	return user.CreateUserReq{
		Name:          gofakeit.Username(),
		Email:         gofakeit.Email(),
		Password:      gofakeit.Password(false, false, false, false, false, 16),
		Image:         fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		Followers:     uint32(gofakeit.Number(0, 1_000_000)),
		PremiumStatus: gofakeit.Bool(),
	}
}

func CreateArtist(repo artistrepo.ArtistRepository) (uuid.UUID, error) {
	return repo.CreateArtist(context.Background(), FakeArtist())
}

func FakeArtist() artist.CreateArtistReq {
	genres := []string{"rock", "pop", "jazz", "hip-hop", "electronic", "classical", "indie", "metal", "folk", "r&b"}
	numGenres := gofakeit.Number(1, 4)
	newGenres := make([]string, numGenres)
	for i := range newGenres {
		newGenres[i] = genres[rand.Intn(len(genres))]
	}

	artistID := uuid.New().String()
	return artist.CreateArtistReq{
		Url:        fmt.Sprintf("https://open.spotify.com/artist/%s", artistID),
		Uri:        fmt.Sprintf("spotify:artist:%s", artistID),
		ArtistName: gofakeit.Name() + "_" + uuid.New().String()[:8],
		Image:      fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		Followers:  gofakeit.Number(100, 5000000),
		Popularity: gofakeit.Number(0, 100),
		Genres:     newGenres,
	}
}

func CreateAlbum(repo albumrepo.AlbumRepository) (uuid.UUID, error) {
	return repo.CreateAlbum(context.Background(), FakeAlbum())
}

func FakeAlbum() album.CreateAlbumReq {
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
	releaseDate := gofakeit.DateRange(time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC), time.Now())
	albumName := fmt.Sprintf("%s - %s (%s)", gofakeit.BuzzWord(), gofakeit.Noun(), albumID[:8])

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

func CreatePlaylist(repo playlistrepo.PlaylistRepository, owner uuid.UUID) (uuid.UUID, error) {
	return repo.CreatePlaylist(context.Background(), FakePlaylist(owner))
}

func FakePlaylist(owner uuid.UUID) playlist.CreatePlaylistReq {
	return playlist.CreatePlaylistReq{
		OwnerId:      owner,
		PlaylistName: gofakeit.Sentence(2),
		Description:  gofakeit.Sentence(10),
		Image:        fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		IsPublic:     gofakeit.Bool(),
	}
}

func CreateRecording(repo recordingrepo.RecordingRepository) (uuid.UUID, error) {
	return repo.CreateRecording(context.Background(), FakeRecording())
}

func FakeRecording() recording.CreateRecordingReq {
	isrc := fmt.Sprintf("US%s%02d%05d", gofakeit.LetterN(3), gofakeit.Number(0, 99), gofakeit.Number(0, 99999))
	recordingID := uuid.New().String()
	return recording.CreateRecordingReq{
		ISRC:       isrc,
		DurationMs: int64(gofakeit.Number(30000, 900000)),
		AudioUri:   fmt.Sprintf("https://audio.cdn.example.com/%s.mp3", recordingID),
	}
}

func CreateTrack(repo trackrepo.TrackRepository, albumID, recordingID uuid.UUID, trackNum, discNum int) (uuid.UUID, error) {
	return repo.CreateTrack(context.Background(), FakeTrack(albumID, recordingID, trackNum, discNum))
}

func FakeTrack(albumID, recordingID uuid.UUID, trackNum, discNum int) track.CreateTrackReq {
	return track.CreateTrackReq{
		AlbumId:     albumID,
		RecordingId: recordingID,
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
