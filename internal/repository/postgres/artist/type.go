package artistrepo

import (
	"time"

	"github.com/Chimder/spoti/internal/domain/artist"
	"github.com/google/uuid"
)

type Artist struct {
	Id         uuid.UUID `db:"id"`
	URL        string    `db:"url"`
	URI        string    `db:"uri"`
	ArtistName string    `db:"artist_name"`
	Image      *string   `db:"image"`
	Followers  int64     `db:"followers"`
	Popularity int16     `db:"popularity"`
	Genres     []string  `db:"genres"`
	CreatedAt  time.Time `db:"created_at"`
}

func (a *Artist) ToDomain() artist.Artist {
	return artist.Artist{
		Id:         a.Id,
		Url:        a.URL,
		URI:        a.URI,
		Name:       a.ArtistName,
		Image:      *a.Image,
		Followers:  uint64(a.Followers),
		Popularity: uint8(a.Popularity),
		Genres:     a.Genres,
		CreatedAt:  a.CreatedAt,
	}
}

func ArtistsToDomain(rows []Artist) []artist.Artist {
	result := make([]artist.Artist, len(rows))
	for i, row := range rows {
		result[i] = row.ToDomain()
	}
	return result
}
