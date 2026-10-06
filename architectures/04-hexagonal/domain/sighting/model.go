package sighting

// Sighting represents an individual sighting of a species.
// It's a record of a specific observation at a particular location and time.
type Sighting struct {
	Key            int     // 1258202889
	SpeciesKey     int     // absent when the ID is only genus-level
	ScientificName string  // Loxodonta africana
	Country        string  // Kenya
	Locality       string  // Maasai Mara
	Latitude       float64 // -1.4061
	Longitude      float64 // 35.0083
	EventDate      string  // "2016-01-11", and not a time.Time
}
