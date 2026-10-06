package openmeteo

import "errors"

// Coordinates represents geographical coordinates with latitude and longitude.
type Coordinates struct {
	Latitude  float64 // Latitude in degrees, must be between -90 and 90.
	Longitude float64 // Longitude in degrees, must be between -180 and 180.
}

// Validate checks whether the latitude and longitude values are within valid ranges.
// Returns an error if either the latitude or longitude is out of bounds.
func (c Coordinates) Validate() error {
	if c.Latitude < -90 || c.Latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if c.Longitude < -180 || c.Longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}
