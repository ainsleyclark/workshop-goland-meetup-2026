package gbif

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

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

// Species fetches GET /species/{key}.
func (g *Client) Species(ctx context.Context, key int) (Species, error) {
	url := fmt.Sprintf("%s/species/%d", BaseURL, key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return Species{}, err
	}

	res, err := g.HTTP.Do(req)
	if err != nil {
		return Species{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Species{}, fmt.Errorf("gbif: fetching species %d: %s", key, res.Status)
	}

	var species Species
	if err = json.NewDecoder(res.Body).Decode(&species); err != nil {
		return Species{}, err
	}

	return species, nil
}
