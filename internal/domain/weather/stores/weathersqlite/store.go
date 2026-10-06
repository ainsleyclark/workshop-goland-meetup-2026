package weathersqlite

import (
	"context"
	"errors"
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
	// TODO: Implement
	return weather.Weather{}, errors.New("weathersqlite: Create not implemented")
}
