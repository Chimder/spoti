package album

import (
	"time"

	"github.com/google/uuid"
)

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

type GetAlbumsByIdsResponse struct {
	Albums []GetAlbumResponse `json:"albums"`
}

type GetAlbumResponse struct {
	AlbumType   string          `json:"album_type"`
	TotalTracks int             `json:"total_tracks"`
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	ReleaseDate time.Time       `json:"release_date"`
	URI         string          `json:"uri"`
	Artists     []ArtistSummary `json:"artists"`
	Tracks      AlbumTracksDTO  `json:"tracks"`
}

type ArtistSummary struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	URI  string    `json:"uri"`
}

type AlbumTracksDTO struct {
	Items []TrackSummary `json:"items"`
}

type TrackSummary struct {
	ID          uuid.UUID       `json:"id"`
	Name        string          `json:"name"`
	TrackNumber int             `json:"track_number"`
	DiscNumber  int             `json:"disc_number"`
	DurationMs  int             `json:"duration_ms"`
	Explicit    bool            `json:"explicit"`
	URI         string          `json:"uri"`
	Artists     []ArtistSummary `json:"artists"`
}
