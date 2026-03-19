package clickhouse_test

import (
	"context"
	"testing"

	"github.com/Chimder/spoti/internal/domain/user"
	"github.com/Chimder/spoti/internal/repository/clickhouse"
	"github.com/Chimder/spoti/internal/repository/clickhouse/testhelpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddListeningEvent(t *testing.T) {
	conn, close := testhelpers.SetupClickHouseContainer(t)
	defer close()

	repo := clickhouse.NewListeningEventRepo(conn)
	ctx := context.Background()

	t.Run("add listening event", func(t *testing.T) {
		event := testhelpers.CreateFakeListeningEvent(t, repo)

		var count uint64
		row := conn.QueryRow(ctx, `SELECT count() FROM spotify.listening_events WHERE user_id = ? AND track_id = ?`, event.UserId, event.TrackId)
		err := row.Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, uint64(1), count)
	})

	t.Run("add skipped event", func(t *testing.T) {
		event := user.ListeningEventReq{
			UserId:     uuid.New(),
			TrackId:    uuid.New(),
			ArtistId:   uuid.New(),
			AlbumId:    uuid.New(),
			DurationMs: 5000,
			IsSkipped:  true,
		}

		err := repo.AddListeningEvent(ctx, event)
		assert.NoError(t, err)

		var isSkipped bool
		row := conn.QueryRow(ctx, `SELECT is_skipped FROM spotify.listening_events WHERE user_id = ? AND track_id = ?`, event.UserId, event.TrackId)
		err = row.Scan(&isSkipped)
		require.NoError(t, err)
		assert.True(t, isSkipped)
	})
}
