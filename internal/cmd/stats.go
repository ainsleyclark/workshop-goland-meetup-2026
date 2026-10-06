package cmd

import (
	"cmp"
	"context"
	"fmt"
	"strings"
	"time"
	"workshop/internal/cmd/internal/view"
	"workshop/internal/common/printer"
	"workshop/internal/common/printer/styles"
	"workshop/internal/domain/sighting"

	"github.com/urfave/cli/v3"
)

// statsTopN is how many rows each leaderboard shows.
const statsTopN = 5

type (
	// statsReport is a summary of everything in the database.
	statsReport struct {
		Sightings    int
		Species      int
		Countries    int
		Media        int
		Earliest     time.Time
		Latest       time.Time
		TopSpecies   []sighting.Tally
		TopCountries []sighting.Tally
	}
)

func statsCmd(app *app) *cli.Command {
	return &cli.Command{
		Name:  "stats",
		Usage: "Show a summary of the saved sightings and species",
		Action: func(ctx context.Context, _ *cli.Command) error {
			all, err := allSightings(ctx, app)
			if err != nil {
				return err
			}

			r := buildStats(all)
			if r.Sightings == 0 {
				view.Render(app.Printer, view.Heading("Stats"), view.Muted("Nothing to report yet, try 'workshop db seed'"))
				return nil
			}

			view.Render(app.Printer,
				view.Card{Heading: "Stats", Rows: []printer.KV{
					{Key: "Sightings", Value: fmt.Sprint(r.Sightings)},
					{Key: "Species", Value: fmt.Sprint(r.Species)},
					{Key: "Countries", Value: fmt.Sprint(r.Countries)},
					{Key: "Media", Value: fmt.Sprint(r.Media)},
					{Key: "Date range", Value: fmt.Sprintf("%s %s %s",
						r.Earliest.Format(view.DateFormat), styles.IconArrow, r.Latest.Format(view.DateFormat))},
				}},
				leaderboardTable("Top species", "SPECIES", r.TopSpecies),
				leaderboardTable("Top countries", "COUNTRY", r.TopCountries),
			)
			return nil
		},
	}
}

// buildStats summarises sightings, which should already have their
// species filled in. The counting comes from the domain; how a species
// or a country is named on this leaderboard is the CLI's own choice, so
// it supplies those labels.
func buildStats(sightings []sighting.Sighting) statsReport {
	summary := sighting.Summarise(sightings)
	bySpecies := sighting.CountBy(sightings, statsSpeciesName)
	byCountry := sighting.CountBy(sightings, statsCountryName)

	return statsReport{
		Sightings: summary.Sightings,
		// Counted from the labelled tallies rather than the summary, so
		// the totals agree with the leaderboards printed beneath them.
		Species:      len(bySpecies),
		Countries:    len(byCountry),
		Media:        summary.Media,
		Earliest:     summary.Earliest,
		Latest:       summary.Latest,
		TopSpecies:   sighting.Top(bySpecies, statsTopN),
		TopCountries: sighting.Top(byCountry, statsTopN),
	}
}

// statsSpeciesName labels a species the way the rest of the CLI does.
func statsSpeciesName(s sighting.Sighting) string {
	return view.SpeciesName(s.Species)
}

// statsCountryName falls back through what GBIF supplied, so a sighting
// with no country still lands somewhere on the leaderboard.
func statsCountryName(s sighting.Sighting) string {
	return cmp.Or(s.Location.Country, s.Location.CountryCode.String(), "Unknown")
}

// leaderboardTable renders each row with a bar scaled against the top count.
func leaderboardTable(title, header string, rows []sighting.Tally) view.Table[sighting.Tally] {
	const barWidth = 20
	bar := styles.Title.UnsetBold()

	return view.Table[sighting.Tally]{
		Title:   title,
		Headers: []string{header, "SIGHTINGS", ""},
		Items:   rows,
		Row: func(row sighting.Tally) []string {
			width := max(1, row.Count*barWidth/max(1, rows[0].Count))
			return []string{
				row.Name,
				fmt.Sprint(row.Count),
				bar.Render(strings.Repeat("█", width)),
			}
		},
	}
}
