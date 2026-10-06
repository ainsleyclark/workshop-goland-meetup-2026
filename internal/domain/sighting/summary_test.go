package sighting

import (
	"testing"
	"time"
	"workshop/internal/domain/species"

	"github.com/stretchr/testify/assert"
)

func TestSummarise(t *testing.T) {
	t.Parallel()

	jan := time.Date(2025, 1, 10, 0, 0, 0, 0, time.UTC)
	jun := time.Date(2026, 6, 20, 0, 0, 0, 0, time.UTC)

	robin := Sighting{
		HappenedAt:      jan,
		RecordedBy:      "Ada",
		IndividualCount: new(2),
		Location:        Location{Country: "United Kingdom"},
		Species:         species.Species{ScientificName: "Erithacus rubecula"},
		Media:           []Media{{Type: "StillImage"}, {Type: "Sound"}},
	}
	whale := Sighting{
		HappenedAt: jun,
		RecordedBy: "Grace",
		Location:   Location{Country: "Kenya"},
		Species:    species.Species{ScientificName: "Megaptera novaeangliae"},
	}

	tt := map[string]struct {
		input []Sighting
		want  Summary
	}{
		"No sightings": {want: Summary{}},
		"One sighting": {
			input: []Sighting{robin},
			want: Summary{
				Sightings: 1, Species: 1, Countries: 1, Media: 2, Individuals: 2,
				Recorders: 1, Earliest: jan, Latest: jan,
			},
		},
		"Two species across two countries": {
			input: []Sighting{whale, robin},
			want: Summary{
				Sightings: 2, Species: 2, Countries: 2, Media: 2, Individuals: 2,
				Recorders: 2, Earliest: jan, Latest: jun,
			},
		},
		"Blank country and recorder are not counted": {
			input: []Sighting{{HappenedAt: jan}},
			want: Summary{
				Sightings: 1, Species: 1, Earliest: jan, Latest: jan,
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, Summarise(test.input))
		})
	}
}

func TestCountBy(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input []Sighting
		want  map[string]int
	}{
		"No sightings": {want: map[string]int{}},
		"Counts repeats": {
			input: []Sighting{
				{Location: Location{Country: "Kenya"}},
				{Location: Location{Country: "Kenya"}},
				{Location: Location{Country: "Chile"}},
			},
			want: map[string]int{"Kenya": 2, "Chile": 1},
		},
		"Skips blanks": {
			input: []Sighting{{Location: Location{Country: ""}}},
			want:  map[string]int{},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, CountBy(test.input, func(s Sighting) string { return s.Location.Country }))
		})
	}
}

func TestCountByYear(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input []Sighting
		want  map[int]int
	}{
		"No sightings": {want: map[int]int{}},
		"Groups by year": {
			input: []Sighting{
				{HappenedAt: time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)},
				{HappenedAt: time.Date(2025, 12, 31, 0, 0, 0, 0, time.UTC)},
				{HappenedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)},
			},
			want: map[int]int{2025: 2, 2026: 1},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, CountByYear(test.input))
		})
	}
}

func TestTop(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input map[string]int
		limit int
		want  []Tally
	}{
		"No counts": {input: map[string]int{}, limit: 3, want: []Tally{}},
		"Highest first": {
			input: map[string]int{"Kenya": 2, "Chile": 9, "Peru": 5},
			limit: 3,
			want:  []Tally{{"Chile", 9}, {"Peru", 5}, {"Kenya", 2}},
		},
		"Ties break by name so the order is stable": {
			input: map[string]int{"Peru": 4, "Chile": 4, "Kenya": 4},
			limit: 3,
			want:  []Tally{{"Chile", 4}, {"Kenya", 4}, {"Peru", 4}},
		},
		"Keeps at most the limit": {
			input: map[string]int{"Kenya": 2, "Chile": 9, "Peru": 5},
			limit: 2,
			want:  []Tally{{"Chile", 9}, {"Peru", 5}},
		},
		"A limit of zero keeps everything": {
			input: map[string]int{"Chile": 9, "Peru": 5},
			limit: 0,
			want:  []Tally{{"Chile", 9}, {"Peru", 5}},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, Top(test.input, test.limit))
		})
	}
}
