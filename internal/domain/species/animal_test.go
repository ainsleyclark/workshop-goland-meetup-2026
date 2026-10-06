package species

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAnimal_TaxonKeys(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Animal
		want  []int
	}{
		"One species":                 {AnimalHumpbackWhale, []int{5220086}},
		"Several species":             {AnimalZebra, []int{2440892, 2440888, 2440894}},
		"A genus key":                 {AnimalGiraffe, []int{7716694}},
		"A group and its own species": {AnimalAlligator, []int{2441370, 2441368}},
		"Unknown animal":              {Animal("pangolin"), nil},
		"Empty animal":                {Animal(""), nil},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, test.input.TaxonKeys())
		})
	}
}

func TestAnimal_TaxonKeysAreCopied(t *testing.T) {
	t.Parallel()

	keys := AnimalZebra.TaxonKeys()
	keys[0] = 0

	assert.Equal(t, []int{2440892, 2440888, 2440894}, AnimalZebra.TaxonKeys())
}

func TestAnimal_Known(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input Animal
		want  bool
	}{
		"Supported":   {AnimalLion, true},
		"Unsupported": {Animal("pangolin"), false},
		"Empty":       {Animal(""), false},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, test.input.Known())
		})
	}
}

func TestAnimals(t *testing.T) {
	t.Parallel()

	all := Animals()

	assert.NotEmpty(t, all)
	assert.IsIncreasing(t, all, "listed in name order so the CLI is stable")
	assert.Contains(t, all, AnimalHumpbackWhale)

	for _, a := range all {
		assert.True(t, a.Known(), "every listed animal resolves to at least one taxon")
	}
}

func TestParseAnimal(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		input   string
		want    Animal
		wantErr bool
	}{
		"Exact name":    {input: "humpback_whale", want: AnimalHumpbackWhale},
		"Title case":    {input: "Humpback Whale", want: AnimalHumpbackWhale},
		"Hyphenated":    {input: "humpback-whale", want: AnimalHumpbackWhale},
		"Padded":        {input: "  lion  ", want: AnimalLion},
		"Unsupported":   {input: "pangolin", wantErr: true},
		"Empty":         {input: "", wantErr: true},
		"Only spacing":  {input: "   ", wantErr: true},
		"Mangled input": {input: "humpback  whale", wantErr: true},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			got, err := ParseAnimal(test.input)
			assert.Equal(t, test.wantErr, err != nil)
			assert.Equal(t, test.want, got)
		})
	}
}

func TestAnimal_String(t *testing.T) {
	t.Parallel()
	assert.Equal(t, "humpback_whale", AnimalHumpbackWhale.String())
}
