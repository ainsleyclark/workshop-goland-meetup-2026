package cmd

import (
	"cmp"
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"uuid"
	"workshop/internal/cmd/internal/lookup"
	"workshop/internal/cmd/internal/prompt"
	"workshop/internal/cmd/internal/view"
	"workshop/internal/common/printer"
	"workshop/internal/common/printer/styles"
	"workshop/internal/domain/sighting"
	"workshop/internal/domain/species"

	"github.com/charmbracelet/huh"
	"github.com/urfave/cli/v3"
)

// recentLimit is how many of a species' sightings the detail view shows.
const recentLimit = 5

// speciesWithCount is a species alongside how
// many sightings of it are saved.
type speciesWithCount struct {
	species.Species
	Sightings int
}

func speciesCmd(app *app) *cli.Command {
	return &cli.Command{
		Name:    "species",
		Aliases: []string{"sp"},
		Usage:   "List and inspect species, which are saved while ingesting sightings",
		Commands: []*cli.Command{
			{
				Name:    "list",
				Aliases: []string{"ls"},
				Usage:   "List saved species, most sighted first",
				Action: func(ctx context.Context, _ *cli.Command) error {
					counted, err := speciesCounts(ctx, app)
					if err != nil {
						return err
					}

					var hint view.Hint
					if len(counted) > 0 {
						hint = view.Hintf("%d species %s workshop species find <gbif key>", len(counted), styles.IconArrow)
					}

					view.Render(app.Printer,
						view.Table[speciesWithCount]{
							Title:   "Species",
							Headers: []string{"GBIF KEY", "NAME", "SCIENTIFIC NAME", "FAMILY", "CLASS", "SIGHTINGS"},
							Items:   counted,
							Empty:   "No species yet, they're saved when you run 'workshop sightings ingest'",
							Row: func(s speciesWithCount) []string {
								return []string{
									fmt.Sprint(s.GBIFKey),
									cmp.Or(s.VernacularName, "—"),
									cmp.Or(s.CanonicalName, s.ScientificName),
									s.Family,
									s.Class,
									fmt.Sprint(s.Sightings),
								}
							},
						},
						hint,
					)
					return nil
				},
			},
			{
				Name:      "find",
				Aliases:   []string{"show"},
				Usage:     "Show a species' taxonomy and its most recent sightings",
				ArgsUsage: "<gbif key, id or id prefix>",
				Metadata:  prompt.With(pickSpecies(app)),
				Action: func(ctx context.Context, cmd *cli.Command) error {
					s, err := findSpecies(ctx, app, cmd.Args().First())
					if err != nil {
						return err
					}

					all, err := allSightings(ctx, app)
					if err != nil {
						return err
					}
					recent := slices.DeleteFunc(all, func(v sighting.Sighting) bool {
						return v.Species.ID != s.ID
					})
					view.Render(app.Printer,
						view.Card{Title: view.SpeciesName(s), Rows: []printer.KV{
							{Key: "ID", Value: s.ID.String()},
							{Key: "GBIF key", Value: fmt.Sprint(s.GBIFKey)},
							{Key: "Scientific", Value: s.ScientificName},
							{Key: "Rank", Value: view.Humanise(s.Rank)},
							{Key: "Sightings", Value: fmt.Sprint(len(recent))},
						}},
						view.Lineage{Title: "Taxonomy", Rows: []printer.KV{
							{Key: "kingdom", Value: s.Kingdom},
							{Key: "phylum", Value: s.Phylum},
							{Key: "class", Value: s.Class},
							{Key: "order", Value: s.Order},
							{Key: "family", Value: s.Family},
							{Key: "genus", Value: s.Genus},
							{Key: "species", Value: s.Species},
						}},
					)
					if len(recent) > 0 {
						recent = recent[:min(len(recent), recentLimit)]
						view.Render(app.Printer,
							sightingsTable("Recent sightings", recent),
							view.Hintf("%s workshop sightings find %s", styles.IconArrow, view.ShortID(recent[0].ID)),
						)
					}
					return nil
				},
			},
		},
	}
}

// speciesCounts lists every saved species with how many sightings
// of it are saved, most sighted first.
func speciesCounts(ctx context.Context, app *app) ([]speciesWithCount, error) {
	all, err := allSpecies(ctx, app)
	if err != nil {
		return nil, err
	}
	sightings, err := allSightings(ctx, app)
	if err != nil {
		return nil, err
	}
	counts := countBySpecies(sightings)

	out := make([]speciesWithCount, 0, len(all))
	for _, s := range all {
		out = append(out, speciesWithCount{Species: s, Sightings: counts[s.ID]})
	}
	slices.SortFunc(out, func(a, b speciesWithCount) int {
		return cmp.Or(cmp.Compare(b.Sightings, a.Sightings), cmp.Compare(view.SpeciesName(a.Species), view.SpeciesName(b.Species)))
	})

	return out, nil
}

func countBySpecies(sightings []sighting.Sighting) map[uuid.UUID]int {
	counts := make(map[uuid.UUID]int)
	for _, s := range sightings {
		counts[s.Species.ID]++
	}
	return counts
}

// pickSpecies lets the user choose a saved species,
// most sighted first, instead of typing its GBIF key.
func pickSpecies(app *app) prompt.Func {
	return func(ctx context.Context) ([]string, error) {
		counted, err := speciesCounts(ctx, app)
		if err != nil {
			return nil, err
		} else if len(counted) == 0 {
			return nil, errors.New("no species yet, ingest some sightings first")
		}

		opts := make([]huh.Option[string], 0, len(counted))
		for _, s := range counted {
			label := fmt.Sprintf("%s %s %d sighting(s)", view.SpeciesName(s.Species), styles.IconDot, s.Sightings)
			opts = append(opts, huh.NewOption(label, strconv.Itoa(s.GBIFKey)))
		}

		key, err := prompt.Choose("Which species?", opts...)
		return []string{key}, err
	}
}

// findSpecies resolves a GBIF key, a full UUID, or a
// unique prefix of one.
func findSpecies(ctx context.Context, app *app, arg string) (species.Species, error) {
	arg = strings.ToLower(strings.TrimSpace(arg))
	if arg == "" {
		return species.Species{}, errors.New("missing species, e.g. 'workshop species find 2440735'")
	}

	var (
		s   species.Species
		err error
	)
	if key, convErr := strconv.Atoi(arg); convErr == nil {
		s, err = app.Species.FindByGbifKey(ctx, key)
	} else if id, parseErr := uuid.Parse(arg); parseErr == nil {
		s, err = app.Species.Find(ctx, id)
	} else {
		return findSpeciesByPrefix(ctx, app, arg)
	}

	if errors.Is(err, species.ErrNotFound) {
		return s, fmt.Errorf("no species matches %q", arg)
	}
	return s, err
}

// findSpeciesByPrefix resolves a unique prefix of a species' UUID.
func findSpeciesByPrefix(ctx context.Context, app *app, prefix string) (species.Species, error) {
	all, err := allSpecies(ctx, app)
	if err != nil {
		return species.Species{}, err
	}
	return lookup.ByIDPrefix(all, prefix, "species", "species", func(s species.Species) uuid.UUID {
		return s.ID
	})
}

// allSpecies lists every saved species.
func allSpecies(ctx context.Context, app *app) ([]species.Species, error) {
	all, err := app.Species.List(ctx)
	if err != nil {
		return nil, fmt.Errorf("listing species: %w", err)
	}
	return all, nil
}
