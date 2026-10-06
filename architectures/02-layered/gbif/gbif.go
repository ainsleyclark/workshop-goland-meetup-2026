package gbif

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"02-layered/models"
)

// BaseURL is where the GBIF API lives.
const BaseURL = "https://api.gbif.org/v1"

// Client represents the type that interacts with the GBIF API.
type Client struct {
	HTTP *http.Client
}

// Occurrences obtains a collection of occurrences.
func (g *Client) Occurrences(ctx context.Context) ([]models.Occurrence, error) {
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
		Results []models.Occurrence `json:"results"`
	}
	if err = json.NewDecoder(res.Body).Decode(&body); err != nil {
		return nil, err
	}

	return body.Results, nil
}

// Species fetches GET /species/{key}.
func (g *Client) Species(ctx context.Context, key int) (models.Species, error) {
	url := fmt.Sprintf("%s/species/%d", BaseURL, key)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return models.Species{}, err
	}

	res, err := g.HTTP.Do(req)
	if err != nil {
		return models.Species{}, err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return models.Species{}, fmt.Errorf("gbif: fetching species %d: %s", key, res.Status)
	}

	var species models.Species
	if err = json.NewDecoder(res.Body).Decode(&species); err != nil {
		return models.Species{}, err
	}

	return species, nil
}
