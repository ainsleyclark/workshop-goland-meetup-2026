package gbif

import (
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPaginationOptions_Query(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input PaginationOptions
		want  url.Values
	}{
		"Zero values omit pagination": {
			input: PaginationOptions{},
			want:  url.Values{},
		},
		"Limit only": {
			input: PaginationOptions{Limit: 100},
			want:  url.Values{"limit": {"100"}},
		},
		"Offset only": {
			input: PaginationOptions{Offset: 200},
			want:  url.Values{"offset": {"200"}},
		},
		"Limit and offset": {
			input: PaginationOptions{Limit: 100, Offset: 200},
			want:  url.Values{"limit": {"100"}, "offset": {"200"}},
		},
		"Negative values are passed through": {
			input: PaginationOptions{Limit: -1, Offset: -2},
			want:  url.Values{"limit": {"-1"}, "offset": {"-2"}},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, test.input.query())
		})
	}
}
