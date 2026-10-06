package gbif

import "net/http"

// BaseURL is where the GBIF API lives.
const BaseURL = "https://api.gbif.org/v1"

// Client represents the type that interacts with the GBIF API.
type Client struct {
	HTTP *http.Client
}

// New returns a Client using the default HTTP client.
func New() *Client {
	return &Client{HTTP: http.DefaultClient}
}
