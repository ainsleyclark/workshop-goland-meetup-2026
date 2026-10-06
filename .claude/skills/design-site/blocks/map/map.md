# Map

The sightings on a map: [Leaflet](https://leafletjs.com) 1.9 from a CDN on OpenStreetMap tiles,
which need no API key. Markers are drawn on a canvas, so the whole list can go on one map, and
each opens a popup that jumps to its record on the page. Two optional plugins group nearby markers
into clusters or draw where sightings are densest as heat.

The map is an extra: the same records are always in the page as text, so without JavaScript or
the CDN the reader loses the map and nothing else.

| File | Copy to |
|------|---------|
| `map.css` | `web/assets/css/map.css` |
| `map.js` | `web/assets/js/map.js` |
| the component below | `web/views/components/map.templ` |
| the helpers below | `web/views/components/map.go` |

## Component

```templ
package components

import "cmp"

type MapProps struct {
	Points  []MapPoint
	ID      string // The id of the JSON the map reads; "map-points" when empty, so set it when a page has two maps.
	Cluster bool   // Group nearby markers, with Leaflet.markercluster.
	Heat    bool   // Draw how dense the sightings are instead of each one, with Leaflet.heat.
	Label   string // What the map shows, for screen readers.
	Class   string
}

templ Map(props MapProps) {
	<div
		class={ "map", props.Class }
		data-map={ cmp.Or(props.ID, "map-points") }
		data-cluster?={ props.Cluster }
		data-heat?={ props.Heat }
		role="region"
		aria-label={ cmp.Or(props.Label, "Map of sightings") }
	>
		<p class="map__fallback">The map needs JavaScript and a connection to load.</p>
	</div>
	@templ.JSONScript(cmp.Or(props.ID, "map-points"), props.Points)
}
```

## Helpers

```go
package components

import "workshop/internal/domain/sighting"

// MapPoint is one marker, handed to map.js as JSON.
type MapPoint struct {
	Lat    float64 `json:"lat"`
	Lng    float64 `json:"lng"`
	Title  string  `json:"title"`
	Date   string  `json:"date"`
	Place  string  `json:"place"`
	Count  string  `json:"count,omitzero"`
	Anchor string  `json:"anchor"`
}

// MapPoints turns sightings into markers, leaving out any with no position.
func MapPoints(sightings []sighting.Sighting) []MapPoint {
	points := make([]MapPoint, 0, len(sightings))
	for _, s := range sightings {
		if s.Location.Latitude == 0 && s.Location.Longitude == 0 {
			continue
		}
		points = append(points, MapPoint{
			Lat:    s.Location.Latitude,
			Lng:    s.Location.Longitude,
			Title:  displayName(s.Species),
			Date:   s.HappenedAt.Format("2 Jan 2006"),
			Place:  place(s.Location),
			Count:  count(s.IndividualCount),
			Anchor: Anchor(s),
		})
	}
	return points
}

// Anchor is a stable element id for a sighting, so a marker can link to its card.
func Anchor(s sighting.Sighting) string {
	return "sighting-" + s.ID.String()
}

// count renders how many were seen, or nothing when GBIF didn't say.
func count(n *int) string {
	if n == nil {
		return ""
	}
	return individuals(*n)
}
```

`displayName`, `place` and `individuals` are already in `components`, beside the sighting card.

## Using it

```templ
@components.Map(components.MapProps{
	Points:  components.MapPoints(sightings),
	Cluster: true,
	Label:   "Where the sightings were",
})
```

Give each record the id the popup links to, on the card or the row: `id={ components.Anchor(s) }`.

In the layout's `<head>`:

```templ
<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/leaflet.min.css" integrity="sha512-h9FcoyWjHcOcmEVkxOfTLnmZFWIH0iZhZT1H2TbOq55xssQGEJHEaIm+PgoUaZbRvQTNTluNOEfb1ZRy6D3BOw==" crossorigin="anonymous" referrerpolicy="no-referrer"/>
<link rel="stylesheet" href="/assets/css/map.css"/>
```

At the end of `<body>`, in this order. `defer` keeps them in order, so Leaflet is there before
`map.js` runs:

```templ
<script src="https://cdnjs.cloudflare.com/ajax/libs/leaflet/1.9.4/leaflet.min.js" integrity="sha512-puJW3E/qXDqYp9IfhAI54BJEaWIfloJ7JWs7OeD5i6ruC9JZL1gERT1wjtwXFlh7CjE7ZJ+/vcRZRkIYIb6p4g==" crossorigin="anonymous" referrerpolicy="no-referrer" defer></script>
<script src="/assets/js/map.js" defer></script>
```

For `Cluster`, add the plugin between those two scripts, and its two stylesheets after Leaflet's:

```templ
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/leaflet.markercluster@1.5.3/dist/MarkerCluster.css" integrity="sha384-pmjIAcz2bAn0xukfxADbZIb3t8oRT9Sv0rvO+BR5Csr6Dhqq+nZs59P0pPKQJkEV" crossorigin="anonymous"/>
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/leaflet.markercluster@1.5.3/dist/MarkerCluster.Default.css" integrity="sha384-wgw+aLYNQ7dlhK47ZPK7FRACiq7ROZwgFNg0m04avm4CaXS+Z9Y7nMu8yNjBKYC+" crossorigin="anonymous"/>
```

```templ
<script src="https://cdn.jsdelivr.net/npm/leaflet.markercluster@1.5.3/dist/leaflet.markercluster.js" integrity="sha384-eXVCORTRlv4FUUgS/xmOyr66XBVraen8ATNLMESp92FKXLAMiKkerixTiBvXriZr" crossorigin="anonymous" defer></script>
```

For `Heat`, this one instead, in the same place:

```templ
<script src="https://cdn.jsdelivr.net/npm/leaflet.heat@0.2.0/dist/leaflet-heat.js" integrity="sha384-mFKkGiGvT5vo1fEyGCD3hshDdKmW3wzXW/x+fWriYJArD0R3gawT6lMvLboM22c0" crossorigin="anonymous" defer></script>
```

Then set `--map-marker` and `--map-marker-stroke` in the site's CSS: `map.js` reads them, so the
markers and clusters take the palette.

## Gotchas

1. **Thousands of points need clusters or heat.** Markers always go on a canvas, which copes with
   the whole list, but seven thousand dots on top of each other tell the reader nothing. Cluster
   them, or draw heat. A page that lists a few hundred, one country or one animal, can show plain
   markers.
2. **The points are in the page as JSON**, about a hundred bytes each, so the full list adds
   close to a megabyte. Drop the fields the popup doesn't need from `MapPoint`, or map a page with
   a `ListFilter`.
3. **Popups are built from DOM nodes with `textContent`**, never from an HTML string, because every
   field in them is GBIF's. Keep it that way when changing what a popup shows.
4. **Restart `make web` after adding the component.** The templ is compiled in, so the page has no
   map until the server restarts. `map.js` and `map.css` are read from disk in development.
5. **The box needs a size before Leaflet starts**, which `map.css` gives it with `aspect-ratio` and
   `min-height`. A map that's a grey strip or nothing at all has lost that.
6. **The map isn't the accessible version; the list is.** Canvas markers can't take keyboard focus,
   so the records stay in the page as text and the map's `Label` says what it shows.
7. **Worldwide points show the whole world.** `fitBounds` frames whatever it's given, and a tight
   set zooms in to `data-max-zoom` (12 by default).
8. **The attribution stays.** OpenStreetMap's tiles are free on the condition that they're credited.

## Common edits

| Request | Change |
|---------|--------|
| Heat instead of markers | `Heat: true` and the heat script in place of the cluster one |
| Marker colour | `--map-marker` and `--map-marker-stroke` in the site's CSS |
| A dark map | `Class: "map--dark"`, which inverts the tiles in `map.css` |
| Taller or shorter | `--map-height: 70vh` on `.map`, instead of the aspect ratio |
| Another basemap | `data-tiles` and `data-attribution` on the `<div>`; `map.js` reads them. OpenStreetMap lists keyless ones |
| Two maps on a page | a different `ID` for each |
| Start zoomed out | `data-max-zoom="4"` on the `<div>` |
