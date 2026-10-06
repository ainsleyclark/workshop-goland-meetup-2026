package sighting

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
)

func TestCreateParams_Validate(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 9, 22, 12, 0, 0, 0, time.UTC)
	past := now.Add(-time.Hour)
	speciesID := uuid.New()

	tt := map[string]struct {
		input   CreateParams
		wantErr bool
	}{
		"Valid":                  {input: CreateParams{SpeciesID: speciesID, HappenedAt: past, Location: Location{Coordinates: Coordinates{Latitude: 51.5, Longitude: -0.12}}}, wantErr: false},
		"Boundary coordinates":   {input: CreateParams{SpeciesID: speciesID, HappenedAt: past, Location: Location{Coordinates: Coordinates{Latitude: -90, Longitude: 180}}}, wantErr: false},
		"Happened now":           {input: CreateParams{SpeciesID: speciesID, HappenedAt: now}, wantErr: false},
		"Zero individual count":  {input: CreateParams{SpeciesID: speciesID, HappenedAt: past, IndividualCount: new(0)}, wantErr: false},
		"Missing species ID":     {input: CreateParams{HappenedAt: past}, wantErr: true},
		"Missing happened at":    {input: CreateParams{SpeciesID: speciesID}, wantErr: true},
		"Happened in the future": {input: CreateParams{SpeciesID: speciesID, HappenedAt: now.Add(time.Second)}, wantErr: true},
		"Latitude too low":       {input: CreateParams{SpeciesID: speciesID, HappenedAt: past, Location: Location{Latitude: -90.1}}, wantErr: true},
		"Latitude too high":      {input: CreateParams{SpeciesID: speciesID, HappenedAt: past, Location: Location{Latitude: 90.1}}, wantErr: true},
		"Longitude too low":      {input: CreateParams{SpeciesID: speciesID, HappenedAt: past, Location: Location{Longitude: -180.1}}, wantErr: true},
		"Longitude too high":     {input: CreateParams{SpeciesID: speciesID, HappenedAt: past, Location: Location{Longitude: 180.1}}, wantErr: true},
		"Negative count":         {input: CreateParams{SpeciesID: speciesID, HappenedAt: past, IndividualCount: new(-1)}, wantErr: true},
		"Multiple rules broken":  {input: CreateParams{Location: Location{Coordinates: Coordinates{Latitude: 100, Longitude: 200}}}, wantErr: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := test.input.Validate(now)
			assert.Equal(t, test.wantErr, err != nil)
		})
	}
}
