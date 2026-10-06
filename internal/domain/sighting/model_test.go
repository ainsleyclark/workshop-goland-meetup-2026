package sighting

import (
	"testing"
	"time"
	"uuid"
	"workshop/internal/domain/species"

	"github.com/stretchr/testify/assert"
)

func TestSighting_CreateParams(t *testing.T) {
	t.Parallel()

	nestedID := uuid.New()

	tt := map[string]struct {
		input Sighting
		want  CreateParams
	}{
		"Zero value": {},
		"Timestamps excluded": {
			input: Sighting{
				CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		"ID only": {
			input: Sighting{
				ID: uuid.New(),
			},
			want: CreateParams{},
		},
		"All fields": {
			input: Sighting{
				ID:              uuid.New(),
				GBIFKey:         123,
				OccurrenceID:    "occurrence-123",
				HappenedAt:      time.Date(2026, 9, 22, 12, 30, 0, 0, time.UTC),
				BasisOfRecord:   "HUMAN_OBSERVATION",
				RecordedBy:      "Observer",
				IndividualCount: new(2),
				Remarks:         "A robin",
				ReferenceURL:    "https://example.com/123",
				Location: Location{
					Locality:                      "London",
					StateProvince:                 "England",
					Country:                       "United Kingdom",
					CountryCode:                   "GB",
					Coordinates:                   Coordinates{Latitude: 51.5, Longitude: -0.12},
					CoordinateUncertaintyInMetres: new(3.5),
					ElevationInMetres:             new(125.25),
				},
				Species: species.Species{
					ID:             nestedID,
					GBIFKey:        2480498,
					VernacularName: "European robin",
				},
				Media: []Media{{
					ID:     nestedID,
					Type:   "StillImage",
					Format: "image/jpeg",
					URL:    "https://example.com/robin.jpg",
					Title:  "Robin",
				}},
			},
			want: CreateParams{
				GBIFKey:         123,
				OccurrenceID:    "occurrence-123",
				HappenedAt:      time.Date(2026, 9, 22, 12, 30, 0, 0, time.UTC),
				BasisOfRecord:   "HUMAN_OBSERVATION",
				RecordedBy:      "Observer",
				IndividualCount: new(2),
				Remarks:         "A robin",
				ReferenceURL:    "https://example.com/123",
				Location: Location{
					Locality:                      "London",
					StateProvince:                 "England",
					Country:                       "United Kingdom",
					CountryCode:                   "GB",
					Coordinates:                   Coordinates{Latitude: 51.5, Longitude: -0.12},
					CoordinateUncertaintyInMetres: new(3.5),
					ElevationInMetres:             new(125.25),
				},
				SpeciesID: nestedID,
				Media: []Media{{
					ID:     nestedID,
					Type:   "StillImage",
					Format: "image/jpeg",
					URL:    "https://example.com/robin.jpg",
					Title:  "Robin",
				}},
			},
		},
		"Empty media": {
			input: Sighting{
				Media: []Media{},
			},
			want: CreateParams{
				Media: []Media{},
			},
		},
		"Optional zero values": {
			input: Sighting{
				IndividualCount: new(0),
				Location: Location{
					CoordinateUncertaintyInMetres: new(0.0),
					ElevationInMetres:             new(0.0),
				},
			},
			want: CreateParams{
				IndividualCount: new(0),
				Location: Location{
					CoordinateUncertaintyInMetres: new(0.0),
					ElevationInMetres:             new(0.0),
				},
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

func TestParseCoordinates(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		latitude  float64
		longitude float64
		want      Coordinates
		wantErr   bool
	}{
		"Valid":              {latitude: 51.5, longitude: -0.12, want: Coordinates{Latitude: 51.5, Longitude: -0.12}},
		"Null island":        {latitude: 0, longitude: 0, want: Coordinates{}},
		"Northern boundary":  {latitude: 90, longitude: 180, want: Coordinates{Latitude: 90, Longitude: 180}},
		"Southern boundary":  {latitude: -90, longitude: -180, want: Coordinates{Latitude: -90, Longitude: -180}},
		"Latitude too high":  {latitude: 90.1, wantErr: true},
		"Latitude too low":   {latitude: -90.1, wantErr: true},
		"Longitude too high": {longitude: 180.1, wantErr: true},
		"Longitude too low":  {longitude: -180.1, wantErr: true},
		"Both out of range":  {latitude: 100, longitude: 200, wantErr: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseCoordinates(test.latitude, test.longitude)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
			if test.wantErr {
				assert.ErrorIs(t, err, ErrInvalid)
			}
		})
	}
}
