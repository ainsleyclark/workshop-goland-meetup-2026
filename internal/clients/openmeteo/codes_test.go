package openmeteo

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWeatherCode_String(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input WeatherCode
		want  string
	}{
		"OK": {
			input: 0,
			want:  "Clear sky",
		},
		"Unknown": {
			input: -1,
			want:  "Unknown weather code",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := test.input.String()
			assert.Equal(t, test.want, got)
		})
	}
}
