package speciesmem

import (
	"fmt"

	"05-feature/domain/species"
)

// Store defines an in memory, species store.
type Store struct {
	data map[int]species.Species
}

var _ species.Store = (*Store)(nil)

// New returns a new in memory, species store.
func New() *Store {
	return &Store{
		data: make(map[int]species.Species),
	}
}

// Save saves a species to the datastore.
// Returns an error if it exists by key.
func (s *Store) Save(sp species.Species) error {
	_, ok := s.data[sp.Key]
	if ok {
		return fmt.Errorf("%w: %d", species.ErrAlreadyExists, sp.Key)
	}
	s.data[sp.Key] = sp
	return nil
}

// Find returns a species by key.
// Returns ErrNotFound if it does not exist.
func (s *Store) Find(key int) (species.Species, error) {
	sp, ok := s.data[key]
	if !ok {
		return species.Species{}, fmt.Errorf("%w: %d", species.ErrNotFound, key)
	}
	return sp, nil
}
