package track

import "spoti/internal/domain/artist"

type Track struct {
	ID               string
	AvailableMarkets []string
	Explicit         bool
	IsPlayable       bool
	Name             string
	Popularity       int
	PreviewURL       string
	DiscNumber       int
	TrackNumber      int
	DurationMs       uint32
	Type             string
	URI              string
	IsLocal          bool
	Artists          artist.Artist
}
