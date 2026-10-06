package weathersqlite

import (
	"errors"
	"testing"
	"time"
	"uuid"
	"workshop/internal/domain/weather"
	"workshop/internal/infra/db/dbtest"
	db "workshop/internal/infra/db/sqlc"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeather_Store(t *testing.T) {
	ctx, conn, teardown := dbtest.Setup(t)
	defer teardown()

	store := New(db.New(conn))

	errInternal := errors.New("internal error")
	failingStore := New(dbtest.FailingQuerier(errInternal))

	setup := func(t *testing.T, id uuid.UUID, in weather.CreateParams) weather.Weather {
		t.Helper()

		got, err := store.Create(ctx, id, in)
		require.NoError(t, err)
		assert.WithinDuration(t, time.Now(), got.CreatedAt, 5*time.Second)
		assert.Equal(t, wantCreated(id, in, got.CreatedAt), got)

		return got
	}

	var fixture weather.Weather

	t.Run("Create", func(t *testing.T) {
		t.Run("OK", func(t *testing.T) {
			fixture = setup(t, uuid.New(), newCreateParams())
		})

		t.Run("Persists the condition and the hour it was observed", func(t *testing.T) {
			in := newCreateParams()
			created := setup(t, uuid.New(), in)

			got, err := store.Find(ctx, created.ID)
			require.NoError(t, err)
			assert.Equal(t, in.Condition, got.Condition)
			assert.Equal(t, observedAt, got.ObservedAt)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.Create(ctx, uuid.New(), newCreateParams())
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
			assert.ErrorIs(t, err, weather.ErrNotFound)
		})

		t.Run("Internal error", func(t *testing.T) {
			got, err := failingStore.Find(ctx, fixture.ID)
			assert.ErrorIs(t, err, errInternal)
			assert.Zero(t, got)
		})
	})
}

// observedAt is the hour every test reading was taken in.
var observedAt = time.Date(2026, time.September, 22, 9, 0, 0, 0, time.UTC)

// newCreateParams returns a fully populated hour's weather.
func newCreateParams() weather.CreateParams {
	return weather.CreateParams{
		ObservedAt:    observedAt,
		Temperature:   weather.Temperature{Actual: 26.8, Apparent: 29.5},
		Precipitation: weather.Precipitation{Total: 1.2, Rain: 1.2, Snowfall: 0},
		Wind:          weather.Wind{Speed: 24.7, Gusts: 41.3, Direction: 160},
		Condition:     weather.Condition{Code: 61, Description: "Slight rain"},
		CloudCover:    68,
		Humidity:      57,
	}
}

// wantCreated is the weather Create should return for in.
func wantCreated(id uuid.UUID, in weather.CreateParams, createdAt time.Time) weather.Weather {
	return weather.Weather{
		ID:            id,
		ObservedAt:    in.ObservedAt,
		Temperature:   in.Temperature,
		Precipitation: in.Precipitation,
		Wind:          in.Wind,
		Condition:     in.Condition,
		CloudCover:    in.CloudCover,
		Humidity:      in.Humidity,
		CreatedAt:     createdAt,
		UpdatedAt:     createdAt,
	}
}
