package albumrepo

import (
	"spoti/internal/domain/album"
	"time"

	"github.com/google/uuid"
)

type Album struct {
	ID          uuid.UUID `db:"id"`
	AlbumType   string    `db:"album_type"`
	TotalTracks int     `db:"total_tracks"`
	Image       string    `db:"image"`
	AlbumName   string    `db:"album_name"`
	URI         string    `db:"uri"`
	Copyrights  string    `db:"copyrights"`
	AlbumLabel  string    `db:"album_label"`
	Popularity  int16     `db:"popularity"`
	ReleaseDate time.Time `db:"release_date"`
	CreatedAt   time.Time `db:"created_at"`
}
func (a *Album) ToDomain() album.Album {
	return album.Album{
		ID:          a.ID,
		AlbumType:   a.AlbumType,
		TotalTracks: a.TotalTracks,
		Image:       a.Image,
		Name:        a.AlbumName,
		URI:         a.URI,
		Copyrights:  a.Copyrights,
		Label:       a.AlbumLabel,
		Popularity:  a.Popularity,
		ReleaseDate: a.ReleaseDate,
		CreatedAt:   a.CreatedAt,
	}
}

func AlbumsToDomain(albums []Album) []album.Album {
	result := make([]album.Album, len(albums))
	for i, a := range albums {
		result[i] = a.ToDomain()
	}
	return result
}
