package user

import "github.com/Chimder/spoti/internal/domain/artist"


type UserTopArtists struct {
	Href     string
	Limit    int
	Next     string
	Offset   int
	Previous string
	Total    int
	Artists  []artist.Artist
}
