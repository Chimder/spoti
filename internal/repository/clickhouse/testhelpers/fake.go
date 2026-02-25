package testhelpers

import (
	"context"
	"spoti/internal/domain/user"
	"spoti/internal/repository/clickhouse"
	"testing"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func createFakeListeningEvent() user.ListeningEventReq {
	return user.ListeningEventReq{
		UserId:     uuid.New(),
		TrackId:    uuid.New(),
		ArtistId:   uuid.New(),
		AlbumId:    uuid.New(),
		DurationMs: gofakeit.Number(90_000, 300_000),
		IsSkipped:  gofakeit.Bool(),
	}
}

func CreateFakeListeningEvent(t *testing.T, repo *clickhouse.ListeningEventRepo) user.ListeningEventReq {
	t.Helper()

	event := createFakeListeningEvent()
	err := repo.AddListeningEvent(context.Background(), event)
	require.NoError(t, err)

	return event
}

// func GetFakeListeningEvent(t *testing.T, repo *clickhouse.ListeningEventRepo) user.ListeningEventReq {
// 	t.Helper()

// 	event,err := repo.GetListeningEvent(context.Background(), event)
// 	require.NoError(t, err)

// 	return event
// }
