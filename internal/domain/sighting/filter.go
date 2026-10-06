package sighting

import (
	"time"
	"workshop/internal/domain/sighting/country"
	"workshop/internal/domain/species"
)

// ListFilter narrows a listing. The zero value matches everything, so
// callers only set the fields they care about.
type ListFilter struct {
	FromDate *time.Time
	ToDate   *time.Time
	Country  country.Code
	Animals  []species.Animal
	// Limit caps how many sightings come back, newest first. Zero means
	// no limit.
	Limit int
}

// TaxonKeys returns the GBIF taxon keys the filter's animals stand for,
// which is what a species is stored against.
func (f ListFilter) TaxonKeys() []int {
	return taxonKeys(f.Animals)
}

// taxonKeys retrieves the integers values according to the GBIF client.
func taxonKeys(animals []species.Animal) []int {
	var keys []int
	seen := make(map[int]bool)
	for _, animal := range animals {
		for _, key := range animal.TaxonKeys() {
			if !seen[key] {
				seen[key] = true
				keys = append(keys, key)
			}
		}
	}
	return keys
}
