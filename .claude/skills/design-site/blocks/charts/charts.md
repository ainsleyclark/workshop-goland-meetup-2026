# Charts

Two small charts from the aggregates the domain already has, with no library: a sparkline of
sightings per year, drawn as inline SVG, and bars for the top few of anything, drawn in CSS. Both
are plain markup, so they take the site's colours and they print.

| File | Copy to |
|------|---------|
| `charts.css` | `web/assets/css/charts.css` |
| the components below | `web/views/components/charts.templ` |
| the helpers below | `web/views/components/charts.go` |

## Components

```templ
package components

import "workshop/internal/domain/sighting"

templ Sparkline(byYear map[int]int, sum sighting.Summary) {
	if sum.Sightings > 0 {
		<svg class="sparkline" viewBox="0 0 100 24" preserveAspectRatio="none" role="img" aria-label={ sparkLabel(byYear, sum) }>
			<polygon class="sparkline__fill" points={ sparkArea(byYear, sum) }></polygon>
			<polyline class="sparkline__line" points={ sparkPoints(byYear, sum) } vector-effect="non-scaling-stroke"></polyline>
		</svg>
	}
}

templ Bars(rows []sighting.Tally) {
	if len(rows) > 0 {
		<ol class="bars">
			for _, row := range rows {
				<li class="bars__row" style={ bar(row.Count, rows[0].Count) }>
					<span class="bars__name">{ row.Name }</span>
					<span class="bars__count">{ row.Count }</span>
				</li>
			}
		</ol>
	}
}
```

## Helpers

```go
package components

import (
	"fmt"
	"strconv"
	"strings"

	"workshop/internal/domain/sighting"
)

// The sparkline's drawing box, matching the viewBox in the templ.
const sparkW, sparkH = 100.0, 24.0

// years lists every year from the first sighting to the last, so a quiet year stays a dip.
func years(sum sighting.Summary) []int {
	ys := make([]int, 0, sum.Latest.Year()-sum.Earliest.Year()+1)
	for y := sum.Earliest.Year(); y <= sum.Latest.Year(); y++ {
		ys = append(ys, y)
	}
	return ys
}

// peak finds the busiest year and its count.
func peak(byYear map[int]int) (year, count int) {
	for y, n := range byYear {
		if n > count || (n == count && y < year) {
			year, count = y, n
		}
	}
	return year, count
}

// sparkPoints plots a count per year across the box, the busiest year at the top.
func sparkPoints(byYear map[int]int, sum sighting.Summary) string {
	ys := years(sum)
	_, high := peak(byYear)
	var b strings.Builder
	for i, y := range ys {
		x := 0.0
		if len(ys) > 1 {
			x = float64(i) / float64(len(ys)-1) * sparkW
		}
		top := sparkH
		if high > 0 {
			top = sparkH - float64(byYear[y])/float64(high)*(sparkH-1)
		}
		fmt.Fprintf(&b, "%.1f,%.1f ", x, top)
	}
	return strings.TrimSpace(b.String())
}

// sparkArea closes the line down to the baseline, for the fill under it.
func sparkArea(byYear map[int]int, sum sighting.Summary) string {
	return fmt.Sprintf("0,%v %s %v,%v", sparkH, sparkPoints(byYear, sum), sparkW, sparkH)
}

// sparkLabel says in words what the line shows.
func sparkLabel(byYear map[int]int, sum sighting.Summary) string {
	y, n := peak(byYear)
	return fmt.Sprintf("Sightings per year, %d to %d, most in %d with %d",
		sum.Earliest.Year(), sum.Latest.Year(), y, n)
}

// bar sets --bar to a row's share of the longest, for charts.css to draw.
func bar(count, longest int) map[string]string {
	if longest <= 0 {
		return map[string]string{"--bar": "0%"}
	}
	return map[string]string{"--bar": strconv.Itoa(count*100/longest) + "%"}
}
```

## Using it

Work the numbers out once, at the top of the page in a templ Go block, and hand them to the
charts and to the figures beside them:

```templ
templ Home(sightings []sighting.Sighting) {
	{{ sum := sighting.Summarise(sightings) }}
	{{ byYear := sighting.CountByYear(sightings) }}
	@layout.Base(layout.BaseProps{Title: "Sightings"}) {
		<section class="stats">
			<p>{ sum.Sightings } sightings of { sum.Species } species in { sum.Countries } countries, { sum.Earliest.Year() } to { sum.Latest.Year() }</p>
			@components.Sparkline(byYear, sum)
			<h2>Most sightings by country</h2>
			@components.Bars(topCountries(sightings, 5))
		</section>
		…
```

`topCountries` is a plain Go helper in a file beside the page:

```go
// topCountries ranks the countries with the most sightings.
func topCountries(sightings []sighting.Sighting, n int) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, func(s sighting.Sighting) string {
		return s.Location.Country
	}), n)
}
```

And the stylesheet in the layout's `<head>`:

```templ
<link rel="stylesheet" href="/assets/css/charts.css"/>
```

## Gotchas

1. **Once per page.** `Summarise`, `CountByYear` and `Top` each read every sighting, so call them
   in a `{{ }}` block at the top, never inside the `for` over the sightings.
2. **Nothing renders for no data.** Both components return nothing on an empty page, so a heading
   above them should read right on its own.
3. **The sparkline stretches to fit**, with `preserveAspectRatio="none"`, and
   `vector-effect="non-scaling-stroke"` is what keeps the line the same thickness while it does.
   Keep both.
4. **A picture plus its words.** The `aria-label` says what the line shows; the figures it sums up
   belong beside it as text, as above, not only in the label.
5. **Empty years are dips, not gaps.** The dataset runs from 1879 and most of those years are flat
   zero, which is the truth of it. For a recent-only line, start `years()` later, or chart a page
   with a `ListFilter{FromDate: …}`.
6. **`--bar` is the one inline style**: a percentage we worked out, set through a
   `map[string]string` so templ writes it safely. Never build a style from GBIF's text.

## Common edits

| Request | Change |
|---------|--------|
| Columns per year instead of a line | an `<ol>` of `<li style={ bar(byYear[y], high) }>` over `years(sum)`, with `height: var(--bar)` in CSS |
| Top animals or recorders | a helper like `topCountries` keyed on `displayName(s.Species)` or `s.RecordedBy`; `displayName` lives in `components`, so put that helper there |
| A dot on the latest year | a `<circle>` at the last pair from `sparkPoints`, split with `strings.Fields` |
| Line only, no fill | drop the `<polygon>` |
| Thicker line | `stroke-width` on `.sparkline__line`; it stays that width at any size |
| Share as a percentage | `row.Count * 100 / sum.Sightings` next to the count |
