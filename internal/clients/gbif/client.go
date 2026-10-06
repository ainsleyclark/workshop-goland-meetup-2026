package gbif

import (
	"workshop/internal/common/httputil"
)

// Client is the GBIF API client.
type Client struct {
	http *httputil.Client
}

// BaseURL is the URL of the GBIF API.
const BaseURL = "https://api.gbif.org/v1"

// New creates a new GBIF API client for the given base URL.
func New(baseURL string) *Client {
	if baseURL == "" {
		baseURL = BaseURL
	}
	return &Client{
		http: httputil.New(httputil.Config{
			BaseURL: baseURL,
		}),
	}
}
