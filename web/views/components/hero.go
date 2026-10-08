package components

import (
	"cmp"
	"fmt"
	"strconv"
	"strings"

	"workshop/internal/domain/sighting"
)

// HeroProps is what the hero shows: the exhibition's totals and its newest sighting.
type HeroProps struct {
	Summary sighting.Summary
	Latest  *sighting.Sighting
}

// NewHeroProps totals the sightings, which arrive newest first, and keeps the newest for the phone.
func NewHeroProps(sightings []sighting.Sighting) HeroProps {
	props := HeroProps{Summary: sighting.Summarise(sightings)}
	if len(sightings) > 0 {
		props.Latest = &sightings[0]
	}
	return props
}

// thousands writes n with a comma between each group of three digits, such as 4,263.
func thousands(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

// plural picks the word that agrees with n.
func plural(n int, one, many string) string {
	if n == 1 {
		return "1 " + one
	}
	return thousands(n) + " " + many
}

// counterDigits pads n to the seven digits of a hit counter, such as 0004263.
func counterDigits(n int) []string {
	return strings.Split(fmt.Sprintf("%07d", n), "")
}

// textSender names who to credit in the phone's message.
func textSender(recordedBy string) string {
	return cmp.Or(recordedBy, "someone")
}
