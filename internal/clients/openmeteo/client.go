package openmeteo

import (
	"workshop/internal/common/httputil"
)

// Client is the Open-Meteo Archive API client.
type Client struct {
	http *httputil.Client
}

// BaseURL is the URL of the Open-Meteo archive API.
const BaseURL = "https://archive-api.open-meteo.com/v1"

// New creates a new Open-Meteo Archive API client. An empty baseURL
// falls back to BaseURL; tests pass the address of a stub server.
func New(baseURL string) *Client {
	if baseURL == "" {
		baseURL = BaseURL
	}
	return &Client{
		http: httputil.New(httputil.Config{BaseURL: baseURL}),
	}
}
