package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistRepo struct {
	db *pgxpool.Pool
}

type Artist struct {
	ID         uuid.UUID `db:"id"`
	URL        string    `db:"url"`
	URI        string    `db:"uri"`
	ArtistName string    `db:"artist_name"`
	Image      *string   `db:"image"`
	Followers  int64     `db:"followers"`
	Popularity int16     `db:"popularity"`
	Genres     []string  `db:"genres"`
	CreatedAt  time.Time `db:"created_at"`
}

func NewArtistRepo(db *pgxpool.Pool) *AlbumRepo {
	return &AlbumRepo{
		db: db,
	}
}

func (art *ArtistRepo) GetArtist(ctx context.Context, artistId string) (Artist, error) {
	query := `SELECT * FROM artists WHERE id = $1`

	rows, err := art.db.Query(ctx, query, artistId)
	if err != nil {
		return Artist{}, fmt.Errorf("err fetch artist by id %w", err)
	}
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Artist])
	if err != nil {
		return Artist{}, fmt.Errorf("err row artist %w", err)
	}
	return data, nil
}

func (art *ArtistRepo) GetArtistsByIDs(ctx context.Context, artistIds []string) ([]Artist, error) {
	query := `
        SELECT * FROM artists
        WHERE id = ANY($1)
        ORDER BY created_at DESC
    `

	rows, err := art.db.Query(ctx, query, artistIds)
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[[]Artist])
	if err != nil {
		return nil, fmt.Errorf("err row artist %w", err)
	}
	return data, nil
}

func (art *ArtistRepo) GetArtistAlbums(ctx context.Context, artistId string) ([]Artist, error) {
	query := `
	SELECT a.*
  FROM albums a
  JOIN album_artists aa ON a.id = aa.album_id
  WHERE aa.artist_id = $1
	`

	rows, err := art.db.Query(ctx, query, artistId)
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[[]Artist])
	if err != nil {
		return nil, fmt.Errorf("err row artist %w", err)
	}
	return data, nil
}
func (art *ArtistRepo) GetArtistTracks(ctx context.Context, artistId string) ([]Track, error) {
	query := `
	SELECT t.*
  FROM tracks t
  JOIN artist_tracks at ON at.track_id = t.id
  WHERE at.artist_id = $1
	`

	rows, err := art.db.Query(ctx, query, artistId)
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[[]Track])
	if err != nil {
		return nil, fmt.Errorf("err row artist %w", err)
	}
	return data, nil
}
