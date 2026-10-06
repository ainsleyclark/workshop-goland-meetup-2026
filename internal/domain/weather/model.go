package weather

import (
	"time"
	"uuid"
)

type (
	// Weather records the conditions at one moment and place. Once the
	// hour's weather has been fetched it never changes, so it is saved
	// and read back rather than fetched again. A sighting points at it.
	Weather struct {
		ID            uuid.UUID     `json:"id"`
		ObservedAt    time.Time     `json:"observedAt"`
		Temperature   Temperature   `json:"temperature"`
		Precipitation Precipitation `json:"precipitation"`
		Wind          Wind          `json:"wind"`
		Condition     Condition     `json:"condition"`
		CloudCover    int           `json:"cloudCover"`
		Humidity      int           `json:"humidity"`
		CreatedAt     time.Time     `json:"createdAt"`
		UpdatedAt     time.Time     `json:"updatedAt"`
	}
	// Temperature is the air temperature in degrees Celsius. Apparent is
	// what it felt like, with wind and humidity taken into account.
	Temperature struct {
		Actual   float64 `json:"actual"`
		Apparent float64 `json:"apparent"`
	}
	// Precipitation is what fell during the hour. Total and Rain are in
	// millimetres, Snowfall in centimetres. Rain is the part of Total
	// that fell as rain rather than snow.
	Precipitation struct {
		Total    float64 `json:"total"`
		Rain     float64 `json:"rain"`
		Snowfall float64 `json:"snowfall"`
	}
	// Wind is the wind at ten metres in km/h, and the direction it blew
	// from in degrees clockwise from north.
	Wind struct {
		Speed     float64 `json:"speed"`
		Gusts     float64 `json:"gusts"`
		Direction int     `json:"direction"`
	}
	// Condition describes the present weather. Code is the WMO code the
	// reading came with; Description is that code in words, translated
	// when the record was created so readers need no lookup table.
	Condition struct {
		Code        int    `json:"code"`
		Description string `json:"description"`
	}
	// CreateParams represents the parameters required to record
	// the weather for a sighting.
	CreateParams struct {
		ObservedAt    time.Time     `json:"observedAt"`
		Temperature   Temperature   `json:"temperature"`
		Precipitation Precipitation `json:"precipitation"`
		Wind          Wind          `json:"wind"`
		Condition     Condition     `json:"condition"`
		CloudCover    int           `json:"cloudCover"`
		Humidity      int           `json:"humidity"`
	}
)

// CreateParams returns the weather fields without its ID or timestamps.
func (w Weather) CreateParams() CreateParams {
	return CreateParams{
		ObservedAt:    w.ObservedAt,
		Temperature:   w.Temperature,
		Precipitation: w.Precipitation,
		Wind:          w.Wind,
		Condition:     w.Condition,
		CloudCover:    w.CloudCover,
		Humidity:      w.Humidity,
	}
}
