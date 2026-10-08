package components

import (
	"net/url"
	"path"
	"strings"
)

// iNaturalist's photo sizes, 240 and 500px on the longest side.
const (
	sizeSmall  = "small"
	sizeMedium = "medium"
)

// photoURL swaps an iNaturalist original for a smaller copy, leaving any other URL as it is.
func photoURL(raw, size string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Host != "inaturalist-open-data.s3.amazonaws.com" {
		return raw
	}

	dir, file := path.Split(u.Path)
	if !strings.HasPrefix(file, "original.") {
		return raw
	}

	u.Path = dir + size + path.Ext(file)
	return u.String()
}
