package user

import "spoti/internal/domain/track"

type UserTopTracks struct {
	Href     string
	Limit    int
	Next     string
	Offset   int
	Previous string
	Total    int
	Tracks   []track.Track
}
