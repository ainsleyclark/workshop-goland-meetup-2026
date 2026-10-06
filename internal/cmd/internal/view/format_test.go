package view

import (
	"testing"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"

	"github.com/stretchr/testify/assert"
)

func TestSpeciesName(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input species.Species
		want  string
	}{
		"Vernacular and canonical": {
			input: species.Species{VernacularName: "Humpback Whale", CanonicalName: "Megaptera novaeangliae"},
			want:  "Humpback Whale (Megaptera novaeangliae)",
		},
		"Vernacular only":     {input: species.Species{VernacularName: "Humpback Whale"}, want: "Humpback Whale"},
		"Scientific fallback": {input: species.Species{ScientificName: "Megaptera novaeangliae"}, want: "Megaptera novaeangliae"},
		"Canonical preferred": {input: species.Species{CanonicalName: "Panthera leo", ScientificName: "Panthera leo (Linnaeus, 1758)"}, want: "Panthera leo"},
		"No names":            {input: species.Species{}, want: "Unknown species"},
		"Lowercase vernacular": {
			input: species.Species{VernacularName: "humpback whale", CanonicalName: "Megaptera novaeangliae"},
			want:  "Humpback Whale (Megaptera novaeangliae)",
		},
		"Mixed case vernacular kept": {input: species.Species{VernacularName: "Risso's dolphin"}, want: "Risso's dolphin"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, SpeciesName(test.input))
		})
	}
}

func TestLocation(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input sighting.Location
		want  string
	}{
		"All parts":         {input: sighting.Location{Locality: "Watamu", StateProvince: "Kilifi", Country: "Kenya"}, want: "Watamu, Kilifi, Kenya"},
		"Country only":      {input: sighting.Location{Country: "Kenya"}, want: "Kenya"},
		"Coordinates only":  {input: sighting.Location{Coordinates: sighting.Coordinates{Latitude: -3.37123, Longitude: 40.02131}}, want: "-3.3712, 40.0213"},
		"Duplicate country": {input: sighting.Location{Locality: "Kwale", StateProvince: "Kenya", Country: "Kenya"}, want: "Kwale, Kenya"},
		"Duplicate state":   {input: sighting.Location{Locality: "kilifi", StateProvince: "Kilifi", Country: "Kenya"}, want: "kilifi, Kenya"},
		"Shouted locality": {
			input: sighting.Location{Locality: "KISITE-MPUNGUTI-MARINE NATIONAL RESERVE", Country: "Kenya"},
			want:  "Kisite-Mpunguti-Marine National Reserve, Kenya",
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, Location(test.input))
		})
	}
}

func TestHumanise(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Empty":      {input: "", want: ""},
		"Single":     {input: "SPECIES", want: "Species"},
		"Underscore": {input: "HUMAN_OBSERVATION", want: "Human observation"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, Humanise(test.input))
		})
	}
}

func TestTidy(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input string
		want  string
	}{
		"Empty":       {input: "", want: ""},
		"Digits only": {input: "123", want: "123"},
		"Mixed case":  {input: "Watamu Marine Park", want: "Watamu Marine Park"},
		"Lower case":  {input: "kwale", want: "kwale"},
		"Shouted":     {input: "DIANI BEACH", want: "Diani Beach"},
		"Hyphenated":  {input: "KISITE-MPUNGUTI", want: "Kisite-Mpunguti"},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, tidy(test.input))
		})
	}
}
