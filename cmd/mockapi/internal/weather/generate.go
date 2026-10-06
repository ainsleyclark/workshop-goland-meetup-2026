package weather

import (
	"fmt"
	"hash/fnv"
	"math"
	"math/rand/v2"
	"time"
)

// Open-Meteo's real weather codes for the subset this mock generates. See:
// https://open-meteo.com/en/docs/historical-weather-api
//
// These are deliberately not internal/clients/openmeteo.WeatherCode: that
// type numbers a different, manned-station WMO table whose descriptions
// don't match what these same integers mean in Open-Meteo's own output
// (its code 3 is "CloudsForming"; Open-Meteo's weather_code 3 means
// "Overcast").
const (
	codeClearSky     = 0
	codeMainlyClear  = 1
	codePartlyCloudy = 2
	codeOvercast     = 3

	codeDrizzle      = 51
	codeRainSlight   = 61
	codeRainModerate = 63
	codeRainHeavy    = 65

	codeSnowSlight   = 71
	codeSnowModerate = 73
	codeSnowHeavy    = 75

	codeShowersSlight   = 80
	codeShowersModerate = 81
	codeShowersViolent  = 82

	codeThunderstorm           = 95
	codeThunderstormSlightHail = 96
	codeThunderstormHeavyHail  = 99
)

// reading is one hour of generated weather, shaped for hourlyValues.
type reading struct {
	Temperature, ApparentTemperature float64
	Humidity                         int
	Precipitation, Rain, Snowfall    float64
	CloudCover                       int
	Code                             int
	WindSpeed, WindGusts             float64
	WindDirection                    int
}

// generate deterministically derives one hour of weather for a place and
// time, so the same request always returns the same reading.
func generate(lat, lon float64, hour time.Time) reading {
	r := rand.New(seed(lat, lon, hour)) // #nosec G404 -- deterministic mock data, not a security use.

	cloudRoll := r.Float64()
	precipRoll := r.Float64()
	windRoll := r.Float64()
	directionRoll := r.Float64()
	jitterRoll := r.Float64()

	temp := temperature(lat, hour, jitterRoll)
	code, cloudCover, precipitation, rain, snowfall := weatherFor(temp, lat, cloudRoll, precipRoll)

	precipitating := precipitation > 0
	thunder := code == codeThunderstorm || code == codeThunderstormSlightHail || code == codeThunderstormHeavyHail

	speed, gusts, direction := windFor(lat, precipitating, thunder, windRoll, directionRoll)
	hum := humidity(cloudCover, precipitating)

	return reading{
		Temperature:         temp,
		ApparentTemperature: apparentTemperature(temp, speed, hum),
		Humidity:            hum,
		Precipitation:       precipitation,
		Rain:                rain,
		Snowfall:            snowfall,
		CloudCover:          cloudCover,
		Code:                code,
		WindSpeed:           speed,
		WindGusts:           gusts,
		WindDirection:       direction,
	}
}

// seed derives a PRNG source from the request, so the same place and hour
// always roll the same sequence. hash/fnv is used rather than hash/maphash,
// whose seed is randomised per process and would break that guarantee
// across runs.
func seed(lat, lon float64, hour time.Time) rand.Source {
	key := fmt.Sprintf("%.6f,%.6f,%s", lat, lon, hour.UTC().Format(hourLayout))

	h1 := fnv.New64a()
	_, _ = h1.Write([]byte(key))

	h2 := fnv.New64a()
	_, _ = h2.Write([]byte(key + "#2"))

	return rand.NewPCG(h1.Sum64(), h2.Sum64())
}

// temperature derives a plausible temperature in degrees Celsius from
// latitude (colder toward the poles), season (Southern Hemisphere-aware,
// since the workshop's Kenya, humpback and Antarctica themes are all south
// of or near the equator), and a small diurnal wobble and jitter roll for
// variety between nearby points.
func temperature(lat float64, hour time.Time, jitterRoll float64) float64 {
	latFactor := math.Abs(lat) / 90

	base := 27.0 - 65.0*math.Pow(latFactor, 1.5) // ~27C at the equator, ~-38C at the poles.
	amplitude := 4.0 + 16.0*latFactor            // small swing near the equator, up to 20C at the poles.

	warmestDay := 197.0 // mid-July: Northern Hemisphere summer.
	if lat <= 0 {
		warmestDay = 15.0 // mid-January: Southern Hemisphere summer.
	}
	seasonal := amplitude * math.Cos(2*math.Pi*(float64(hour.YearDay())-warmestDay)/365.25)

	diurnal := 3.0 * math.Cos(2*math.Pi*(float64(hour.UTC().Hour())-15.0)/24.0)
	jitter := jitterRoll*4.0 - 2.0

	return round1(base + seasonal + diurnal + jitter)
}

// weatherFor picks an internally consistent weather code and cloud cover,
// precipitation, rain and snowfall for a temperature and latitude: it never
// returns snowfall above freezing, or a thunderstorm code in the polar
// latitudes.
func weatherFor(temp, lat, cloudRoll, precipRoll float64) (code, cloudCover int, precipitation, rain, snowfall float64) {
	cloudCover = int(math.Round(cloudRoll * 100))

	switch {
	case cloudCover < 15:
		code = codeClearSky
	case cloudCover <= 40:
		code = codeMainlyClear
	case cloudCover <= 70:
		code = codePartlyCloudy
	default:
		code = codeOvercast
	}

	// Precipitation only happens above 40% cloud cover.
	if cloudCover <= 40 {
		return code, cloudCover, 0, 0, 0
	}
	chance := float64(cloudCover-40) / 60 * 0.75
	if precipRoll >= chance {
		return code, cloudCover, 0, 0, 0
	}

	amount := round1(0.1 + (chance-precipRoll)/chance*9.9)
	freezing := temp <= 0
	polar := math.Abs(lat) > 60

	switch {
	case freezing:
		snowfall = amount
		switch {
		case amount < 2:
			code = codeSnowSlight
		case amount < 6:
			code = codeSnowModerate
		default:
			code = codeSnowHeavy
		}
	case !polar && temp >= 18 && amount > 7:
		rain = amount
		switch {
		case amount > 9:
			code = codeThunderstormHeavyHail
		case amount > 8:
			code = codeThunderstormSlightHail
		default:
			code = codeThunderstorm
		}
	case cloudCover < 85:
		rain = amount
		switch {
		case amount < 2:
			code = codeShowersSlight
		case amount < 6:
			code = codeShowersModerate
		default:
			code = codeShowersViolent
		}
	default:
		rain = amount
		switch {
		case amount < 1:
			code = codeDrizzle
		case amount < 4:
			code = codeRainSlight
		case amount < 7:
			code = codeRainModerate
		default:
			code = codeRainHeavy
		}
	}

	return code, cloudCover, rain + snowfall, rain, snowfall
}

// windFor derives wind speed, gusts and direction in km/h and degrees. Wind
// strengthens toward the poles (a nod to the Southern Ocean and Antarctica)
// and under precipitation or thunder.
func windFor(lat float64, precipitating, thunder bool, windRoll, directionRoll float64) (speed, gusts float64, direction int) {
	latFactor := math.Abs(lat) / 90
	speed = (3.0 + windRoll*22.0) * (1.0 + latFactor*0.6)

	switch {
	case thunder:
		speed += 15
	case precipitating:
		speed += 5
	}

	gusts = speed * (1.3 + windRoll*0.4)
	direction = int(directionRoll*360) % 360

	return round1(speed), round1(gusts), direction
}

// humidity derives a relative humidity percentage from cloud cover, raised
// further when it's precipitating.
func humidity(cloudCover int, precipitating bool) int {
	h := 30 + cloudCover*40/100
	if precipitating {
		h += 20
	}
	return min(max(h, 0), 100)
}

// apparentTemperature approximates how the temperature feels: wind cools
// it, humidity nudges it. This is a simplification for teaching material,
// not a physical model.
func apparentTemperature(temp, windSpeed float64, humidity int) float64 {
	return round1(temp - windSpeed*0.05 + (float64(humidity)-50)/100)
}

// round1 rounds to one decimal place, matching the precision Open-Meteo's
// real API returns.
func round1(v float64) float64 {
	return math.Round(v*10) / 10
}
