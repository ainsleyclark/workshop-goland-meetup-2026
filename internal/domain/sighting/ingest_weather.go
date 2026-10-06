package sighting

import (
	"context"
	"uuid"
	"workshop/internal/clients/gbif"
)

// weatherPerSighting records the conditions where and when it happened.
func (s Service) weatherPerSighting(ctx context.Context, v gbif.Occurrence) (uuid.UUID, error) {
	// TODO: Implement
	// 1. Call OpenMeteo Client
	// 2. Save the record through to the Store.
	// 3. Return the ID.
	return uuid.Nil(), nil
}
