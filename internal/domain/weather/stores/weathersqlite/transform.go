package weathersqlite

import (
	"uuid"
	"workshop/internal/domain/weather"
	"workshop/internal/infra/db/sqlc"
)

func transform(r db.Weather) weather.Weather {
	return weather.Weather{
		ID:         r.ID,
		ObservedAt: r.ObservedAt,
		Temperature: weather.Temperature{
			Actual:   r.Temperature,
			Apparent: r.ApparentTemperature,
		},
		Precipitation: weather.Precipitation{
			Total:    r.Precipitation,
			Rain:     r.Rain,
			Snowfall: r.Snowfall,
		},
		Wind: weather.Wind{
			Speed:     r.WindSpeed,
			Gusts:     r.WindGusts,
			Direction: int(r.WindDirection),
		},
		Condition: weather.Condition{
			Code:        int(r.ConditionCode),
			Description: r.Condition,
		},
		CloudCover: int(r.CloudCover),
		Humidity:   int(r.Humidity),
		CreatedAt:  r.CreatedAt,
		UpdatedAt:  r.UpdatedAt,
	}
}

func toWeatherCreateParams(id uuid.UUID, in weather.CreateParams) db.WeatherCreateParams {
	return db.WeatherCreateParams{
		ID:                  id,
		ObservedAt:          in.ObservedAt,
		Temperature:         in.Temperature.Actual,
		ApparentTemperature: in.Temperature.Apparent,
		Humidity:            int64(in.Humidity),
		Precipitation:       in.Precipitation.Total,
		Rain:                in.Precipitation.Rain,
		Snowfall:            in.Precipitation.Snowfall,
		CloudCover:          int64(in.CloudCover),
		WindSpeed:           in.Wind.Speed,
		WindGusts:           in.Wind.Gusts,
		WindDirection:       int64(in.Wind.Direction),
		ConditionCode:       int64(in.Condition.Code),
		Condition:           in.Condition.Description,
	}
}
