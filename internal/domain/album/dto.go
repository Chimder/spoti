package album

import "time"

type CreateAlbumReq struct {
	AlbumType   string
	TotalTracks int
	Image       string
	AlbumName   string
	Uri         string
	Copyrights  string
	AlbumLabel  string
	Popularity  int
	ReleaseDate time.Time
}
