package weather

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
)

func TestWeather_CreateParams(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Weather
		want  CreateParams
	}{
		"Zero value": {},
		"Timestamps excluded": {
			input: Weather{
				CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		"ID only": {
			input: Weather{
				ID: uuid.New(),
			},
			want: CreateParams{},
		},
		"All fields": {
			input: Weather{
				ID:            uuid.New(),
				ObservedAt:    time.Date(2026, 9, 22, 14, 0, 0, 0, time.UTC),
				Temperature:   Temperature{Actual: 24.2, Apparent: 26.0},
				Precipitation: Precipitation{Total: 2.4, Rain: 2.4},
				Wind:          Wind{Speed: 21.3, Gusts: 38.9, Direction: 135},
				Condition:     Condition{Code: 61, Description: "Slight rain"},
				CloudCover:    72,
				Humidity:      64,
				CreatedAt:     time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC),
				UpdatedAt:     time.Date(2026, 9, 23, 0, 0, 0, 0, time.UTC),
			},
			want: CreateParams{
				ObservedAt:    time.Date(2026, 9, 22, 14, 0, 0, 0, time.UTC),
				Temperature:   Temperature{Actual: 24.2, Apparent: 26.0},
				Precipitation: Precipitation{Total: 2.4, Rain: 2.4},
				Wind:          Wind{Speed: 21.3, Gusts: 38.9, Direction: 135},
				Condition:     Condition{Code: 61, Description: "Slight rain"},
				CloudCover:    72,
				Humidity:      64,
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, test.input.CreateParams())
		})
	}
}
