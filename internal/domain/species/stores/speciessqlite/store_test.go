package speciessqlite

import (
	"errors"
	"testing"
	"time"
	"uuid"
	"workshop/internal/domain/species"
	"workshop/internal/infra/db/dbtest"
	db "workshop/internal/infra/db/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSpecies_Store(t *testing.T) {
	ctx, conn, teardown := dbtest.Setup(t)
	defer teardown()

	store := New(db.New(conn))

	errInternal := errors.New("internal error")
	failingStore := New(dbtest.FailingQuerier(errInternal))

	setup := func(t *testing.T, in species.CreateParams) species.Species {
		t.Helper()

		got, err := store.Create(ctx, in)
		require.NoError(t, err)
		assert.NotEqual(t, uuid.Nil(), got.ID)
		assert.WithinDuration(t, time.Now(), got.CreatedAt, 5*time.Second)
		assert.Equal(t, wantCreated(got.ID, in, got.CreatedAt), got)

		return got
	}

	var fixture, genusFixture species.Species

	t.Run("Create", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			fixture = setup(t, newCreateParams(5219404))
		})

		t.Run("Persists a genus without species or vernacular name", func(t *testing.T) {
			// Genus-level records have no species, and GBIF only sometimes
			// knows a common name, so both columns must accept empty text.
			in := newCreateParams(2435194)
			in.Rank = "GENUS"
			in.Species = ""
			in.ScientificName = "Panthera Oken, 1816"
			in.CanonicalName = "Panthera"
			in.VernacularName = ""

			genusFixture = setup(t, in)
		})

		t.Run("Assigns a new ID to each species", func(t *testing.T) {
			assert.NotEqual(t, fixture.ID, genusFixture.ID)
		})

		t.Run("Rejects duplicate GBIF key", func(t *testing.T) {
			_, err := store.Create(ctx, newCreateParams(fixture.GBIFKey))
			assert.ErrorIs(t, err, species.ErrAlreadyExists)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.Create(ctx, newCreateParams(5219405))
			assert.ErrorIs(t, err, errInternal)
			assert.Zero(t, got)
		})
	})

	t.Run("Find", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			got, err := store.Find(ctx, fixture.ID)
			require.NoError(t, err)
			assert.Equal(t, fixture, got)
		})

		t.Run("Not found", func(t *testing.T) {
			_, err := store.Find(ctx, uuid.New())
			assert.ErrorIs(t, err, species.ErrNotFound)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.Find(ctx, fixture.ID)
			assert.ErrorIs(t, err, errInternal)
			assert.Zero(t, got)
		})
	})

	t.Run("FindByGbifKey", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			got, err := store.FindByGbifKey(ctx, fixture.GBIFKey)
			require.NoError(t, err)
			assert.Equal(t, fixture, got)
		})

		t.Run("Not found", func(t *testing.T) {
			_, err := store.FindByGbifKey(ctx, 999)
			assert.ErrorIs(t, err, species.ErrNotFound)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.FindByGbifKey(ctx, fixture.GBIFKey)
			assert.ErrorIs(t, err, errInternal)
			assert.Zero(t, got)
		})
	})

	t.Run("List", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			got, err := store.List(ctx)
			require.NoError(t, err)
			assert.ElementsMatch(t, []species.Species{fixture, genusFixture}, got)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.List(ctx)
			assert.ErrorIs(t, err, errInternal)
			assert.Nil(t, got)
		})
	})
}

// newCreateParams returns a fully populated lion with the given GBIF key.
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

// wantCreated is the species Create should return for in.
func wantCreated(id uuid.UUID, in species.CreateParams, createdAt time.Time) species.Species {
	return species.Species{
		ID:             id,
		GBIFKey:        in.GBIFKey,
		Kingdom:        in.Kingdom,
		Phylum:         in.Phylum,
		Class:          in.Class,
		Order:          in.Order,
		Family:         in.Family,
		Genus:          in.Genus,
		Species:        in.Species,
		Rank:           in.Rank,
		ScientificName: in.ScientificName,
		CanonicalName:  in.CanonicalName,
		VernacularName: in.VernacularName,
		CreatedAt:      createdAt,
		UpdatedAt:      createdAt,
	}
}
