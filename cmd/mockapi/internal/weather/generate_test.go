package weather

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestTemperature(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		lat  float64
		hour time.Time
	}{
		"Equator, Southern summer": {lat: -1.29, hour: time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)},
		"Kenya, Southern summer":   {lat: -3.4, hour: time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)},
		"Antarctica, Southern summer": {
			lat: -75, hour: time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC),
		},
		"Antarctica, Southern winter": {
			lat: -75, hour: time.Date(2024, time.July, 15, 12, 0, 0, 0, time.UTC),
		},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := temperature(test.lat, test.hour, 0.5)
			assert.InDelta(t, temperature(test.lat, test.hour, 0.5), got, 0, "deterministic for the same inputs")
		})
	}

	t.Run("Equator is warmer than Antarctica at the same time", func(t *testing.T) {
		t.Parallel()
		hour := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
		assert.Greater(t, temperature(-1.29, hour, 0.5), temperature(-75, hour, 0.5))
	})

	t.Run("Antarctica is colder in its winter than its summer", func(t *testing.T) {
		t.Parallel()
		summer := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)
		winter := time.Date(2024, time.July, 15, 12, 0, 0, 0, time.UTC)
		assert.Greater(t, temperature(-75, summer, 0.5), temperature(-75, winter, 0.5))
	})
}

func TestWeatherFor(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		temp, lat, cloudRoll, precipRoll float64
	}{
		"Clear sky, no rolls to precipitate": {temp: 25, lat: -1, cloudRoll: 0.05, precipRoll: 0.9},
		"Overcast, freezing":                 {temp: -5, lat: -75, cloudRoll: 0.95, precipRoll: 0.01},
		"Overcast, tropical heat":            {temp: 30, lat: -1, cloudRoll: 0.95, precipRoll: 0.01},
		"Overcast, polar, warm-code rolls":   {temp: 25, lat: -75, cloudRoll: 0.95, precipRoll: 0.01},
		"Partly cloudy, light rain":          {temp: 20, lat: -1, cloudRoll: 0.6, precipRoll: 0.05},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			code, cloudCover, precipitation, rain, snowfall := weatherFor(test.temp, test.lat, test.cloudRoll, test.precipRoll)

			assert.Equal(t, rain+snowfall, precipitation, "precipitation is always the sum of rain and snowfall")

			if test.temp <= 0 {
				assert.Zero(t, rain, "no rain while freezing")
			} else {
				assert.Zero(t, snowfall, "no snowfall above freezing")
			}

			if code == codeClearSky || code == codeMainlyClear {
				assert.Zero(t, precipitation, "clear and mainly-clear codes never precipitate")
			}

			isThunder := code == codeThunderstorm || code == codeThunderstormSlightHail || code == codeThunderstormHeavyHail
			if isThunder {
				assert.LessOrEqual(t, cloudCover, 100)
				assert.LessOrEqual(t, absFloat(test.lat), 60.0, "no thunderstorms in the polar latitudes")
			}
		})
	}
}

func TestWindFor(t *testing.T) {
	t.Parallel()

	t.Run("Calmer at the equator than the poles for the same rolls", func(t *testing.T) {
		t.Parallel()
		equatorSpeed, _, _ := windFor(-1, false, false, 0.5, 0.5)
		polarSpeed, _, _ := windFor(-89, false, false, 0.5, 0.5)
		assert.Less(t, equatorSpeed, polarSpeed)
	})

	t.Run("Thunder is windier than calm precipitation, which is windier than dry", func(t *testing.T) {
		t.Parallel()
		dry, _, _ := windFor(-1, false, false, 0.5, 0.5)
		wet, _, _ := windFor(-1, true, false, 0.5, 0.5)
		stormy, _, _ := windFor(-1, true, true, 0.5, 0.5)
		assert.Less(t, dry, wet)
		assert.Less(t, wet, stormy)
	})

	t.Run("Gusts are always at least as strong as sustained wind", func(t *testing.T) {
		t.Parallel()
		speed, gusts, _ := windFor(-30, true, true, 0.9, 0.1)
		assert.GreaterOrEqual(t, gusts, speed)
	})

	t.Run("Direction is a compass bearing", func(t *testing.T) {
		t.Parallel()
		_, _, direction := windFor(-1, false, false, 0.5, 0.999)
		assert.GreaterOrEqual(t, direction, 0)
		assert.Less(t, direction, 360)
	})
}

func TestHumidity(t *testing.T) {
	t.Parallel()

	tt := map[string]struct {
		cloudCover    int
		precipitating bool
		want          int
	}{
		"Clear sky":                  {cloudCover: 0, precipitating: false, want: 30},
		"Overcast":                   {cloudCover: 100, precipitating: false, want: 70},
		"Overcast and precipitating": {cloudCover: 100, precipitating: true, want: 90},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			got := humidity(test.cloudCover, test.precipitating)
			assert.Equal(t, test.want, got)
			assert.GreaterOrEqual(t, got, 0)
			assert.LessOrEqual(t, got, 100)
		})
	}
}

func TestApparentTemperature(t *testing.T) {
	t.Parallel()

	t.Run("Wind cools the apparent temperature", func(t *testing.T) {
		t.Parallel()
		still := apparentTemperature(20, 0, 50)
		windy := apparentTemperature(20, 40, 50)
		assert.Less(t, windy, still)
	})

	t.Run("Higher humidity raises the apparent temperature", func(t *testing.T) {
		t.Parallel()
		dry := apparentTemperature(20, 10, 20)
		humid := apparentTemperature(20, 10, 80)
		assert.Less(t, dry, humid)
	})
}

func TestGenerate(t *testing.T) {
	t.Parallel()

	hour := time.Date(2024, time.January, 15, 12, 0, 0, 0, time.UTC)

	t.Run("Deterministic for the same request", func(t *testing.T) {
		t.Parallel()
		first := generate(-3.4, 39.95, hour)
		second := generate(-3.4, 39.95, hour)
		assert.Equal(t, first, second)
	})

	t.Run("Different coordinates give a different reading", func(t *testing.T) {
		t.Parallel()
		kenya := generate(-3.4, 39.95, hour)
		antarctica := generate(-75, 0, hour)
		assert.Greater(t, kenya.Temperature, antarctica.Temperature)
	})
}

func absFloat(v float64) float64 {
	if v < 0 {
		return -v
	}
	return v
}
