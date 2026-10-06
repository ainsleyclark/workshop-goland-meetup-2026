package gbif

import (
	"net/url"
	"strconv"
)

// PaginatedResponse defines a response from the GBIF API
// that contains more than one record.
type PaginatedResponse[T any] struct {
	Offset       int   `json:"offset"`
	Limit        int   `json:"limit"`
	EndOfRecords bool  `json:"endOfRecords"`
	Count        int64 `json:"count"`
	Results      []T   `json:"results"`
}

// PaginationOptions controls paging for GBIF collection requests.
// Zero values omit the parameters, preserving the API defaults.
//
// See: https://techdocs.gbif.org/en/openapi/#paging
type PaginationOptions struct {
	Limit  int
	Offset int
}

func (p PaginationOptions) query() url.Values {
	q := url.Values{}
	if p.Limit != 0 {
		q.Set("limit", strconv.Itoa(p.Limit))
	}
	if p.Offset != 0 {
		q.Set("offset", strconv.Itoa(p.Offset))
	}
	return q
}
