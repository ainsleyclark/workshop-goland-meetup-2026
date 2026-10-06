package sighting

import (
	"context"
	"uuid"
	"workshop/internal/clients/gbif"
	"workshop/internal/clients/openmeteo"
	"workshop/internal/domain/weather"
)

// weatherPerSighting records the conditions where and when it happened.
func (s Service) weatherPerSighting(ctx context.Context, v gbif.Occurrence) (uuid.UUID, error) {
	reading, err := s.openmeteo.Reading(ctx, openmeteo.Request{
		Coordinates: openmeteo.Coordinates{
			Latitude:  v.DecimalLatitude,
			Longitude: v.DecimalLongitude,
		},
		Time: v.EventDate,
	})
	if err != nil {
		return uuid.Nil(), err
	}

	created, err := s.weather.Create(ctx, weather.CreateParams{
		ObservedAt: reading.Time,
		Temperature: weather.Temperature{
			Actual:   reading.Temperature,
			Apparent: reading.ApparentTemperature,
		},
		Precipitation: weather.Precipitation{
			Total:    reading.Precipitation,
			Rain:     reading.Rain,
			Snowfall: reading.Snowfall,
		},
		Wind: weather.Wind{
			Speed:     reading.WindSpeed,
			Gusts:     reading.WindGusts,
			Direction: reading.WindDirection,
		},
		Condition: weather.Condition{
			Code:        int(reading.Code),
			Description: reading.Code.String(),
		},
		CloudCover: reading.CloudCover,
		Humidity:   reading.Humidity,
	})
	if err != nil {
		return uuid.Nil(), err
	}

	return created.ID, nil
}
