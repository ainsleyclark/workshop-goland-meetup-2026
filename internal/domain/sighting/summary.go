package sighting

import (
	"cmp"
	"slices"
	"time"
)

type (
	// Summary is an aggregate over a set of sightings. Anything
	// reporting on sightings tends to lead with it, so the numbers are
	// worked out once here rather than wherever they are shown.
	Summary struct {
		Sightings   int
		Species     int
		Countries   int
		Media       int
		Individuals int
		Recorders   int
		Earliest    time.Time
		Latest      time.Time
	}
	// Tally is one row of a leaderboard: a name and how often it appeared.
	Tally struct {
		Name  string
		Count int
	}
)

// Summarise reduces sightings to the headline figures. It reads every
// sighting once, so callers should hold on to the result.
func Summarise(sightings []Sighting) Summary {
	s := Summary{Sightings: len(sightings)}

	species := make(map[string]struct{})
	countries := make(map[string]struct{})
	recorders := make(map[string]struct{})

	for _, sight := range sightings {
		species[sight.Species.ScientificName] = struct{}{}
		if sight.Location.Country != "" {
			countries[sight.Location.Country] = struct{}{}
		}
		if sight.RecordedBy != "" {
			recorders[sight.RecordedBy] = struct{}{}
		}

		s.Media += len(sight.Media)
		if sight.IndividualCount != nil {
			s.Individuals += *sight.IndividualCount
		}

		if s.Earliest.IsZero() || sight.HappenedAt.Before(s.Earliest) {
			s.Earliest = sight.HappenedAt
		}
		if sight.HappenedAt.After(s.Latest) {
			s.Latest = sight.HappenedAt
		}
	}

	s.Species = len(species)
	s.Countries = len(countries)
	s.Recorders = len(recorders)

	return s
}

// CountBy tallies sightings by the key the given function returns,
// skipping any sighting the function gives no key for. How a thing is
// named is a presentation decision, so the caller supplies it: the CLI
// labels a species differently from the web pages, and both count the
// same way.
func CountBy(sightings []Sighting, key func(Sighting) string) map[string]int {
	counts := make(map[string]int)
	for _, s := range sightings {
		if k := key(s); k != "" {
			counts[k]++
		}
	}
	return counts
}

// CountByYear returns how many sightings happened in each year.
func CountByYear(sightings []Sighting) map[int]int {
	counts := make(map[int]int)
	for _, s := range sightings {
		counts[s.HappenedAt.Year()]++
	}
	return counts
}

// Top ranks counts highest first, breaking ties by name so the order is
// stable between renders, and keeps at most n of them.
func Top(counts map[string]int, n int) []Tally {
	rows := make([]Tally, 0, len(counts))
	for name, count := range counts {
		rows = append(rows, Tally{Name: name, Count: count})
	}

	slices.SortFunc(rows, func(a, b Tally) int {
		if c := cmp.Compare(b.Count, a.Count); c != 0 {
			return c
		}
		return cmp.Compare(a.Name, b.Name)
	})

	if n > 0 && len(rows) > n {
		rows = rows[:n]
	}

	return rows
}
