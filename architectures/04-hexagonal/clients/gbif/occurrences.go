package gbif

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

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

// Occurrences obtains a collection of occurrences.
func (g *Client) Occurrences(ctx context.Context) ([]Occurrence, error) {
	url := fmt.Sprintf("%s/occurrence/search?limit=%d&hasCoordinate=true", BaseURL, 100)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := g.HTTP.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gbif: searching occurrences: %s", res.Status)
	}

	var body struct {
		Results []Occurrence `json:"results"`
	}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Results, nil
}
