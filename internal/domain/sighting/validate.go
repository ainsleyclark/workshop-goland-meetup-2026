package sighting

import (
	"errors"
	"fmt"
	"time"
	"uuid"
)

// Validate checks the params against the rules every sighting must
// follow. Every broken rule is returned, each wrapping ErrInvalid.
func (p CreateParams) Validate(now time.Time) error {
	var errs []error

	if p.SpeciesID == uuid.Nil() {
		errs = append(errs, fmt.Errorf("%w: species ID is required", ErrInvalid))
	}
	if p.HappenedAt.IsZero() {
		errs = append(errs, fmt.Errorf("%w: happened at is required", ErrInvalid))
	}
	if p.HappenedAt.After(now) {
		errs = append(errs, fmt.Errorf("%w: happened at cannot be in the future", ErrInvalid))
	}
	// The coordinate rules live in ParseCoordinates, not here. A Location
	// built through it is already in range; this catches the one built by
	// hand, without stating the same bounds in two places.
	if _, err := ParseCoordinates(p.Location.Latitude, p.Location.Longitude); err != nil {
		errs = append(errs, err)
	}
	if p.IndividualCount != nil && *p.IndividualCount < 0 {
		errs = append(errs, fmt.Errorf("%w: individual count cannot be negative", ErrInvalid))
	}

	return errors.Join(errs...)
}
