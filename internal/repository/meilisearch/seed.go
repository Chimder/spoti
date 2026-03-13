package meilisearchrepo

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func StartSeedMeiliSearch() error {
	fmt.Print("start seed meilisearch")

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		return err
	}
	defer pool.Close()

	conn := NewMeiliDB("http://localhost:7700")
	if conn == nil {
		return fmt.Errorf("err conn to meilisearch")
	}
	defer conn.Close()

	meiliRepo := NewMeiliRepository(conn)

	seeds := []struct {
		query string
		typ   string
	}{
		{"SELECT id, artist_name AS name FROM artists", "artist"},
		{"SELECT id, album_name AS name FROM albums", "album"},
		{"SELECT id, track_name AS name FROM tracks", "track"},
		{"SELECT id, playlist_name AS name FROM playlists", "playlist"},
	}

	for _, seed := range seeds {

		data, err := getData(ctx, pool, seed.query)
		if err != nil {
			return err
		}

		for _, item := range data {
			err := meiliRepo.Add(ctx, Document{
				ID:   item.ID,
				Type: seed.typ,
				Name: item.Name,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

type Data struct {
	ID   string `db:"id"`
	Name string `db:"name"`
}

func getData(ctx context.Context, pool *pgxpool.Pool, query string) ([]Data, error) {
	rows, err := pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return pgx.CollectRows(rows, pgx.RowToStructByName[Data])
}
