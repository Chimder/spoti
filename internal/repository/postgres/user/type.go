package userrepo

import (
	"time"

	"github.com/Chimder/spoti/internal/domain/user"
	"github.com/google/uuid"
)

type User struct {
	ID            uuid.UUID `db:"id"`
	UserName      string    `db:"user_name"`
	Email         string    `db:"email"`
	Image         string    `db:"image"`
	PasswordHash  string    `db:"password_hash"`
	Followers     int64     `db:"followers"`
	CreatedAt     time.Time `db:"created_at"`
	PremiumStatus bool      `db:"premium_status"`
}

func (u *User) ToDomain() user.User {
	return user.User{
		Id:            u.ID,
		Name:          u.UserName,
		PasswordHash:  u.PasswordHash,
		Email:         u.Email,
		Image:         u.Image,
		Followers:     uint32(u.Followers),
		PremiumStatus: u.PremiumStatus,
	}
}

func UsersToDomain(rows []User) []user.User {
	result := make([]user.User, len(rows))
	for i, row := range rows {
		result[i] = row.ToDomain()
	}
	return result
}
