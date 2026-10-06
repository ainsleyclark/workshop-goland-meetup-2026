package weathersqlite

import (
	"context"
	"database/sql"
	"errors"
	"uuid"
	"workshop/internal/domain/weather"
	"workshop/internal/infra/db/sqlc"
	"workshop/internal/infra/db/sqlite"
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
	row, err := s.queries.WeatherFind(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return weather.Weather{}, weather.ErrNotFound
	} else if err != nil {
		return weather.Weather{}, err
	}

	return transform(row), nil
}

// Create persists a new weather record under id, or returns
// weather.ErrAlreadyExists if one is already there.
func (s Store) Create(ctx context.Context, id uuid.UUID, in weather.CreateParams) (weather.Weather, error) {
	row, err := s.queries.WeatherCreate(ctx, toWeatherCreateParams(id, in))
	if sqlite.IsUniqueViolation(err) {
		return weather.Weather{}, weather.ErrAlreadyExists
	} else if err != nil {
		return weather.Weather{}, err
	}

	return transform(row), nil
}
