package species_test

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"
	"time"
	"uuid"
	"workshop/internal/common/logs"
	"workshop/internal/domain/species"
	"workshop/internal/domain/species/stores/speciessqlite"
	"workshop/internal/infra/db/dbtest"
	db "workshop/internal/infra/db/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var errInternal = errors.New("internal error")

func setup(t *testing.T, fail bool) (context.Context, *species.Service) {
	t.Helper()

	ctx, conn, teardown := dbtest.Setup(t)
	t.Cleanup(teardown)

	q := db.Querier(db.New(conn))
	if fail {
		q = dbtest.FailingQuerier(errInternal)
	}

	return ctx, species.NewService(
		logs.NewProduction(io.Discard, slog.LevelError),
		speciessqlite.New(q),
	)
}

// newCreateParams returns a fully populated species.
func newCreateParams(gbifKey int) species.CreateParams {
	return species.CreateParams{
		GBIFKey:        gbifKey,
		Kingdom:        "Animalia",
		Phylum:         "Chordata",
		Class:          "Mammalia",
		Order:          "Carnivora",
		Family:         "Felidae",
		Genus:          "Panthera",
		Species:        "Panthera leo",
		Rank:           "SPECIES",
		ScientificName: "Panthera leo (Linnaeus, 1758)",
		CanonicalName:  "Panthera leo",
		VernacularName: "lion",
	}
}

func create(t *testing.T, ctx context.Context, svc *species.Service, in species.CreateParams) species.Species {
	t.Helper()

	got, err := svc.Create(ctx, in)
	require.NoError(t, err)
	require.NotEqual(t, uuid.Nil(), got.ID)
	assert.Equal(t, in, got.CreateParams())

	return got
}

func TestService_Find(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		ctx, svc := setup(t, false)
		lion := create(t, ctx, svc, newCreateParams(5219404))

		got, err := svc.Find(ctx, lion.ID)
		require.NoError(t, err)
		assert.Equal(t, lion, got)
	})

	t.Run("Not found", func(t *testing.T) {
		ctx, svc := setup(t, false)

		_, err := svc.Find(ctx, uuid.New())
		assert.ErrorIs(t, err, species.ErrNotFound)
	})

	t.Run("Store error", func(t *testing.T) {
		ctx, svc := setup(t, true)

		_, err := svc.Find(ctx, uuid.New())
		assert.ErrorIs(t, err, errInternal)
	})
}

func TestService_FindByGbifKey(t *testing.T) {
	t.Run("Found", func(t *testing.T) {
		ctx, svc := setup(t, false)
		lion := create(t, ctx, svc, newCreateParams(5219404))

		got, err := svc.FindByGbifKey(ctx, 5219404)
		require.NoError(t, err)
		assert.Equal(t, lion, got)
	})

	t.Run("Not found", func(t *testing.T) {
		ctx, svc := setup(t, false)

		_, err := svc.FindByGbifKey(ctx, 999)
		assert.ErrorIs(t, err, species.ErrNotFound)
	})

	t.Run("Store error", func(t *testing.T) {
		ctx, svc := setup(t, true)

		_, err := svc.FindByGbifKey(ctx, 5219404)
		assert.ErrorIs(t, err, errInternal)
	})
}

func TestService_List(t *testing.T) {
	t.Run("Empty", func(t *testing.T) {
		ctx, svc := setup(t, false)

		got, err := svc.List(ctx)
		require.NoError(t, err)
		assert.Empty(t, got)
	})

	t.Run("Lists every species", func(t *testing.T) {
		ctx, svc := setup(t, false)
		lion := create(t, ctx, svc, newCreateParams(5219404))
		tiger := create(t, ctx, svc, newCreateParams(5219448))

		got, err := svc.List(ctx)
		require.NoError(t, err)
		assert.ElementsMatch(t, []species.Species{lion, tiger}, got)
	})

	t.Run("Store error", func(t *testing.T) {
		ctx, svc := setup(t, true)

		_, err := svc.List(ctx)
		assert.ErrorIs(t, err, errInternal)
	})
}

func TestService_Create(t *testing.T) {
	t.Run("Creates species", func(t *testing.T) {
		ctx, svc := setup(t, false)
		in := newCreateParams(5219404)

		got, err := svc.Create(ctx, in)
		require.NoError(t, err)
		assert.Equal(t, in, got.CreateParams())
		assert.NotEqual(t, uuid.Nil(), got.ID)
		assert.WithinDuration(t, time.Now(), got.CreatedAt, 5*time.Second)
		assert.Equal(t, got.CreatedAt, got.UpdatedAt)
	})

	t.Run("Is retrievable once created", func(t *testing.T) {
		ctx, svc := setup(t, false)
		created := create(t, ctx, svc, newCreateParams(2435099))

		got, err := svc.Find(ctx, created.ID)
		require.NoError(t, err)
		assert.Equal(t, created, got)
	})

	t.Run("Already exists", func(t *testing.T) {
		ctx, svc := setup(t, false)
		in := newCreateParams(9999)
		create(t, ctx, svc, in)

		_, err := svc.Create(ctx, in)
		assert.ErrorIs(t, err, species.ErrAlreadyExists)
	})

	t.Run("Store error", func(t *testing.T) {
		ctx, svc := setup(t, true)

		_, err := svc.Create(ctx, newCreateParams(5219404))
		assert.ErrorIs(t, err, errInternal)
	})
}
