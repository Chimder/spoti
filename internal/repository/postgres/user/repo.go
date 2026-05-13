package userrepo

import (
	"context"

	"github.com/Chimder/spoti/internal/domain/user"
	"github.com/Chimder/spoti/internal/repository/postgres/pgiface"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/rs/zerolog/log"
)

type UserRepository interface {
	CreateUser(ctx context.Context, user user.CreateUserReq, hashPass string) (uuid.UUID, error)
	GetUserById(ctx context.Context, userId uuid.UUID) (user.User, error)
	GetUserByEmail(ctx context.Context, userEmail string) (user.User, error)
	FollowUserToPlaylist(ctx context.Context, userId, playlistId uuid.UUID) error
	UnfollowUserFromPlaylist(ctx context.Context, userId, playlistId uuid.UUID) error
	FollowUserToArtist(ctx context.Context, userId, artistId uuid.UUID) error
	UnfollowUserFromArtist(ctx context.Context, userId, artistId uuid.UUID) error
}
type UserRepo struct {
	db pgiface.Querier
}

func NewUserRepo(db pgiface.Querier) *UserRepo {
	return &UserRepo{
		db: db,
	}
}

func (ur *UserRepo) CreateUser(ctx context.Context, user user.CreateUserReq, hashPass string) (uuid.UUID, error) {
	query := `
			INSERT INTO users (user_name, email, image, password_hash)
			VALUES ($1, $2, $3, $4)
			RETURNING id
		`

	var id uuid.UUID
	err := ur.db.QueryRow(ctx, query, user.Name, user.Email, user.Image, hashPass).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("err create user")
		return uuid.UUID{}, err
	}

	return id, err
}

func (ur *UserRepo) GetUserByEmail(ctx context.Context, userEmail string) (user.User, error) {
	rows, err := ur.db.Query(ctx, "SELECT * FROM users WHERE email = $1", userEmail)
	if err != nil {
		log.Error().Err(err).Msg("err get user_pass by email")
		return user.User{}, err
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		log.Error().Err(err).Msg("err collect user info by id")
		return user.User{}, err
	}

	return data.ToDomain(), nil
}

func (ur *UserRepo) GetUserById(ctx context.Context, userId uuid.UUID) (user.User, error) {
	rows, err := ur.db.Query(ctx, "SELECT * FROM users WHERE id = $1", userId)
	if err != nil {
		log.Error().Err(err).Msg("err get user by id")
		return user.User{}, err
	}

	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[User])
	if err != nil {
		log.Error().Err(err).Msg("err collect user info by id")
		return user.User{}, err
	}

	return data.ToDomain(), nil
}

func (ur *UserRepo) FollowUserToPlaylist(ctx context.Context, userId, playlistId uuid.UUID) error {
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
func (ur *UserRepo) UnfollowUserFromPlaylist(ctx context.Context, userId, playlistId uuid.UUID) error {
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

func (ur *UserRepo) FollowUserToArtist(ctx context.Context, userId, artistId uuid.UUID) error {
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
func (ur *UserRepo) UnfollowUserFromArtist(ctx context.Context, userId, artistId uuid.UUID) error {
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
