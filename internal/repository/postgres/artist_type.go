package postgres

import (
	"spoti/internal/domain/artist"
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	ID         uuid.UUID `db:"id"`
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
		ID:         a.ID,
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
	for _, row := range rows {
		result = append(result, row.ToDomain())
	}
	return result
}
