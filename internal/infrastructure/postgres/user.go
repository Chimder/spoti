package postgres

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type User struct {
	ID            uuid.UUID `db:"id"`
	UserName      string    `db:"user_name"`
	Email         string    `db:"email"`
	Image         string    `db:"image"`
	Followers     int64     `db:"followers"`
	CreatedAt     time.Time `db:"created_at"`
	PremiumStatus bool      `db:"premium_status"`
}

type UserRepo struct {
	db *pgxpool.Pool
}

func NewUserRepo(db *pgxpool.Pool) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) GetUserById(ctx context.Context, userId string) (User, error) {

	rows, err := ur.db.Query(ctx, "SELECT * FROM users WHERE id = $1", userId)
	if err != nil {
		log.Error().Err(err).Msg("err get user by id")
		return User{}, err
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		log.Error().Err(err).Msg("err collect user info by id")
		return User{}, err
	}

	return data, nil
}

func (ur *UserRepo) FollowUserToPlaylist(ctx context.Context, userId, playlistId string) error {
	query := `
	INSERT INTO user_saved_playlists (user_id, playlist_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`

	_, err := ur.db.Exec(ctx, query, userId, playlistId)
	if err != nil {
		log.Error().Err(err).Msg("err follow user to playlist")
		return err
	}

	return nil
}
func (ur *UserRepo) UnfollowUserFromPlaylist(ctx context.Context, userId, playlistId string) error {
	query := `
	DELETE FROM user_saved_playlists
	WHERE user_id = $1 AND playlist_id = $2
	`

	_, err := ur.db.Exec(ctx, query, userId, playlistId)
	if err != nil {
		log.Error().Err(err).Msg("err unfollow user from playlist")
		return err
	}

	return nil
}

func (ur *UserRepo) FollowUserToArtist(ctx context.Context, userId, artistId string) error {
	query := `
	INSERT INTO user_saved_artists (user_id, artist_id)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	`

	_, err := ur.db.Exec(ctx, query, userId, artistId)
	if err != nil {
		log.Error().Err(err).Msg("err follow user to artist")
		return err
	}

	return nil
}
func (ur *UserRepo) UnfollowUserFromArtist(ctx context.Context, userId, artistId string) error {
	query := `
	DELETE FROM user_saved_artists
	WHERE user_id = $1 AND artist_id = $2
	`

	_, err := ur.db.Exec(ctx, query, userId, artistId)
	if err != nil {
		log.Error().Err(err).Msg("err unfollow user from artist")
		return err
	}

	return nil
}
