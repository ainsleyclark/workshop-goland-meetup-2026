package models

// See: https://techdocs.gbif.org/en/openapi/.

// Occurrence is GBIF's word for one organism recorded at a place and time.
type Occurrence struct {
	Key            int     `json:"key"`              // 1258202889
	SpeciesKey     int     `json:"speciesKey"`       // absent when the ID is only genus-level
	ScientificName string  `json:"scientificName"`   // Loxodonta africana
	Country        string  `json:"country"`          // Kenya
	Locality       string  `json:"locality"`         // Maasai Mara
	Latitude       float64 `json:"decimalLatitude"`  // -1.4061
	Longitude      float64 `json:"decimalLongitude"` // 35.0083
	EventDate      string  `json:"eventDate"`        // "2016-01-11", and not a time.Time
}
