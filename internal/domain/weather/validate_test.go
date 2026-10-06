package weather

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateParams_Validate(t *testing.T) {
	t.Parallel()

	valid := CreateParams{
		ObservedAt:    time.Date(2026, 9, 22, 14, 0, 0, 0, time.UTC),
		Temperature:   Temperature{Actual: 24.2, Apparent: 26.0},
		Precipitation: Precipitation{Total: 2.4, Rain: 2.4},
		Wind:          Wind{Speed: 21.3, Gusts: 38.9, Direction: 135},
		Condition:     Condition{Code: 61, Description: "Slight rain"},
		CloudCover:    72,
		Humidity:      64,
	}

	tt := map[string]struct {
		mutate   func(p *CreateParams)
		wantErrs []string
	}{
		"Valid": {
			mutate: func(*CreateParams) {},
		},
		"Missing observed at": {
			mutate:   func(p *CreateParams) { p.ObservedAt = time.Time{} },
			wantErrs: []string{"observed at is required"},
		},
		"Negative precipitation": {
			mutate:   func(p *CreateParams) { p.Precipitation.Total = -1 },
			wantErrs: []string{"precipitation cannot be negative"},
		},
		"Negative rain": {
			mutate:   func(p *CreateParams) { p.Precipitation.Rain = -0.1 },
			wantErrs: []string{"precipitation cannot be negative"},
		},
		"Negative snowfall": {
			mutate:   func(p *CreateParams) { p.Precipitation.Snowfall = -0.1 },
			wantErrs: []string{"precipitation cannot be negative"},
		},
		"Wind direction below range": {
			mutate:   func(p *CreateParams) { p.Wind.Direction = -1 },
			wantErrs: []string{"wind direction must be between 0 and 360"},
		},
		"Wind direction above range": {
			mutate:   func(p *CreateParams) { p.Wind.Direction = 361 },
			wantErrs: []string{"wind direction must be between 0 and 360"},
		},
		"Cloud cover out of range": {
			mutate:   func(p *CreateParams) { p.CloudCover = 101 },
			wantErrs: []string{"cloud cover must be between 0 and 100"},
		},
		"Humidity out of range": {
			mutate:   func(p *CreateParams) { p.Humidity = -1 },
			wantErrs: []string{"humidity must be between 0 and 100"},
		},
		"Reports every broken rule": {
			mutate: func(p *CreateParams) {
				p.ObservedAt = time.Time{}
				p.Wind.Direction = 400
			},
			wantErrs: []string{
				"observed at is required",
				"wind direction must be between 0 and 360",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			in := valid
			test.mutate(&in)

			err := in.Validate()
			if len(test.wantErrs) == 0 {
				assert.NoError(t, err)
				return
			}

			require.ErrorIs(t, err, ErrInvalid)
			for _, want := range test.wantErrs {
				assert.ErrorContains(t, err, want)
			}
		})
	}
}
