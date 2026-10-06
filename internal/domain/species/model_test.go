package species

import (
	"testing"
	"time"
	"uuid"

	"github.com/stretchr/testify/assert"
)

func TestSpecies_CreateParams(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Species
		want  CreateParams
	}{
		"Zero value": {},
		"Timestamps excluded": {
			input: Species{
				CreatedAt: time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC),
				UpdatedAt: time.Date(2026, 1, 2, 0, 0, 0, 0, time.UTC),
			},
		},
		"ID only": {
			input: Species{
				ID: uuid.New(),
			},
			want: CreateParams{},
		},
		"All fields": {
			input: Species{
				ID:             uuid.New(),
				GBIFKey:        2480498,
				Kingdom:        "Animalia",
				Phylum:         "Chordata",
				Class:          "Aves",
				Order:          "Passeriformes",
				Family:         "Muscicapidae",
				Genus:          "Erithacus",
				Species:        "Erithacus rubecula",
				Rank:           "SPECIES",
				ScientificName: "Erithacus rubecula (Linnaeus, 1758)",
				CanonicalName:  "Erithacus rubecula",
				VernacularName: "European robin",
			},
			want: CreateParams{
				GBIFKey:        2480498,
				Kingdom:        "Animalia",
				Phylum:         "Chordata",
				Class:          "Aves",
				Order:          "Passeriformes",
				Family:         "Muscicapidae",
				Genus:          "Erithacus",
				Species:        "Erithacus rubecula",
				Rank:           "SPECIES",
				ScientificName: "Erithacus rubecula (Linnaeus, 1758)",
				CanonicalName:  "Erithacus rubecula",
				VernacularName: "European robin",
			},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, test.input.CreateParams())
		})
	}
}
