package components

import (
	"cmp"
	"fmt"
	"slices"

	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"
)

// Room is one species' sightings, hung together as a room of the exhibition.
type Room struct {
	Number    int
	Species   species.Species
	Sightings []sighting.Sighting
	Summary   sighting.Summary
}

// NewRooms gives each species a room, the most seen first, keeping the sightings' order inside each.
func NewRooms(sightings []sighting.Sighting) []Room {
	index := make(map[string]int)
	var rooms []Room

	for _, s := range sightings {
		key := cmp.Or(s.Species.ScientificName, displayName(s.Species))
		i, ok := index[key]
		if !ok {
			i = len(rooms)
			index[key] = i
			rooms = append(rooms, Room{Species: s.Species})
		}
		rooms[i].Sightings = append(rooms[i].Sightings, s)
	}

	slices.SortStableFunc(rooms, func(a, b Room) int {
		return cmp.Or(
			cmp.Compare(len(b.Sightings), len(a.Sightings)),
			cmp.Compare(displayName(a.Species), displayName(b.Species)),
		)
	})

	for i := range rooms {
		rooms[i].Number = i + 1
		rooms[i].Summary = sighting.Summarise(rooms[i].Sightings)
	}

	return rooms
}

// roomID is the anchor a floor-plan tile links to, such as room-01.
func roomID(r Room) string {
	return fmt.Sprintf("room-%02d", r.Number)
}

// roomLabel numbers a room the way a gallery does, such as Room 01.
func roomLabel(r Room) string {
	return fmt.Sprintf("Room %02d", r.Number)
}

// dateRange spans a room's first and latest sighting, such as Mar 2019 – Sep 2026.
func dateRange(s sighting.Summary) string {
	from, to := s.Earliest.Format("Jan 2006"), s.Latest.Format("Jan 2006")
	if from == to {
		return to
	}
	return from + " – " + to
}

// observedBy names who recorded a sighting, or says nobody did.
func observedBy(recordedBy string) string {
	return cmp.Or(recordedBy, "Unknown")
}

// headcount shows how many animals a sighting counted, or nothing when it didn't say.
func headcount(n *int) string {
	if n == nil {
		return ""
	}
	return thousands(*n)
}
