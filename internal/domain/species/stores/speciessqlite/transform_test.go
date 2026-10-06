package speciessqlite

import (
	"testing"
	"time"
	"uuid"
	"workshop/internal/domain/species"
	"workshop/internal/infra/db/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransform(t *testing.T) {
	t.Parallel()

	t.Run("Maps all fields", func(t *testing.T) {
		t.Parallel()

		id := uuid.New()
		createdAt := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
		updatedAt := createdAt.Add(time.Hour)
		input := db.Species{
			ID:             id,
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
			GbifKey:        5219404,
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
		want := species.Species{
			ID:             id,
			GBIFKey:        5219404,
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
			CreatedAt:      createdAt,
			UpdatedAt:      updatedAt,
		}

		assert.Equal(t, want, transform(input))
	})

	t.Run("Preserves zero values", func(t *testing.T) {
		t.Parallel()

		assert.Equal(t, species.Species{}, transform(db.Species{}))
	})
}

func TestToSpeciesCreateParams(t *testing.T) {
	t.Parallel()

	t.Run("Maps all fields", func(t *testing.T) {
		t.Parallel()

		input := species.CreateParams{
			GBIFKey:        5219404,
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
		want := db.SpeciesCreateParams{
			GbifKey:        5219404,
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

		got := toSpeciesCreateParams(input)
		require.NotEqual(t, uuid.Nil(), got.ID)
		// The ID is random, so compare everything else.
		got.ID = uuid.Nil()
		assert.Equal(t, want, got)
	})

	t.Run("Preserves zero values but still assigns an ID", func(t *testing.T) {
		t.Parallel()

		got := toSpeciesCreateParams(species.CreateParams{})
		require.NotEqual(t, uuid.Nil(), got.ID)
		got.ID = uuid.Nil()
		assert.Equal(t, db.SpeciesCreateParams{}, got)
	})

	t.Run("Assigns a new ID on each call", func(t *testing.T) {
		t.Parallel()

		in := species.CreateParams{GBIFKey: 5219404}
		assert.NotEqual(t, toSpeciesCreateParams(in).ID, toSpeciesCreateParams(in).ID)
	})
}
