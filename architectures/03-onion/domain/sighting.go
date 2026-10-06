package domain

// Sighting represents an individual sighting of a species.
// It's a record of a specific observation at a particular location and time.
type Sighting struct {
	Key            int     `json:"key"`            // 1258202889
	SpeciesKey     int     `json:"speciesKey"`     // absent when the ID is only genus-level
	ScientificName string  `json:"scientificName"` // Loxodonta africana
	Country        string  `json:"country"`        // Kenya
	Locality       string  `json:"locality"`       // Maasai Mara
	Latitude       float64 `json:"latitude"`       // -1.4061
	Longitude      float64 `json:"longitude"`      // 35.0083
	EventDate      string  `json:"eventDate"`      // "2016-01-11", and not a time.Time
}

// SightingStore defines all possible actions for sightings.
type SightingStore interface {
	Save(s Sighting) error
	List() []Sighting
}
