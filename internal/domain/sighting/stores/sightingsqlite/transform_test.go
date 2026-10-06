package sightingsqlite

import (
	"testing"
	"time"
	"uuid"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
	"workshop/internal/infra/db/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTransform(t *testing.T) {
	t.Parallel()

	t.Run("Maps all fields and decodes media", func(t *testing.T) {
		t.Parallel()

		id, speciesID := uuid.New(), uuid.New()
		happenedAt := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
		createdAt, updatedAt := happenedAt.Add(time.Hour), happenedAt.Add(2*time.Hour)
		input := db.Sighting{
			ID:                            id,
			GbifKey:                       123,
			OccurrenceID:                  "occurrence-1",
			HappenedAt:                    happenedAt,
			BasisOfRecord:                 "HUMAN_OBSERVATION",
			RecordedBy:                    "Observer",
			IndividualCount:               new(3),
			Remarks:                       "Near the river",
			ReferenceUrl:                  "https://example.com/observation",
			Locality:                      "Richmond",
			StateProvince:                 "London",
			Country:                       "United Kingdom",
			CountryCode:                   "GB",
			Latitude:                      51.45,
			Longitude:                     -0.3,
			CoordinateUncertaintyInMetres: new(15.5),
			ElevationInMetres:             new(120.5),
			SpeciesID:                     speciesID,
			Media:                         `[{"type":"StillImage","format":"image/jpeg","url":"https://example.com/image.jpg","title":"Otter","description":"Swimming","created":"2026-09-20T12:00:00Z","creator":"Photographer","license":"CC0","rightsHolder":"Owner"},{"type":"Sound","format":"audio/mpeg","url":"https://example.com/audio.mp3"}]`,
			CreatedAt:                     createdAt,
			UpdatedAt:                     updatedAt,
		}
		sp := db.Species{
			ID:             speciesID,
			GbifKey:        5219404,
			ScientificName: "Lutra lutra",
		}
		want := sighting.Sighting{
			ID:              id,
			GBIFKey:         123,
			OccurrenceID:    "occurrence-1",
			HappenedAt:      happenedAt,
			BasisOfRecord:   "HUMAN_OBSERVATION",
			RecordedBy:      "Observer",
			IndividualCount: new(3),
			Remarks:         "Near the river",
			ReferenceURL:    "https://example.com/observation",
			Location: sighting.Location{
				Locality:                      "Richmond",
				StateProvince:                 "London",
				Country:                       "United Kingdom",
				CountryCode:                   "GB",
				Coordinates:                   sighting.Coordinates{Latitude: 51.45, Longitude: -0.3},
				CoordinateUncertaintyInMetres: new(15.5),
				ElevationInMetres:             new(120.5),
			},
			Species: species.Species{
				ID:             speciesID,
				GBIFKey:        5219404,
				ScientificName: "Lutra lutra",
			},
			Media: []sighting.Media{
				{
					Type:         "StillImage",
					Format:       "image/jpeg",
					URL:          "https://example.com/image.jpg",
					Title:        "Otter",
					Description:  "Swimming",
					Created:      happenedAt,
					Creator:      "Photographer",
					License:      "CC0",
					RightsHolder: "Owner",
				},
				{
					Type:   "Sound",
					Format: "audio/mpeg",
					URL:    "https://example.com/audio.mp3",
				},
			},
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		}

		got, err := transform(db.SightingListRow{Sighting: input, Species: sp})
		require.NoError(t, err)
		assert.Equal(t, want, got)
	})

	t.Run("Preserves zero values with empty media", func(t *testing.T) {
		t.Parallel()

		got, err := transform(db.SightingListRow{Sighting: db.Sighting{Media: "[]"}})
		require.NoError(t, err)
		assert.Equal(t, sighting.Sighting{Media: []sighting.Media{}}, got)
	})

	t.Run("Normalises null media to an empty slice", func(t *testing.T) {
		t.Parallel()

		got, err := transform(db.SightingListRow{Sighting: db.Sighting{Media: "null"}})
		require.NoError(t, err)
		assert.Equal(t, []sighting.Media{}, got.Media)
	})

	t.Run("Rejects malformed media JSON", func(t *testing.T) {
		t.Parallel()

		got, err := transform(db.SightingListRow{Sighting: db.Sighting{Media: "["}})
		require.ErrorContains(t, err, "decode media")
		assert.Equal(t, sighting.Sighting{}, got)
	})

	t.Run("Rejects empty media JSON", func(t *testing.T) {
		t.Parallel()

		got, err := transform(db.SightingListRow{})
		require.ErrorContains(t, err, "decode media")
		assert.Equal(t, sighting.Sighting{}, got)
	})

	t.Run("Rejects media that is not an array", func(t *testing.T) {
		t.Parallel()

		got, err := transform(db.SightingListRow{Sighting: db.Sighting{Media: "{}"}})
		require.ErrorContains(t, err, "decode media")
		assert.Equal(t, sighting.Sighting{}, got)
	})

	t.Run("Rejects invalid media timestamps", func(t *testing.T) {
		t.Parallel()

		got, err := transform(db.SightingListRow{Sighting: db.Sighting{Media: `[{"created":"invalid"}]`}})
		require.ErrorContains(t, err, "decode media")
		assert.Equal(t, sighting.Sighting{}, got)
	})
}

func TestToSightingCreateParams(t *testing.T) {
	t.Parallel()

	t.Run("Maps all fields and encodes media", func(t *testing.T) {
		t.Parallel()

		id, speciesID := uuid.New(), uuid.New()
		happenedAt := time.Date(2026, 9, 20, 12, 0, 0, 0, time.UTC)
		input := sighting.CreateParams{
			GBIFKey:         123,
			OccurrenceID:    "occurrence-1",
			HappenedAt:      happenedAt,
			BasisOfRecord:   "HUMAN_OBSERVATION",
			RecordedBy:      "Observer",
			IndividualCount: new(3),
			Remarks:         "Near the river",
			ReferenceURL:    "https://example.com/observation",
			SpeciesID:       speciesID,
			Location: sighting.Location{
				Locality:                      "Richmond",
				StateProvince:                 "London",
				Country:                       "United Kingdom",
				CountryCode:                   "GB",
				Coordinates:                   sighting.Coordinates{Latitude: 51.45, Longitude: -0.3},
				CoordinateUncertaintyInMetres: new(15.5),
				ElevationInMetres:             new(120.5),
			},
			Media: []sighting.Media{
				{
					Type:         "StillImage",
					Format:       "image/jpeg",
					URL:          "https://example.com/image.jpg",
					Title:        "Otter",
					Description:  "Swimming",
					Created:      happenedAt,
					Creator:      "Photographer",
					License:      "CC0",
					RightsHolder: "Owner",
				},
				{
					Type:   "Sound",
					Format: "audio/mpeg",
					URL:    "https://example.com/audio.mp3",
				},
			},
		}
		want := db.SightingCreateParams{
			ID:                            id,
			GbifKey:                       123,
			OccurrenceID:                  "occurrence-1",
			HappenedAt:                    happenedAt,
			BasisOfRecord:                 "HUMAN_OBSERVATION",
			RecordedBy:                    "Observer",
			IndividualCount:               new(3),
			Remarks:                       "Near the river",
			ReferenceUrl:                  "https://example.com/observation",
			Locality:                      "Richmond",
			StateProvince:                 "London",
			Country:                       "United Kingdom",
			CountryCode:                   "GB",
			Latitude:                      51.45,
			Longitude:                     -0.3,
			CoordinateUncertaintyInMetres: new(15.5),
			ElevationInMetres:             new(120.5),
			SpeciesID:                     speciesID,
			Media:                         `[{"id":"00000000-0000-0000-0000-000000000000","type":"StillImage","format":"image/jpeg","url":"https://example.com/image.jpg","title":"Otter","description":"Swimming","created":"2026-09-20T12:00:00Z","creator":"Photographer","license":"CC0","rightsHolder":"Owner"},{"id":"00000000-0000-0000-0000-000000000000","type":"Sound","format":"audio/mpeg","url":"https://example.com/audio.mp3"}]`,
		}

		got, err := toSightingCreateParams(id, input)
		require.NoError(t, err)
		assert.JSONEq(t, want.Media, got.Media)
		got.Media = want.Media
		assert.Equal(t, want, got)
	})

	t.Run("Preserves zero values and encodes nil media as an empty array", func(t *testing.T) {
		t.Parallel()

		got, err := toSightingCreateParams(uuid.UUID{}, sighting.CreateParams{})
		require.NoError(t, err)
		assert.Equal(t, db.SightingCreateParams{
			Media: "[]",
		}, got)
	})

	t.Run("Encodes empty media as an empty array", func(t *testing.T) {
		t.Parallel()

		got, err := toSightingCreateParams(uuid.UUID{}, sighting.CreateParams{
			Media: []sighting.Media{},
		})
		require.NoError(t, err)
		assert.Equal(t, "[]", got.Media)
	})

	t.Run("Rejects media that cannot be encoded", func(t *testing.T) {
		t.Parallel()

		input := sighting.CreateParams{
			Media: []sighting.Media{{
				Created: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC),
			}},
		}

		got, err := toSightingCreateParams(uuid.New(), input)
		require.ErrorContains(t, err, "encode media")
		assert.Equal(t, db.SightingCreateParams{}, got)
	})
}
