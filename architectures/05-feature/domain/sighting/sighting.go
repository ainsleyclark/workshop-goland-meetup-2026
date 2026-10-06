package sighting

import (
	"errors"

	"05-feature/clients/gbif"
	"05-feature/domain/species"
)

// Store defines the operations used to persist and query sighting data.
type Store interface {
	Save(s Sighting) error
	List() []Sighting
}

// SpeciesReadWriter is what ingest needs from species: look one up,
// and save it if it is new.
type SpeciesReadWriter interface {
	Find(key int) (species.Species, error)
	Save(sp species.Species) error
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
	species SpeciesReadWriter
	gbif    *gbif.Client
}

// NewService returns a Service backed by the given store, the species
// reader it needs to name what was seen, and a GBIF client to ingest from.
func NewService(store Store, speciesSvc SpeciesReadWriter, client *gbif.Client) *Service {
	return &Service{
		store:   store,
		species: speciesSvc,
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
