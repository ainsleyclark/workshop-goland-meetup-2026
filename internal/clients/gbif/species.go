package gbif

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"workshop/internal/common/httputil"
)

// Species fetches GET /species/{key}.
//
// Obtains the taxonomic name usage record for a single species by
// its usage key.
//
// See: https://techdocs.gbif.org/en/openapi/v1/species
func (c Client) Species(ctx context.Context, key int) (Species, error) {
	return c.http.Do[Species](ctx, httputil.Request{
		Method: http.MethodGet,
		Path:   "/species/" + strconv.Itoa(key),
	})
}

// Species defines a taxonomic name usage record.
type Species struct {
	Key                 int       `json:"key"`
	NubKey              int       `json:"nubKey"`
	NameKey             int       `json:"nameKey"`
	TaxonID             string    `json:"taxonID"`
	SourceTaxonKey      int       `json:"sourceTaxonKey"`
	Kingdom             string    `json:"kingdom"`
	Phylum              string    `json:"phylum"`
	Order               string    `json:"order"`
	Family              string    `json:"family"`
	Genus               string    `json:"genus"`
	Species             string    `json:"species"`
	KingdomKey          int       `json:"kingdomKey"`
	PhylumKey           int       `json:"phylumKey"`
	ClassKey            int       `json:"classKey"`
	OrderKey            int       `json:"orderKey"`
	FamilyKey           int       `json:"familyKey"`
	GenusKey            int       `json:"genusKey"`
	SpeciesKey          int       `json:"speciesKey"`
	DatasetKey          string    `json:"datasetKey"`
	ConstituentKey      string    `json:"constituentKey"`
	ParentKey           int       `json:"parentKey"`
	Parent              string    `json:"parent"`
	BasionymKey         int       `json:"basionymKey"`
	Basionym            string    `json:"basionym"`
	ScientificName      string    `json:"scientificName"`
	CanonicalName       string    `json:"canonicalName"`
	VernacularName      string    `json:"vernacularName"`
	Authorship          string    `json:"authorship"`
	NameType            string    `json:"nameType"`
	Rank                string    `json:"rank"`
	Origin              string    `json:"origin"`
	TaxonomicStatus     string    `json:"taxonomicStatus"`
	NomenclaturalStatus []any     `json:"nomenclaturalStatus"`
	Remarks             string    `json:"remarks"`
	PublishedIn         string    `json:"publishedIn"`
	NumDescendants      int       `json:"numDescendants"`
	LastCrawled         time.Time `json:"lastCrawled"`
	LastInterpreted     time.Time `json:"lastInterpreted"`
	Issues              []any     `json:"issues"`
	Class               string    `json:"class"`
}
