package album

import (
	"time"

	"github.com/google/uuid"
)

// type Album struct {
// 	ID          uuid.UUID
// 	AlbumType   string
// 	TotalTracks int
// 	Images      string
// 	Name        string
// 	ReleaseDate string
// 	URI         string
// 	Tracks      []AlbumTracks
// 	Artists     []AlbumArtist
// 	Copyrights  string
// 	Genres      []string
// 	Label       string
// 	Popularity  int
// }

type Album struct {
	ID          uuid.UUID
	AlbumType   string
	TotalTracks int
	Image       string
	Name        string
	URI         string
	Copyrights  string
	Label       string
	Popularity  int16
	ReleaseDate time.Time
	CreatedAt   time.Time
}
