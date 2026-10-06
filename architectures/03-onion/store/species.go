package store

import (
	"fmt"

	"03-onion/domain"
)

// Species defines an in memory, species store.
type Species struct {
	data map[int]domain.Species
}

var _ domain.SpeciesStore = (*Species)(nil)

// NewSpeciesStore returns a new in memory, species store.
func NewSpeciesStore() *Species {
	return &Species{
		data: make(map[int]domain.Species),
	}
}

// Save saves a species to the datastore.
// Returns an error if it exists by key.
func (s *Species) Save(sp domain.Species) error {
	_, ok := s.data[sp.Key]
	if ok {
		return fmt.Errorf("species already exists: %d", sp.Key)
	}
	s.data[sp.Key] = sp
	return nil
}

// Find returns a species by key.
// Returns an error if it does not exist.
func (s *Species) Find(key int) (domain.Species, error) {
	sp, ok := s.data[key]
	if !ok {
		return domain.Species{}, fmt.Errorf("species not found: %d", key)
	}
	return sp, nil
}
