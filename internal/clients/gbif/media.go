package gbif

import (
	"context"
	"net/http"
	"strconv"
	"time"
	"workshop/internal/common/httputil"
)

// Media fetches GET /species/{key}/media.
//
// Obtains a collection of media items (images, sounds, etc.) linked
// to a single species by its usage key.
// See: https://techdocs.gbif.org/en/openapi/images
func (c Client) Media(ctx context.Context, key int) (PaginatedResponse[Media], error) {
	return c.http.Do[PaginatedResponse[Media]](ctx, httputil.Request{
		Method: http.MethodGet,
		Path:   "/species/" + strconv.Itoa(key) + "/media",
	})
}

// Media defines a single media item (e.g. an image) linked to a species.
type Media struct {
	Type           string    `json:"type"`
	Format         string    `json:"format"`
	Source         string    `json:"source,omitempty"`
	Title          string    `json:"title,omitempty"`
	Description    string    `json:"description,omitempty"`
	Created        time.Time `json:"created"`
	Creator        string    `json:"creator,omitempty"`
	Contributor    string    `json:"contributor,omitempty"`
	Publisher      string    `json:"publisher,omitempty"`
	License        string    `json:"license,omitempty"`
	RightsHolder   string    `json:"rightsHolder,omitempty"`
	TaxonKey       int       `json:"taxonKey,omitempty"`
	SourceTaxonKey int       `json:"sourceTaxonKey,omitempty"`
	Identifier     string    `json:"identifier"`
	References     string    `json:"references,omitempty"`
}
