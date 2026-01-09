package album

import "github.com/google/uuid"

type Album struct {
	ID          uuid.UUID
	AlbumType   string //"album", "single", "compilation"
	TotalTracks int
	Images      string
	Name        string
	ReleaseDate string
	URI         string
	Tracks      []AlbumTracks
	Artists     []AlbumArtist
	Copyrights  string
	Genres      []string
	Label       string
	Popularity  int
}
