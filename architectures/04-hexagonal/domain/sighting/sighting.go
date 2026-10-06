package sighting

import (
	"errors"

	"04-hexagonal/clients/gbif"
	"04-hexagonal/domain/species"
)

// Store defines the operations used to persist and query sighting data.
type Store interface {
	Save(s Sighting) error
	List() []Sighting
}

var (
	// ErrNotFound represents an error where a sighting is not found
	// in the data store.
	ErrNotFound = errors.New("sighting not found")

	// ErrAlreadyExists represents an error where a sighting with the
	// same GBIF key is already in the data store.
	ErrAlreadyExists = errors.New("sighting already exists")

	// ErrInvalid represents an error where a sighting breaks one of
	// the rules in Validate.
	ErrInvalid = errors.New("invalid sighting")
)

// Service is the entry point for working with sighting data. Handlers
// should use this type instead of interacting with the store directly.
type Service struct {
	store   Store
	species species.Store
	gbif    *gbif.Client
}

// NewService returns a Service backed by the given store.
func NewService(store Store, speciesStore species.Store, client *gbif.Client) *Service {
	return &Service{
		store:   store,
		species: speciesStore,
		gbif:    client,
	}
}

// List returns every sighting held in the store.
func (s Service) List() []Sighting {
	return s.store.List()
}

// Save validates the sighting before handing it to the store.
func (s Service) Save(si Sighting) error {
	if err := si.Validate(); err != nil {
		return err
	}
	return s.store.Save(si)
}
