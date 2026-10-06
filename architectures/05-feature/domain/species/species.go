package species

import "errors"

// Store defines the operations used to persist and query species data.
// It is declared here, by the package that needs it, and implemented in
// stores/speciesmem.
type Store interface {
	Save(s Species) error
	Find(key int) (Species, error)
}

var (
	// ErrNotFound represents an error where a species is not found
	// in the data store.
	ErrNotFound = errors.New("species not found")

	// ErrAlreadyExists represents an error where a species with the
	// same key is already in the data store.
	ErrAlreadyExists = errors.New("species already exists")
)

// Service is the entry point for working with species data. Handlers
// and other domains should use this instead of the store directly.
type Service struct {
	store Store
}

// NewService returns a Service backed by the given store.
func NewService(store Store) *Service {
	return &Service{store: store}
}

// Find returns the species with the given key, or ErrNotFound.
func (s Service) Find(key int) (Species, error) {
	return s.store.Find(key)
}

// Save stores the species.
func (s Service) Save(sp Species) error {
	return s.store.Save(sp)
}
