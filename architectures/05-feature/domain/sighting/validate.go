package sighting

import (
	"errors"
	"fmt"
)

// Validate checks the sighting against the rules every sighting must
// follow. Every broken rule is returned, each wrapping ErrInvalid.
func (s Sighting) Validate() error {
	var errs []error

	if s.Key == 0 {
		errs = append(errs, fmt.Errorf("%w: key is required", ErrInvalid))
	}
	if s.Species.Key == 0 {
		errs = append(errs, fmt.Errorf("%w: species is required", ErrInvalid))
	}
	if s.Latitude == 0 && s.Longitude == 0 {
		errs = append(errs, fmt.Errorf("%w: coordinates are required", ErrInvalid))
	}
	if s.Latitude < -90 || s.Latitude > 90 {
		errs = append(errs, fmt.Errorf("%w: latitude must be between -90 and 90", ErrInvalid))
	}
	if s.Longitude < -180 || s.Longitude > 180 {
		errs = append(errs, fmt.Errorf("%w: longitude must be between -180 and 180", ErrInvalid))
	}
	if s.EventDate == "" {
		errs = append(errs, fmt.Errorf("%w: event date is required", ErrInvalid))
	}

	return errors.Join(errs...)
}
