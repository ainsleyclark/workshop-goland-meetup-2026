package cmd

import (
	"context"
	"fmt"
	"time"
	"workshop/internal/clients/openmeteo"
	"workshop/internal/cmd/internal/view"
	"workshop/internal/common/printer"

	"github.com/urfave/cli/v3"
)

func weatherCmd(app *app) *cli.Command {
	london := openmeteo.Coordinates{
		Latitude:  51.5072,
		Longitude: -0.1276,
	}

	return &cli.Command{
		Name:  "weather",
		Usage: "Show the current weather in London",
		Action: func(ctx context.Context, _ *cli.Command) error {
			res, err := app.OpenMeteo.Reading(ctx, openmeteo.Request{
				Coordinates: london,
				Time:        time.Now().UTC(),
			})
			if err != nil {
				return fmt.Errorf("fetching weather: %w", err)
			}

			view.Render(app.Printer,
				view.Card{
					Heading: "London Weather",
					Title:   res.Time.Format(view.DateFormat),
					Rows: []printer.KV{
						{Key: "Conditions", Value: res.Code.String()},
						{Key: "Temperature", Value: fmt.Sprintf("%.1f°C (feels like %.1f°C)", res.Temperature, res.ApparentTemperature)},
						{Key: "Humidity", Value: fmt.Sprintf("%d%%", res.Humidity)},
						{Key: "Precipitation", Value: fmt.Sprintf("%.1fmm", res.Precipitation)},
						{Key: "Cloud Cover", Value: fmt.Sprintf("%d%%", res.CloudCover)},
						{Key: "Wind", Value: fmt.Sprintf("%.1f km/h, gusts %.1f km/h, %d°", res.WindSpeed, res.WindGusts, res.WindDirection)},
					},
				},
			)

			return nil
		},
	}
}
