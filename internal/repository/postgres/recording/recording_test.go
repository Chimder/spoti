package recordingrepo_test

import (
	"context"
	"testing"

	recordingrepo "github.com/Chimder/spoti/internal/repository/postgres/recording"
	"github.com/Chimder/spoti/internal/repository/postgres/testhelpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRecordingRepo_CreateRecording(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	repo := recordingrepo.NewRecordingRepo(db)
	ctx := context.Background()

	t.Run("ok create recording", func(t *testing.T) {
		id, err := testhelpers.CreateRecording(repo)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil, id)
	})

	t.Run("duplicate isrc", func(t *testing.T) {
		fake := testhelpers.FakeRecording()

		_, err := repo.CreateRecording(ctx, fake)
		require.NoError(t, err)

		_, err = repo.CreateRecording(ctx, fake)
		assert.Error(t, err)
	})
}

func TestRecordingRepo_GetRecordingById(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	repo := recordingrepo.NewRecordingRepo(db)
	ctx := context.Background()

	t.Run("Ok get recording", func(t *testing.T) {
		id, err := testhelpers.CreateRecording(repo)
		require.NoError(t, err)

		got, err := repo.GetRecordingById(ctx, id.String())
		require.NoError(t, err)
		assert.Equal(t, id, got.Id)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetRecordingById(ctx, uuid.New().String())
		assert.Error(t, err)
	})

	t.Run("bad uuid", func(t *testing.T) {
		_, err := repo.GetRecordingById(ctx, "BadUuid")
		assert.Error(t, err)
	})
}
