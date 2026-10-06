package country

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParse(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		want    Code
		wantErr bool
	}{
		"Empty":      {input: "", want: ""},
		"Upper case": {input: "KE", want: Kenya},
		"Lower case": {input: " ke ", want: Kenya},
		"Too long":   {input: "Kenya", wantErr: true},
		"Too short":  {input: "K", wantErr: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := Parse(test.input)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}
