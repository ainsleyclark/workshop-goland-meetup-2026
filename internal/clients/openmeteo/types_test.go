package openmeteo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoordinates_Validate(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   Coordinates
		wantErr bool
	}{
		"Valid Coordinates": {
			input: Coordinates{
				Latitude:  45.0,
				Longitude: 90.0,
			},
			wantErr: false,
		},
		"Latitude Out Of Bounds Negative": {
			input: Coordinates{
				Latitude:  -95.0,
				Longitude: 45.0,
			},
			wantErr: true,
		},
		"Latitude Out Of Bounds Positive": {
			input: Coordinates{
				Latitude:  100.0,
				Longitude: 45.0,
			},
			wantErr: true,
		},
		"Longitude Out Of Bounds Negative": {
			input: Coordinates{
				Latitude:  45.0,
				Longitude: -190.0,
			},
			wantErr: true,
		},
		"Longitude Out Of Bounds Positive": {
			input: Coordinates{
				Latitude:  45.0,
				Longitude: 200.0,
			},
			wantErr: true,
		},
		"Both Latitude and Longitude Out Of Bounds": {
			input: Coordinates{
				Latitude:  95.0,
				Longitude: 200.0,
			},
			wantErr: true,
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.input.Validate()
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
