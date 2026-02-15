package artistrepo

import (
	"context"
	"fmt"
	"spoti/internal/domain/artist"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

type ArtistRepo struct {
	db *pgxpool.Pool
}

func NewArtistRepo(db *pgxpool.Pool) *ArtistRepo {
	return &ArtistRepo{
		db: db,
	}
}

func (art *ArtistRepo) CreateArtist(ctx context.Context, a artist.CreateArtistReq) (uuid.UUID, error) {
	query := `
	INSERT INTO artists (url, uri, artist_name, image)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`
	var id uuid.UUID
	err := art.db.QueryRow(ctx, query, a.Url, a.Uri, a.ArtistName, a.Image).Scan(&id)
	if err != nil {
		log.Error().Err(err).Msg("err create artist")
		return uuid.UUID{}, err
	}

	return id, err
}

func (art *ArtistRepo) GetArtist(ctx context.Context, artistId string) (artist.Artist, error) {
	query := `SELECT * FROM artists WHERE id = $1`

	rows, err := art.db.Query(ctx, query, artistId)
	if err != nil {
		return artist.Artist{}, fmt.Errorf("err fetch artist by id %w", err)
	}
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[Artist])
	if err != nil {
		return artist.Artist{}, fmt.Errorf("err row artist %w", err)
	}
	return data.ToDomain(), nil
}

func (art *ArtistRepo) GetArtistsByIDs(ctx context.Context, artistIds []string) ([]artist.Artist, error) {
	query := `
        SELECT * FROM artists
        WHERE id = ANY($1)
        ORDER BY created_at DESC
    `

	rows, err := art.db.Query(ctx, query, artistIds)
	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Artist])
	if err != nil {
		return nil, fmt.Errorf("err row artist %w", err)
	}

	return ArtistsToDomain(data), nil
}

func (art *ArtistRepo) GetArtistAlbums(ctx context.Context, artistId string) ([]artist.Artist, error) {
	query := `
	SELECT a.*
  FROM albums a
  JOIN album_artists aa ON a.id = aa.album_id
  WHERE aa.artist_id = $1
	`

	rows, err := art.db.Query(ctx, query, artistId)
	data, err := pgx.CollectRows(rows, pgx.RowToStructByName[Artist])
	if err != nil {
		return nil, fmt.Errorf("err row artist %w", err)
	}

	return ArtistsToDomain(data), nil
}
