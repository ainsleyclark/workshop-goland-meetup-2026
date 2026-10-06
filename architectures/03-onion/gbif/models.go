package gbif

// See: https://techdocs.gbif.org/en/openapi/.

type (
	// Occurrence is GBIF's word for one organism recorded at a place and time.
	Occurrence struct {
		Key            int     `json:"key"`              // 1258202889
		SpeciesKey     int     `json:"speciesKey"`       // absent when the ID is only genus-level
		ScientificName string  `json:"scientificName"`   // Loxodonta africana
		Country        string  `json:"country"`          // Kenya
		Locality       string  `json:"locality"`         // Maasai Mara
		Latitude       float64 `json:"decimalLatitude"`  // -1.4061
		Longitude      float64 `json:"decimalLongitude"` // 35.0083
		EventDate      string  `json:"eventDate"`        // "2016-01-11", and not a time.Time
	}
	// Species is GBIF's taxonomy record.
	Species struct {
		Key            int    `json:"key"`            // 5231190
		Kingdom        string `json:"kingdom"`        // Animalia
		Phylum         string `json:"phylum"`         // Chordata
		Class          string `json:"class"`          // Aves
		Order          string `json:"order"`          // Passeriformes
		Family         string `json:"family"`         // Passeridae
		Genus          string `json:"genus"`          // Passer
		Species        string `json:"species"`        // Passer domesticus
		Rank           string `json:"rank"`           // SPECIES, or GENUS for a coarser ID
		ScientificName string `json:"scientificName"` // Passer domesticus (Linnaeus, 1758)
		CanonicalName  string `json:"canonicalName"`  // Passer domesticus
		VernacularName string `json:"vernacularName"` // House Sparrow - this endpoint only
	}
)
