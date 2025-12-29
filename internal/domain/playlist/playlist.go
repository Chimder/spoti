package playlist

import "github.com/google/uuid"

type Playlist struct {
	Id          uuid.UUID
	Owner       uuid.UUID
	Name        string
	Description string
	DiscNumber  int
	Img         string
	Public      bool
	Total       uint
}
