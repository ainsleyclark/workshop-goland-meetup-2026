# Weather

The conditions a sighting was seen in, from the weather already saved beside it: an icon for the
sky, the temperature, Open-Meteo's description in words, and the wind as a small arrow. Plain
templ and CSS, no library.

Weather is on every sighting once Block 2's ingest has run. Until then it's zero, and the
component says so rather than hiding the sighting.

| File | Copy to |
|------|---------|
| `weather.css` | `web/assets/css/weather.css` |
| the components below | `web/views/components/weather.templ` |
| the helpers below | `web/views/components/weather.go` |

## Components

```templ
package components

import "workshop/internal/domain/weather"

templ Weather(w weather.Weather) {
	if w.ObservedAt.IsZero() {
		<p class="weather weather--none">No weather recorded</p>
	} else {
		<div class={ "weather", "weather--" + sky(w.Condition.Code) }>
			<svg class="weather__icon" aria-hidden="true"><use href={ templ.URL("#wx-" + sky(w.Condition.Code)) }></use></svg>
			<span class="weather__temp">{ degrees(w.Temperature.Actual) }</span>
			<span class="weather__condition">{ w.Condition.Description }</span>
			<span class={ "weather__wind", "weather__wind--" + compass(w.Wind.Direction) } aria-label={ windLabel(w.Wind) }>
				<span class="weather__arrow" aria-hidden="true"></span>
				{ kmh(w.Wind.Speed) }
			</span>
		</div>
	}
}

templ WeatherIcons() {
	<svg class="weather-icons" aria-hidden="true" focusable="false">
		<symbol id="wx-clear" viewBox="0 0 24 24"><circle cx="12" cy="12" r="5"></circle><path d="M12 1v2M12 21v2M4.22 4.22l1.42 1.42M18.36 18.36l1.42 1.42M1 12h2M21 12h2M4.22 19.78l1.42-1.42M18.36 5.64l1.42-1.42"></path></symbol>
		<symbol id="wx-cloudy" viewBox="0 0 24 24"><path d="M18 10h-1.26A8 8 0 1 0 9 20h9a5 5 0 0 0 0-10z"></path></symbol>
		<symbol id="wx-fog" viewBox="0 0 24 24"><path d="M20 14.58A5 5 0 0 0 18 5h-1.26A8 8 0 1 0 4 13.25M4 18h16M7 22h10"></path></symbol>
		<symbol id="wx-drizzle" viewBox="0 0 24 24"><path d="M8 19v2M8 13v2M16 19v2M16 13v2M12 21v2M12 15v2M20 16.58A5 5 0 0 0 18 7h-1.26A8 8 0 1 0 4 15.25"></path></symbol>
		<symbol id="wx-rain" viewBox="0 0 24 24"><path d="M16 13v8M8 13v8M12 15v8M20 16.58A5 5 0 0 0 18 7h-1.26A8 8 0 1 0 4 15.25"></path></symbol>
		<symbol id="wx-snow" viewBox="0 0 24 24"><path d="M20 17.58A5 5 0 0 0 18 8h-1.26A8 8 0 1 0 4 16.25M8 16h.01M8 20h.01M12 18h.01M12 22h.01M16 16h.01M16 20h.01"></path></symbol>
		<symbol id="wx-storm" viewBox="0 0 24 24"><path d="M19 16.9A5 5 0 0 0 18 7h-1.26a8 8 0 1 0-11.62 9M13 11l-4 6h6l-4 6"></path></symbol>
		<symbol id="wx-dust" viewBox="0 0 24 24"><path d="M9.59 4.59A2 2 0 1 1 11 8H2m10.59 11.41A2 2 0 1 0 14 16H2m15.73-8.27A2.5 2.5 0 1 1 19.5 12H2"></path></symbol>
	</svg>
}
```

The icon paths are from [Feather](https://feathericons.com) (MIT), drawn as strokes so they take
the text colour.

## Helpers

```go
package components

import (
	"fmt"

	"workshop/internal/domain/weather"
)

var winds = [...]struct{ class, name string }{
	{"n", "north"}, {"ne", "north-east"}, {"e", "east"}, {"se", "south-east"},
	{"s", "south"}, {"sw", "south-west"}, {"w", "west"}, {"nw", "north-west"},
}

// sky groups a WMO 4677 weather code into the icon that stands for it.
func sky(code int) string {
	switch {
	case code == 0:
		return "clear"
	case code <= 3:
		return "cloudy"
	case code <= 9, code >= 30 && code <= 35:
		return "dust"
	case code <= 12, code == 28, code >= 40 && code <= 49:
		return "fog"
	case code <= 19, code == 29, code >= 91 && code <= 99:
		return "storm"
	case code == 20, code >= 50 && code <= 59:
		return "drizzle"
	case code == 22, code == 26, code >= 36 && code <= 39, code >= 70 && code <= 79, code == 85, code == 86:
		return "snow"
	case code <= 27, code >= 60 && code <= 69, code >= 80 && code <= 90:
		return "rain"
	default:
		return "cloudy"
	}
}

// degrees renders a temperature to the nearest whole degree.
func degrees(c float64) string {
	return fmt.Sprintf("%.0f °C", c)
}

// kmh renders a wind speed to the nearest whole km/h.
func kmh(v float64) string {
	return fmt.Sprintf("%.0f km/h", v)
}

// point picks the eight-point compass entry for a bearing in degrees.
func point(deg int) int {
	return (((deg%360)+360)%360 + 22) / 45 % 8
}

// compass names the direction the wind blew from, as a class suffix.
func compass(deg int) string {
	return winds[point(deg)].class
}

// windLabel spells the wind out for screen readers.
func windLabel(w weather.Wind) string {
	return fmt.Sprintf("Wind %s from the %s", kmh(w.Speed), winds[point(w.Direction)].name)
}
```

## Using it

The sprite goes in the layout once, just inside `<body>`, and the stylesheet in `<head>`:

```templ
<link rel="stylesheet" href="/assets/css/weather.css"/>
```

```templ
<body>
	@components.WeatherIcons()
	…
```

Then the readout wherever a sighting is shown, such as the card's `<dl>`:

```templ
<dt>Weather</dt>
<dd>@Weather(s.Weather)</dd>
```

## Gotchas

1. **The sprite once, a `<use>` per card.** Each card references a symbol by id instead of
   carrying eight paths, which matters when thousands of cards render. If every icon is blank,
   `WeatherIcons` isn't in the layout.
2. **Zero weather shows as "No weather recorded".** Before Block 2 is finished, or in a copy that
   hasn't ingested, every card says it. Never hide the sighting because its weather is missing.
3. **`Condition.Description` is already in words.** It was translated when the record was saved, so
   show it as it is. `sky()` only groups codes to pick an icon, and it covers the whole WMO table
   rather than Open-Meteo's short list, because the stored data spans it.
4. **The direction is where the wind came *from*.** The class says so (`weather__wind--ne`), and the
   arrow points where the wind was going, as a weather map does. The label spells it out.
5. `Temperature.Apparent` is what it felt like. `Precipitation.Rain` is in millimetres and
   `Snowfall` in centimetres; `Humidity` and `CloudCover` are percentages.

## Common edits

| Request | Change |
|---------|--------|
| Feels like | `<span>feels like { degrees(w.Temperature.Apparent) }</span>` |
| Humidity or cloud cover | `<span>{ strconv.Itoa(w.Humidity) }% humidity</span>` |
| Icon only | drop the spans and put `role="img" aria-label={ w.Condition.Description }` on the `<svg>` |
| A colour per sky | `.weather--clear .weather__icon { color: #e0a12b }`, and so on for `rain`, `storm`, `snow` |
| Filter by weather | `"data-sky": sky(s.Weather.Condition.Code)` in the filter block's `FilterAttrs` |
