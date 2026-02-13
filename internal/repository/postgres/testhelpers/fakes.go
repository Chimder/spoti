package testhelpers

import (
	"fmt"
	"math/rand"
	"spoti/internal/domain/album"
	"spoti/internal/domain/artist"
	"spoti/internal/domain/recording"
	"spoti/internal/domain/user"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

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

func GetFakeRecordings() recording.CreateRecordingReq {
	isrc := fmt.Sprintf("US%s%02d%05d",
		gofakeit.LetterN(3),
		gofakeit.Number(0, 99),
		gofakeit.Number(0, 99999))

	durationMs := gofakeit.Number(30000, 900000)
	recordingID := uuid.New().String()

	return recording.CreateRecordingReq{
		ISRC:       isrc,
		DurationMs: int64(durationMs),
		AudioUri:   fmt.Sprintf("https://audio.cdn.example.com/%s.mp3", recordingID),
	}
}
