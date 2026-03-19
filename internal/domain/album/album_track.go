package album

import "github.com/Chimder/spoti/internal/domain/track"


type AlbumTracks struct {
	Href     string
	Limit    int
	Offset   int
	Next     any
	Previous any
	Total    int
	Tracks   []track.Track
}
