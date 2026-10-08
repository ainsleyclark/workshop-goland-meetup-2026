package weathersqlite

import (
	"context"
	"errors"
	"math"
	"uuid"
	"workshop/internal/domain/weather"
	db "workshop/internal/infra/db/sqlc"
)

// Store persists weather using sqlc-generated SQLite queries.
type Store struct {
	queries db.Querier
}

var _ weather.Store = (*Store)(nil)

// New constructs a weather SQLite store.
func New(queries db.Querier) *Store {
	return &Store{queries: queries}
}

// Find returns the weather record for id, or weather.ErrNotFound if it
// doesn't exist.
func (s Store) Find(ctx context.Context, id uuid.UUID) (weather.Weather, error) {
	// TODO: Implement
	return weather.Weather{}, errors.New("weathersqlite: Find not implemented")
}

// Create persists a new weather record under id, or returns
// weather.ErrAlreadyExists if one is already there.
func (s Store) Create(ctx context.Context, id uuid.UUID, in weather.CreateParams) (weather.Weather, error) {
	create, err := s.queries.WeatherCreate(ctx, db.WeatherCreateParams{
		ID:                  id,
		ObservedAt:          in.ObservedAt,
		Temperature:         wholeDegrees(in.Temperature.Actual),
		ApparentTemperature: wholeDegrees(in.Temperature.Apparent),
	})
	if err != nil {
		return weather.Weather{}, err
	}
	return weather.Weather{
		ID:         create.ID,
		ObservedAt: create.ObservedAt,
		Temperature: weather.Temperature{
			Actual:   float64(deref(create.Temperature)),
			Apparent: float64(deref(create.ApparentTemperature)),
		},
		CreatedAt: create.CreatedAt,
		UpdatedAt: create.UpdatedAt,
	}, nil
}

// wholeDegrees rounds c to the nearest degree, as the temperature
// columns are INTEGER.
func wholeDegrees(c float64) *int {
	d := int(math.Round(c))
	return &d
}

// deref returns the value p points at, or the zero value if p is nil.
func deref[T any](p *T) T {
	if p == nil {
		var zero T
		return zero
	}
	return *p
}
