package lookup

import (
	"testing"
	"uuid"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func parse(t *testing.T, s string) uuid.UUID {
	t.Helper()
	id, err := uuid.Parse(s)
	require.NoError(t, err)
	return id
}

func TestByIDPrefix(t *testing.T) {
	t.Parallel()

	a := parse(t, "3f2a9c1b-0000-4000-8000-000000000001")
	b := parse(t, "3f2a9c1b-0000-4000-8000-000000000002")
	c := parse(t, "7e1d0a2c-0000-4000-8000-000000000003")

	tt := map[string]struct {
		input   string
		want    uuid.UUID
		wantErr bool
	}{
		"Unique prefix": {input: "7e1d", want: c},
		"Full ID":       {input: a.String(), want: a},
		"Ambiguous":     {input: "3f2a9c1b", wantErr: true},
		"No match":      {input: "ffff", wantErr: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got, err := ByIDPrefix([]uuid.UUID{a, b, c}, test.input, "thing", "things", func(id uuid.UUID) uuid.UUID {
				return id
			})
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}
