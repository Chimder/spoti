package artist

import (
	"time"

	"github.com/google/uuid"
)

type Artist struct {
	Id         uuid.UUID
	Url        string
	Followers  uint64
	Genres     []string
	Image      string
	Name       string
	Popularity uint8
	URI        string
	CreatedAt  time.Time
}
