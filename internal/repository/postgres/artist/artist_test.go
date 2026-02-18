package artistrepo_test

import (
	"context"
	artistrepo "spoti/internal/repository/postgres/artist"
	"spoti/internal/repository/postgres/testhelpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtistRepo_CreateArtist(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	repo := artistrepo.NewArtistRepo(db)

	t.Run("create artist", func(t *testing.T) {
		id := testhelpers.CreateTestArtist(t, repo)
		assert.NotEqual(t, uuid.Nil, id)
	})

	t.Run("duplicate uri", func(t *testing.T) {
		ctx := context.Background()
		fake := testhelpers.GetFakeArtist()

		_, err := repo.CreateArtist(ctx, fake)
		require.NoError(t, err)

		_, err = repo.CreateArtist(ctx, fake)
		assert.Error(t, err)
	})
}

func TestArtistRepo_GetArtist(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	repo := artistrepo.NewArtistRepo(db)
	ctx := context.Background()

	t.Run("existing artist", func(t *testing.T) {
		id := testhelpers.CreateTestArtist(t, repo)

		got, err := repo.GetArtist(ctx, id.String())

		require.NoError(t, err)
		assert.Equal(t, id, got.Id)
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetArtist(ctx, uuid.New().String())
		assert.Error(t, err)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		_, err := repo.GetArtist(ctx, "not-a-uuid")
		assert.Error(t, err)
	})
}

func TestArtistRepo_GetArtistsByIDs(t *testing.T) {
	db := testhelpers.SetupContainerDB()
	repo := artistrepo.NewArtistRepo(db)
	ctx := context.Background()

	t.Run("success multiple artists", func(t *testing.T) {
		id1 := testhelpers.CreateTestArtist(t, repo)
		id2 := testhelpers.CreateTestArtist(t, repo)

		got, err := repo.GetArtistsByIDs(ctx, []string{id1.String(), id2.String()})

		require.NoError(t, err)

		gotIDs := make([]uuid.UUID, len(got))
		for i, a := range got {
			gotIDs[i] = a.Id
		}
		assert.Contains(t, gotIDs, id1)
		assert.Contains(t, gotIDs, id2)
	})

	t.Run("one valid one missing", func(t *testing.T) {
		id := testhelpers.CreateTestArtist(t, repo)
		missingId := uuid.New()

		got, err := repo.GetArtistsByIDs(ctx, []string{id.String(), missingId.String()})

		require.NoError(t, err)

		gotIDs := make([]uuid.UUID, len(got))
		for i, a := range got {
			gotIDs[i] = a.Id
		}
		assert.Contains(t, gotIDs, id)
		assert.NotContains(t, gotIDs, missingId)
		assert.NotContains(t, gotIDs, uuid.Nil)
	})

	t.Run("empty ids returns empty", func(t *testing.T) {
		got, err := repo.GetArtistsByIDs(ctx, []string{})

		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("invalid uuid returns error", func(t *testing.T) {
		_, err := repo.GetArtistsByIDs(ctx, []string{"not-a-uuid"})
		assert.Error(t, err)
	})
}
