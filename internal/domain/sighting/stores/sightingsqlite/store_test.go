package sightingsqlite

import (
	"context"
	"errors"
	"testing"
	"time"
	"uuid"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/sighting/country"
	"workshop/internal/domain/species"
	"workshop/internal/domain/species/stores/speciessqlite"
	"workshop/internal/domain/weather"
	"workshop/internal/domain/weather/stores/weathersqlite"
	"workshop/internal/infra/db/dbtest"
	db "workshop/internal/infra/db/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSighting_Store(t *testing.T) {
	ctx, conn, teardown := dbtest.Setup(t)
	defer teardown()

	queries := db.New(conn)
	store := New(queries)
	lion := createLion(t, ctx, queries)
	conditions := createWeather(t, ctx, queries)

	errInternal := errors.New("internal error")
	failingStore := New(dbtest.FailingQuerier(errInternal))

	setup := func(t *testing.T, in sighting.CreateParams) sighting.Sighting {
		t.Helper()

		id := uuid.New()
		got, err := store.Create(ctx, id, in)
		require.NoError(t, err)
		assert.WithinDuration(t, time.Now(), got.CreatedAt, 5*time.Second)
		assert.Equal(t, wantCreated(id, in, got.CreatedAt), got)

		got.Species = lion
		got.Weather = conditions
		return got
	}

	var fixture, nullableFixture sighting.Sighting

	t.Run("Create", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			fixture = setup(t, newCreateParams(lion.ID, conditions.ID, 123))
		})

		t.Run("Persists nil optional values", func(t *testing.T) {
			in := newCreateParams(lion.ID, conditions.ID, 124)
			in.IndividualCount = nil
			in.Location.CoordinateUncertaintyInMetres = nil
			in.Location.ElevationInMetres = nil
			in.Media = nil

			nullableFixture = setup(t, in)

			var storedMedia string
			require.NoError(t, conn.QueryRowContext(ctx, "SELECT media FROM sightings WHERE id = ?", nullableFixture.ID).Scan(&storedMedia))
			assert.Equal(t, "[]", storedMedia)
		})

		t.Run("Round-trips empty media", func(t *testing.T) {
			in := newCreateParams(lion.ID, conditions.ID, 125)
			in.Media = []sighting.Media{}

			created := setup(t, in)
			t.Cleanup(func() {
				_, err := conn.ExecContext(ctx, "DELETE FROM sightings WHERE id = ?", created.ID)
				assert.NoError(t, err)
			})

			got, err := store.Find(ctx, created.ID)
			require.NoError(t, err)
			assert.Equal(t, []sighting.Media{}, got.Media)
		})

		t.Run("Rejects duplicate GBIF key", func(t *testing.T) {
			_, err := store.Create(ctx, uuid.New(), newCreateParams(lion.ID, conditions.ID, fixture.GBIFKey))
			assert.ErrorIs(t, err, sighting.ErrAlreadyExists)
		})

		t.Run("Rejects media that cannot be encoded", func(t *testing.T) {
			in := newCreateParams(lion.ID, conditions.ID, 126)
			in.Media = []sighting.Media{{Created: time.Date(10000, 1, 1, 0, 0, 0, 0, time.UTC)}}

			id := uuid.New()
			_, err := store.Create(ctx, id, in)
			require.ErrorContains(t, err, "encode media")

			_, err = store.Find(ctx, id)
			assert.ErrorIs(t, err, sighting.ErrNotFound)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.Create(ctx, uuid.New(), newCreateParams(lion.ID, conditions.ID, 127))
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
			assert.ErrorIs(t, err, sighting.ErrNotFound)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.Find(ctx, fixture.ID)
			assert.ErrorIs(t, err, errInternal)
			assert.Zero(t, got)
		})
	})

	t.Run("FindByGbifKey", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			got, err := store.FindByGbifKey(ctx, int(fixture.GBIFKey))
			require.NoError(t, err)
			assert.Equal(t, fixture, got)
		})

		t.Run("Not found", func(t *testing.T) {
			_, err := store.FindByGbifKey(ctx, 999)
			assert.ErrorIs(t, err, sighting.ErrNotFound)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.FindByGbifKey(ctx, int(fixture.GBIFKey))
			assert.ErrorIs(t, err, errInternal)
			assert.Zero(t, got)
		})
	})

	t.Run("List", func(t *testing.T) {
		before := fixture.HappenedAt.Add(-time.Hour)
		after := fixture.HappenedAt.Add(time.Hour)
		all := []sighting.Sighting{fixture, nullableFixture}

		tt := map[string]struct {
			filter sighting.ListFilter
			want   []sighting.Sighting
		}{
			"No filter returns all sightings": {
				filter: sighting.ListFilter{},
				want:   all,
			},
			"Filters by date range and country": {
				filter: sighting.ListFilter{
					FromDate: &before,
					ToDate:   &after,
					Country:  fixture.Location.CountryCode,
				},
				want: all,
			},
			"Filters out sightings before the start date": {
				filter: sighting.ListFilter{FromDate: &after},
			},
			"Filters out sightings after the end date": {
				filter: sighting.ListFilter{ToDate: &before},
			},
			"Filters out sightings from another country": {
				filter: sighting.ListFilter{Country: country.France},
			},
			"Filters by animal": {
				filter: sighting.ListFilter{Animals: []species.Animal{species.AnimalLion}},
				want:   all,
			},
			"Filters out sightings of another animal": {
				filter: sighting.ListFilter{Animals: []species.Animal{species.AnimalBlueWhale}},
			},
			"Filters by any of several animals": {
				filter: sighting.ListFilter{Animals: []species.Animal{species.AnimalLion, species.AnimalBlueWhale}},
				want:   all,
			},
		}

		for name, test := range tt {
			t.Run(name, func(t *testing.T) {
				got, err := store.List(ctx, test.filter)
				require.NoError(t, err)
				assert.ElementsMatch(t, test.want, got)
			})
		}

		t.Run("Limits the number of sightings", func(t *testing.T) {
			// Several animals and a limit together is the case sqlc.slice
			// got wrong, binding a taxon key as the limit.
			got, err := store.List(ctx, sighting.ListFilter{
				Animals: []species.Animal{species.AnimalLion, species.AnimalBlueWhale},
				Limit:   1,
			})
			require.NoError(t, err)
			require.Len(t, got, 1)
			assert.Contains(t, all, got[0])
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.List(ctx, sighting.ListFilter{})
			assert.ErrorIs(t, err, errInternal)
			assert.Nil(t, got)
		})
	})
}

func createLion(t *testing.T, ctx context.Context, queries db.Querier) species.Species {
	t.Helper()

	lion, err := speciessqlite.New(queries).Create(ctx, species.CreateParams{
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
	})
	require.NoError(t, err)
	return lion
}

// newCreateParams returns a fully populated sighting of the given species.
func newCreateParams(speciesID, weatherID uuid.UUID, gbifKey int64) sighting.CreateParams {
	return sighting.CreateParams{
		GBIFKey:         gbifKey,
		OccurrenceID:    "occurrence-123",
		HappenedAt:      time.Date(2026, time.September, 22, 11, 12, 13, 0, time.UTC),
		BasisOfRecord:   "HUMAN_OBSERVATION",
		RecordedBy:      "observer",
		IndividualCount: new(2),
		Remarks:         "remarks",
		ReferenceURL:    "https://example.com/sighting/123",
		Location: sighting.Location{
			Locality:                      "London",
			StateProvince:                 "England",
			Country:                       "United Kingdom",
			CountryCode:                   "GB",
			Coordinates:                   sighting.Coordinates{Latitude: 51.5, Longitude: -0.12},
			CoordinateUncertaintyInMetres: new(3.5),
			ElevationInMetres:             new(125.25),
		},
		SpeciesID: speciesID,
		WeatherID: weatherID,
		Media: []sighting.Media{
			{
				ID:           uuid.New(),
				Type:         "StillImage",
				Format:       "image/jpeg",
				URL:          "https://example.com/lion.jpg",
				Title:        "Lion",
				Description:  "A lion in London",
				Created:      time.Date(2026, time.September, 22, 11, 0, 0, 0, time.UTC),
				Creator:      "Observer",
				License:      "CC-BY-4.0",
				RightsHolder: "Photographer",
			},
			{
				ID:     uuid.New(),
				Type:   "Sound",
				Format: "audio/mpeg",
				URL:    "https://example.com/lion.mp3",
			},
		},
	}
}

// wantCreated is the sighting Create should return for in.
// INSERT ... RETURNING cannot join, so the species carries only its ID.
func wantCreated(id uuid.UUID, in sighting.CreateParams, createdAt time.Time) sighting.Sighting {
	media := in.Media
	if media == nil {
		media = []sighting.Media{}
	}
	return sighting.Sighting{
		ID:              id,
		GBIFKey:         in.GBIFKey,
		OccurrenceID:    in.OccurrenceID,
		HappenedAt:      in.HappenedAt,
		BasisOfRecord:   in.BasisOfRecord,
		RecordedBy:      in.RecordedBy,
		IndividualCount: in.IndividualCount,
		Remarks:         in.Remarks,
		ReferenceURL:    in.ReferenceURL,
		Location:        in.Location,
		Species:         species.Species{ID: in.SpeciesID},
		Weather:         weather.Weather{ID: in.WeatherID},
		Media:           media,
		CreatedAt:       createdAt,
		UpdatedAt:       createdAt,
	}
}

// createWeather saves one hour's conditions for test sightings to point at.
func createWeather(t *testing.T, ctx context.Context, queries db.Querier) weather.Weather {
	t.Helper()

	got, err := weathersqlite.New(queries).Create(ctx, uuid.New(), weather.CreateParams{
		ObservedAt:    time.Date(2026, time.September, 22, 9, 0, 0, 0, time.UTC),
		Temperature:   weather.Temperature{Actual: 26.8, Apparent: 29.5},
		Precipitation: weather.Precipitation{Total: 1.2, Rain: 1.2},
		Wind:          weather.Wind{Speed: 24.7, Gusts: 41.3, Direction: 160},
		Condition:     weather.Condition{Code: 61, Description: "Slight rain"},
		CloudCover:    68,
		Humidity:      57,
	})
	require.NoError(t, err)
	return got
}
