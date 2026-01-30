package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
)

type fakeUser struct {
	name          string
	email         string
	image         string
	followers     uint32
	premiumStatus bool
}

func GetFakeUser() fakeUser {
	return fakeUser{
		name:          gofakeit.Username(),
		email:         gofakeit.Email(),
		image:         fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		followers:     uint32(gofakeit.Number(0, 1_000_000)),
		premiumStatus: gofakeit.Bool(),
	}
}

type fakeArtist struct {
	url         string
	uri         string
	artist_name string
	image       string
	followers   int
	popularity  int
	genres      []string
}

func GetFakeArtist() fakeArtist {
	genres := []string{"rock", "pop", "jazz", "hip-hop", "electronic", "classical", "indie", "metal", "folk", "r&b"}
	numGenres := gofakeit.Number(1, 4)
	newGenres := make([]string, numGenres)
	for i := range numGenres {
		newGenres[i] = genres[rand.Intn(len(genres))]
	}

	artistID := uuid.New().String()
	return fakeArtist{
		url:         fmt.Sprintf("https://open.spotify.com/artist/%s", artistID),
		uri:         fmt.Sprintf("spotify:artist:%s", artistID),
		artist_name: gofakeit.Name(),
		image:       fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		followers:   gofakeit.Number(100, 5000000),
		popularity:  gofakeit.Number(0, 100),
		genres:      newGenres,
	}
}

type fakeAlbums struct {
	album_type   string
	total_tracks int
	image        string
	album_name   string
	uri          string
	copyrights   string
	album_label  string
	popularity   int
	release_date time.Time
}

func GetFakeAlbums() fakeAlbums {
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
	return fakeAlbums{
		album_type:   albumType,
		total_tracks: totalTracks,
		image:        fmt.Sprintf("https://i.scdn.co/image/%s", gofakeit.UUID()),
		album_name:   albumName,
		uri:          fmt.Sprintf("spotify:album:%s", albumID),
		copyrights:   fmt.Sprintf("© %d %s", releaseDate.Year(), gofakeit.Company()),
		album_label:  gofakeit.Company(),
		popularity:   gofakeit.Number(0, 100),
		release_date: releaseDate,
	}
}

type FakeRecordings struct {
	isrc        string
	duration_ms int
	popularity  int
	play_count  uint64
	audio_uri   string
	preview_uri string
}

func GetFakeRecordings() FakeRecordings {
	isrc := fmt.Sprintf("US%s%02d%05d",
		gofakeit.LetterN(3),
		gofakeit.Number(0, 99),
		gofakeit.Number(0, 99999))

	durationMs := gofakeit.Number(30000, 900000)
	recordingID := uuid.New().String()

	return FakeRecordings{
		isrc:        isrc,
		duration_ms: durationMs,
		popularity:  gofakeit.Number(0, 100),
		play_count:  uint64(gofakeit.Number(0, 2_000_000_000)),
		audio_uri:   fmt.Sprintf("https://audio.cdn.example.com/%s.mp3", recordingID),
		preview_uri: fmt.Sprintf("https://audio.cdn.example.com/preview/%s.mp3", recordingID),
	}
}
