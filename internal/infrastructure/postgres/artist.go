package postgres

import (
	"context"
	"fmt"
	"spoti/internal/domain/artist"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ArtistRepo struct {
	db *pgxpool.Pool
}

func NewArtistRepo(db *pgxpool.Pool) *AlbumRepo {
	return &AlbumRepo{
		db: db,
	}
}

func (art *ArtistRepo) GetArtist(ctx context.Context, artistId string) (artist.Artist, error) {
	query := `SELECT * FROM artists WHERE id = $1`

	rows, err := art.db.Query(ctx, query, artistId)
	if err != nil {
		return artist.Artist{}, fmt.Errorf("err fetch artist by id %w", err)
	}
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[artist.Artist])
	if err != nil {
		return artist.Artist{}, fmt.Errorf("err row artist %w", err)
	}
	return data, nil
}
