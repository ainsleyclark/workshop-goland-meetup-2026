// Package lookup resolves the short ID prefixes that commands print
// back into the records they belong to.
package lookup

import (
	"fmt"
	"slices"
	"strings"
	"uuid"
)

// ByIDPrefix returns the single item whose ID starts with prefix,
// erroring when nothing matches or the prefix is ambiguous. Noun and
// plural name the items in those errors, e.g. "sighting", "sightings".
func ByIDPrefix[T any](items []T, prefix, noun, plural string, id func(T) uuid.UUID) (T, error) {
	var zero T

	matches := slices.DeleteFunc(items, func(v T) bool {
		return !strings.HasPrefix(id(v).String(), prefix)
	})

	switch len(matches) {
	case 0:
		return zero, fmt.Errorf("no %s matches %q", noun, prefix)
	case 1:
		return matches[0], nil
	default:
		return zero, fmt.Errorf("%q matches %d %s, use a longer prefix", prefix, len(matches), plural)
	}
}
