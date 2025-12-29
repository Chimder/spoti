package user

import (
	"github.com/google/uuid"
)

type User struct {
	id            uuid.UUID
	name          string
	email         string
	image         string
	followers     uint32
	premiumStatus bool
}
