package speciessqlite

import (
	"context"
	"database/sql"
	"errors"
	"uuid"
	"workshop/internal/domain/species"
	"workshop/internal/infra/db/sqlc"
	"workshop/internal/infra/db/sqlite"
)

// Store persists species using sqlc-generated SQLite queries.
type Store struct {
	queries db.Querier
}

var _ species.Store = (*Store)(nil)

// New constructs a species SQLite store.
func New(queries db.Querier) *Store { return &Store{queries: queries} }

func (s Store) Find(ctx context.Context, id uuid.UUID) (species.Species, error) {
	row, err := s.queries.SpeciesFind(ctx, id)
	if errors.Is(err, sql.ErrNoRows) {
		return species.Species{}, species.ErrNotFound
	} else if err != nil {
		return species.Species{}, err
	}

	return transform(row), nil
}

func (s Store) FindByGbifKey(ctx context.Context, gbifKey int) (species.Species, error) {
	row, err := s.queries.SpeciesFindByGbifKey(ctx, int64(gbifKey))
	if errors.Is(err, sql.ErrNoRows) {
		return species.Species{}, species.ErrNotFound
	} else if err != nil {
		return species.Species{}, err
	}

	return transform(row), nil
}

func (s Store) List(ctx context.Context) ([]species.Species, error) {
	rows, err := s.queries.SpeciesList(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]species.Species, len(rows))
	for i, row := range rows {
		result[i] = transform(row)
	}

	return result, nil
}

func (s Store) Create(ctx context.Context, in species.CreateParams) (species.Species, error) {
	row, err := s.queries.SpeciesCreate(ctx, toSpeciesCreateParams(in))
	if sqlite.IsUniqueViolation(err) {
		return species.Species{}, species.ErrAlreadyExists
	} else if err != nil {
		return species.Species{}, err
	}

	return transform(row), nil
}
