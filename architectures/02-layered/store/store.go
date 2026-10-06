package store

import (
	"02-layered/models"
)

// Storage represents all possible actions available to deal with data.
type Storage interface {
	SaveOccurrence(o models.Occurrence) error
	ListOccurrences() []models.Occurrence
	SaveSpecies(sp models.Species) error
	FindSpecies(key int) (models.Species, error)
}
