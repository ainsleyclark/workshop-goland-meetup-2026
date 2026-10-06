package species

// Species represents a taxonomy record, it can be an animal
// or plant with various data points.
type Species struct {
	Key            int
	Kingdom        string
	Phylum         string
	Class          string
	Order          string
	Family         string
	Genus          string
	Species        string
	Rank           string
	ScientificName string
	CanonicalName  string
	VernacularName string
}

// DisplayName returns the name to show a reader: the common name
// where GBIF has one, and the scientific name where it does not.
func (s Species) DisplayName() string {
	if s.VernacularName != "" {
		return s.VernacularName
	}
	return s.ScientificName
}
