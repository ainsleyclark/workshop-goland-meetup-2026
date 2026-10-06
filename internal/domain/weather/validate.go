package weather

import (
	"errors"
	"fmt"
)

// Validate checks the params against the rules every weather record
// must follow. Every broken rule is returned, each wrapping ErrInvalid.
func (p CreateParams) Validate() error {
	var errs []error

	if p.ObservedAt.IsZero() {
		errs = append(errs, fmt.Errorf("%w: observed at is required", ErrInvalid))
	}
	if p.Precipitation.Total < 0 || p.Precipitation.Rain < 0 || p.Precipitation.Snowfall < 0 {
		errs = append(errs, fmt.Errorf("%w: precipitation cannot be negative", ErrInvalid))
	}
	if p.Wind.Direction < 0 || p.Wind.Direction > 360 {
		errs = append(errs, fmt.Errorf("%w: wind direction must be between 0 and 360", ErrInvalid))
	}
	if p.CloudCover < 0 || p.CloudCover > 100 {
		errs = append(errs, fmt.Errorf("%w: cloud cover must be between 0 and 100", ErrInvalid))
	}
	if p.Humidity < 0 || p.Humidity > 100 {
		errs = append(errs, fmt.Errorf("%w: humidity must be between 0 and 100", ErrInvalid))
	}

	return errors.Join(errs...)
}
