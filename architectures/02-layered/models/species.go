package models

// See: https://techdocs.gbif.org/en/openapi/.

// Species is GBIF's taxonomy record.
type Species struct {
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
