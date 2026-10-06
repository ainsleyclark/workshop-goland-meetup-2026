package species

import (
	"fmt"
	"maps"
	"slices"
	"strings"
)

// Animal names an animal, or a group of animals, that the application
// can search for and filter by.
type Animal string

// The animals this application supports. It's not everything that the
// GBIF API supports, just a small selection.
const (
	AnimalLion           Animal = "lion"
	AnimalLeopard        Animal = "leopard"
	AnimalCheetah        Animal = "cheetah"
	AnimalElephant       Animal = "elephant"
	AnimalRhinoceros     Animal = "rhinoceros"
	AnimalGiraffe        Animal = "giraffe"
	AnimalZebra          Animal = "zebra"
	AnimalHippopotamus   Animal = "hippopotamus"
	AnimalAfricanBuffalo Animal = "african_buffalo"
	AnimalSpottedHyena   Animal = "spotted_hyena"
	AnimalAfricanWildDog Animal = "african_wild_dog"
	AnimalWildebeest     Animal = "wildebeest"
	AnimalImpala         Animal = "impala"
	AnimalCommonWarthog  Animal = "common_warthog"

	AnimalBlueWhale               Animal = "blue_whale"
	AnimalHumpbackWhale           Animal = "humpback_whale"
	AnimalSpermWhale              Animal = "sperm_whale"
	AnimalFinWhale                Animal = "fin_whale"
	AnimalCommonMinkeWhale        Animal = "common_minke_whale"
	AnimalGreyWhale               Animal = "grey_whale"
	AnimalNorthAtlanticRightWhale Animal = "north_atlantic_right_whale"
	AnimalBowheadWhale            Animal = "bowhead_whale"
	AnimalBeluga                  Animal = "beluga"
	AnimalNarwhal                 Animal = "narwhal"
	AnimalOrca                    Animal = "orca"
	AnimalCommonBottlenoseDolphin Animal = "common_bottlenose_dolphin"
	AnimalCommonDolphin           Animal = "common_dolphin"
	AnimalSpinnerDolphin          Animal = "spinner_dolphin"
	AnimalDuskyDolphin            Animal = "dusky_dolphin"
	AnimalRissosDolphin           Animal = "rissos_dolphin"
	AnimalAmazonRiverDolphin      Animal = "amazon_river_dolphin"

	AnimalAlligator          Animal = "alligator"
	AnimalAmericanAlligator  Animal = "american_alligator"
	AnimalChineseAlligator   Animal = "chinese_alligator"
	AnimalNileCrocodile      Animal = "nile_crocodile"
	AnimalSaltwaterCrocodile Animal = "saltwater_crocodile"
	AnimalAmericanCrocodile  Animal = "american_crocodile"
	AnimalGharial            Animal = "gharial"
)

// String implements fmt.Stringer on Animal.
func (a Animal) String() string {
	return string(a)
}

// TaxonKeys returns the GBIF Backbone taxon IDs this animal stands for.
// A broad name such as "elephant" covers several species, so this is a
// set.
func (a Animal) TaxonKeys() []int {
	return slices.Clone(animalTaxonKeys[a])
}

// Known reports whether this is an animal the application supports. Prefer
// ParseAnimal where the value came from outside; Known is for a value that
// arrived some other way.
func (a Animal) Known() bool {
	return len(animalTaxonKeys[a]) > 0
}

// ParseAnimal converts a name such as "humpback_whale", "Humpback Whale" or
// "humpback-whale" into a supported Animal, and is the only way an Animal
// should be obtained from user input. An unsupported name is an error rather
// than an Animal nobody can look up.
func ParseAnimal(value string) (Animal, error) {
	name := strings.NewReplacer(" ", "_", "-", "_").Replace(strings.ToLower(strings.TrimSpace(value)))

	animal := Animal(name)
	if !animal.Known() {
		return "", fmt.Errorf("unknown animal %q, run 'workshop animals' to see the options", value)
	}

	return animal, nil
}

// Animals returns every supported animal, in name order.
func Animals() []Animal {
	return slices.Sorted(maps.Keys(animalTaxonKeys))
}

var animalTaxonKeys = map[Animal][]int{
	AnimalLion:           {5219404},                                     // Panthera leo
	AnimalLeopard:        {5219436},                                     // Panthera pardus
	AnimalCheetah:        {2435270},                                     // Acinonyx jubatus
	AnimalElephant:       {2435350, 2435349, 5219461},                   // Loxodonta africana, L. cyclotis, Elephas maximus
	AnimalRhinoceros:     {5220111, 2440880, 5220113, 5220112, 2440878}, // Diceros bicornis, Ceratotherium simum, Rhinoceros unicorns, R. sondaicus, Dicerorhinus sumatrensis
	AnimalGiraffe:        {7716694},                                     // Giraffa
	AnimalZebra:          {2440892, 2440888, 2440894},                   // Equus quagga, E. zebra, E. grevyi
	AnimalHippopotamus:   {2441247},                                     // Hippopotamus amphibius
	AnimalAfricanBuffalo: {2441034},                                     // Syncerus caffer
	AnimalSpottedHyena:   {5218781},                                     // Crocuta crocuta
	AnimalAfricanWildDog: {5219317},                                     // Lycaon pictus
	AnimalWildebeest:     {2441104},                                     // Connochaetes
	AnimalImpala:         {2441144},                                     // Aepyceros melampus
	AnimalCommonWarthog:  {2441212},                                     // Phacochoerus africanus

	AnimalBlueWhale:               {2440735}, // Balaenoptera musculus
	AnimalHumpbackWhale:           {5220086}, // Megaptera novaeangliae
	AnimalSpermWhale:              {8123917}, // Physeter macrocephalus
	AnimalFinWhale:                {2440718}, // Balaenoptera physalus
	AnimalCommonMinkeWhale:        {2440728}, // Balaenoptera acutorostrata
	AnimalGreyWhale:               {2440704}, // Eschrichtius robustus
	AnimalNorthAtlanticRightWhale: {2440339}, // Eubalaena glacialis
	AnimalBowheadWhale:            {2440330}, // Balaena mysticetus
	AnimalBeluga:                  {5220003}, // Delphinapterus leucas
	AnimalNarwhal:                 {5220008}, // Monodon monoceros
	AnimalOrca:                    {2440483}, // Orcinus orca
	AnimalCommonBottlenoseDolphin: {2440447}, // Tursiops truncatus
	AnimalCommonDolphin:           {8324617}, // Delphinus delphis
	AnimalSpinnerDolphin:          {5220060}, // Stenella longirostris
	AnimalDuskyDolphin:            {5220077}, // Lagenorhynchus obscurus
	AnimalRissosDolphin:           {5220027}, // Grampus griseus
	AnimalAmazonRiverDolphin:      {2440314}, // Inia geoffrensis

	AnimalAlligator:          {2441370, 2441368}, // Alligator mississippiensis, A. sinensis
	AnimalAmericanAlligator:  {2441370},          // Alligator mississippiensis
	AnimalChineseAlligator:   {2441368},          // Alligator sinensis
	AnimalNileCrocodile:      {2441341},          // Crocodylus niloticus
	AnimalSaltwaterCrocodile: {2441333},          // Crocodylus porosus
	AnimalAmericanCrocodile:  {2441329},          // Crocodylus acutus
	AnimalGharial:            {2441360},          // Gavialis gangeticus
}
