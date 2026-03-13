package meilisearchrepo_test

import (
	"context"
	"fmt"
	meilisearchrepo "spoti/internal/repository/meilisearch"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func setupMeili(t *testing.T) (*meilisearchrepo.MeiliRepository, testcontainers.Container) {

	ctx := context.Background()

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "getmeili/meilisearch:latest",
			ExposedPorts: []string{"7700/tcp"},
			Env: map[string]string{
				"MEILI_ENV": "development",
			},
			WaitingFor: wait.ForHTTP("/health").
				WithPort("7700/tcp"),
		},
		Started: true,
	})
	if err != nil {
		t.Fatal(err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		t.Fatal(err)
	}

	port, err := container.MappedPort(ctx, "7700")
	if err != nil {
		t.Fatal(err)
	}

	url := fmt.Sprintf("http://%s:%s", host, port.Port())

	conn := meilisearchrepo.NewMeiliDB(url)
	repo := meilisearchrepo.NewMeiliRepository(conn)

	return repo, container
}

func TestIndexAndSearch(t *testing.T) {

	repo, container := setupMeili(t)

	defer container.Terminate(context.Background())

	ctx := context.Background()

	docs := []meilisearchrepo.Document{
		{ID: "1", Type: "artist", Name: "Drake"},
		{ID: "2", Type: "artist", Name: "Drake Bell"},
		{ID: "3", Type: "album", Name: "Views"},
		{ID: "4", Type: "track", Name: "Hotline Bling"},
	}

	for _, d := range docs {
		err := repo.Add(ctx, d)
		require.NoError(t, err)
	}

	time.Sleep(2 * time.Second)

	t.Run("search artist", func(t *testing.T) {
		res, err := repo.Search(ctx, "drake")
		require.NoError(t, err)
		require.NotEmpty(t, res)
	})

	t.Run("search artist type", func(t *testing.T) {
		res, err := repo.SearchByType(ctx, "drake", "artist")
		require.NoError(t, err)
		require.Len(t, res, 2)

		for _, r := range res {
			require.Equal(t, "artist", r.Type)
		}
	})

	t.Run("search album type", func(t *testing.T) {
		res, err := repo.SearchByType(ctx, "views", "album")
		require.NoError(t, err)
		require.Len(t, res, 1)
	})
}
