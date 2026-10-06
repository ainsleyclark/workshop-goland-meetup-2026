package domain

// Species represents a taxonomy record, it can be an animal
// or plant with various data points.
type Species struct {
	Key            int    `json:"key"`
	Kingdom        string `json:"kingdom"`
	Phylum         string `json:"phylum"`
	Class          string `json:"class"`
	Order          string `json:"order"`
	Family         string `json:"family"`
	Genus          string `json:"genus"`
	Species        string `json:"species"`
	Rank           string `json:"rank"`
	ScientificName string `json:"scientificName"`
	CanonicalName  string `json:"canonicalName"`
	VernacularName string `json:"vernacularName"`
}

// SpeciesStore represents all possible actions for species.
type SpeciesStore interface {
	Save(s Species) error
	Find(key int) (Species, error)
}
