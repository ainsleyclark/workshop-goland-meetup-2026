package cmd

import (
	"testing"
	"workshop/internal/infra/config"

	"github.com/stretchr/testify/assert"
)

func TestWithLocalAPIs(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		local bool
		want  config.Config
	}{
		"Local": {
			local: true,
			want:  config.Config{GBIFBaseURL: localGBIFBaseURL, OpenMeteoBaseURL: localOpenMeteoBaseURL},
		},
		"Not local": {local: false, want: config.Config{}},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := withLocalAPIs(config.Config{}, test.local)
			assert.Equal(t, test.want, got)
		})
	}
}
