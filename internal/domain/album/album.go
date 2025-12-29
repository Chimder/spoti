package album

import "github.com/google/uuid"

type Album struct {
	ID               uuid.UUID
	AlbumType        string
	TotalTracks      int
	AvailableMarkets []string
	Images           string
	Name             string
	ReleaseDate      string
	URI              string
	Artists          []AlbumArtist
	Copyrights []struct {
		Text string
		Type string
	}
	ExternalIds struct {
		Upc string
	}
	Genres     []interface{}
	Label      string
	Popularity int
}
