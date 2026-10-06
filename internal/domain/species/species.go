package species

import (
	"context"
	"errors"
	"log/slog"
	"uuid"
)

// Store defines the operations used to persist and query
// species data.
type Store interface {
	Find(context.Context, uuid.UUID) (Species, error)
	FindByGbifKey(context.Context, int) (Species, error)
	List(context.Context) ([]Species, error)
	Create(context.Context, CreateParams) (Species, error)
}

// Service is the entry point for working with species data.
// Handlers and other domains should use this type instead
// of interacting with the store directly.
type Service struct {
	logger *slog.Logger
	store  Store
}

// NewService returns a Service backed by the given repository.
func NewService(
	logger *slog.Logger,
	store Store,
) *Service {
	return &Service{
		logger: logger,
		store:  store,
	}
}

var (
	// ErrNotFound represents an error where a species is
	// not found in the data store.
	ErrNotFound = errors.New("species not found")

	// ErrAlreadyExists represents an error where a species with
	// the same GBIF key is already in the data store.
	ErrAlreadyExists = errors.New("species already exists")

	// ErrInvalid represents an error where a species breaks
	// one of the rules in Validate.
	ErrInvalid = errors.New("invalid species")
)

// Find returns the species with the given ID, or ErrNotFound.
func (s Service) Find(ctx context.Context, u uuid.UUID) (Species, error) {
	return s.store.Find(ctx, u)
}

// FindByGbifKey returns the species with the given GBIF key, or ErrNotFound.
func (s Service) FindByGbifKey(ctx context.Context, key int) (Species, error) {
	return s.store.FindByGbifKey(ctx, key)
}

// List returns every species in the data store.
func (s Service) List(ctx context.Context) ([]Species, error) {
	return s.store.List(ctx)
}

// Create adds a new species to the data store.
func (s Service) Create(ctx context.Context, params CreateParams) (Species, error) {
	return s.store.Create(ctx, params)
}
