package sighting

import (
	"testing"
	"workshop/internal/domain/species"

	"github.com/stretchr/testify/assert"
)

func TestListFilter_TaxonKeys(t *testing.T) {
	t.Parallel()

	filter := ListFilter{Animals: []species.Animal{
		species.AnimalLion,
		species.AnimalElephant,
	}}

	assert.Equal(t, []int{5219404, 2435350, 2435349, 5219461}, filter.TaxonKeys())
}

func TestTaxonKeys(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input []species.Animal
		want  []int
	}{
		"No animals": {},
		"One animal with one species": {
			input: []species.Animal{species.AnimalLion},
			want:  []int{5219404},
		},
		"One animal with several species": {
			input: []species.Animal{species.AnimalElephant},
			want:  []int{2435350, 2435349, 5219461},
		},
		"Several animals preserve their order": {
			input: []species.Animal{species.AnimalHumpbackWhale, species.AnimalLion},
			want:  []int{5220086, 5219404},
		},
		"Overlapping animals are deduplicated": {
			input: []species.Animal{species.AnimalAlligator, species.AnimalAmericanAlligator},
			want:  []int{2441370, 2441368},
		},
		"Repeated animals are deduplicated": {
			input: []species.Animal{species.AnimalLion, species.AnimalLion},
			want:  []int{5219404},
		},
		"Unknown animals are ignored": {
			input: []species.Animal{species.Animal("pangolin"), species.AnimalLion},
			want:  []int{5219404},
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, taxonKeys(test.input))
		})
	}
}
