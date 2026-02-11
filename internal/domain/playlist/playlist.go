package playlist

import (
	"time"

	"github.com/google/uuid"
)

type Playlist struct {
	Id          uuid.UUID
	Owner       uuid.UUID
	Name        string
	Description string
	DiscNumber  int
	Img         string
	Public      bool
	Total       uint
	Tracks      []PlayListTrack
}

type PlaylistJson struct {
	Id          uuid.UUID      `json:"id"`
	Owner       uuid.UUID      `json:"owner_id"`
	Name        string         `json:"playlist_name"`
	Description string         `json:"description"`
	Img         string         `json:"image"`
	Public      bool           `json:"is_public"`
	Total       uint           `json:"total"`
	Tracks      PlaylistTracks `json:"tracks"`
}

type PlaylistTracks struct {
	Total  uint            `json:"total"`
	Limit  int             `json:"limit"`
	Offset int             `json:"offset"`
	Items  []PlayListTrack `json:"items"`
}

type PlayListTrack struct {
	AddedAt time.Time `json:"added_at"`
	IsLocal bool      `json:"is_local"`
	Track   Track     `json:"track"`
}

type Track struct {
	Id          uuid.UUID `json:"id"`
	Name        string    `json:"track_name"`
	TrackNumber int       `json:"track_number"`
	DiscNumber  int       `json:"disc_number"`
	Explicit    bool      `json:"explicit"`
	Artists     []Artist  `json:"artists"`
}

type Artist struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Uri  string    `json:"uri"`
}
