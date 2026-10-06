# Filter

Chips and a search box that narrow the list in the browser: tick an animal or a country, type a
word, and the sightings that don't match are hidden, with a count of those that do. Nothing is
fetched or re-rendered, so the page is the same page that was served, every sighting still in it.
Plain JavaScript, no library.

| File | Copy to |
|------|---------|
| `filter.css` | `web/assets/css/filter.css` |
| `filter.js` | `web/assets/js/filter.js` |
| the component below | `web/views/components/filter.templ` |
| the helpers below | `web/views/components/filter.go` |

## Component

```templ
package components

import "workshop/internal/domain/sighting"

type FilterProps struct {
	Target string        // The id of the list whose items are filtered.
	Search bool          // A search box too, matching each item's data-search.
	Groups []FilterGroup // A row of chips each.
}

type FilterGroup struct {
	Key     string           // The data attribute on each item: "animal" reads data-animal. One lowercase word.
	Label   string
	Options []sighting.Tally // Each name and how many items carry it, from FilterOptions.
}

templ Filter(props FilterProps) {
	<form class="filter" data-filter={ props.Target } role="search" hidden>
		if props.Search {
			<label class="filter__search">
				<span class="visually-hidden">Search sightings</span>
				<input type="search" name="q" placeholder="Search…" autocomplete="off"/>
			</label>
		}
		for _, g := range props.Groups {
			<fieldset class="filter__group">
				<legend>{ g.Label }</legend>
				for _, o := range g.Options {
					<label class="chip">
						<input type="checkbox" name={ g.Key } value={ o.Name }/>
						{ o.Name } <small>{ o.Count }</small>
					</label>
				}
			</fieldset>
		}
		<p class="filter__status">
			<output class="filter__count" aria-live="polite"></output>
			<button type="reset">Clear</button>
		</p>
	</form>
}
```

## Helpers

The chips and the items are labelled by the same functions, so a chip's value always matches
what it filters:

```go
package components

import (
	"strconv"
	"strings"

	"github.com/a-h/templ"

	"workshop/internal/domain/sighting"
)

// Animal names the animal seen, as the card does.
func Animal(s sighting.Sighting) string { return displayName(s.Species) }

// Country is where it was seen.
func Country(s sighting.Sighting) string { return s.Location.Country }

// Year is when, to the year.
func Year(s sighting.Sighting) string { return strconv.Itoa(s.HappenedAt.Year()) }

// FilterOptions counts how often each key appears, keeping the n most common; 0 keeps all.
func FilterOptions(sightings []sighting.Sighting, n int, key func(sighting.Sighting) string) []sighting.Tally {
	return sighting.Top(sighting.CountBy(sightings, key), n)
}

// FilterAttrs marks a list item with what the filter matches on.
func FilterAttrs(s sighting.Sighting) templ.Attributes {
	return templ.Attributes{
		"data-animal":  Animal(s),
		"data-country": Country(s),
		"data-year":    Year(s),
		"data-search": strings.ToLower(strings.Join([]string{
			Animal(s), s.Species.ScientificName, place(s.Location), s.RecordedBy,
		}, " ")),
	}
}
```

## Using it

Above the list, with the list given the id the form points at and each item its attributes:

```templ
@components.Filter(components.FilterProps{
	Target: "sightings",
	Search: true,
	Groups: []components.FilterGroup{
		{Key: "animal", Label: "Animal", Options: components.FilterOptions(sightings, 12, components.Animal)},
		{Key: "country", Label: "Country", Options: components.FilterOptions(sightings, 12, components.Country)},
	},
})
<ul id="sightings" class="sightings">
	for _, s := range sightings {
		<li { components.FilterAttrs(s)... }>
			@components.SightingCard(s)
		</li>
	}
</ul>
```

In the layout's `<head>`, and at the end of `<body>`:

```templ
<link rel="stylesheet" href="/assets/css/filter.css"/>
```

```templ
<script src="/assets/js/filter.js" defer></script>
```

## Gotchas

1. **It hides, it never removes.** The HTML still carries every sighting, so what the Profiling
   block measures is unchanged, and so is the page without JavaScript: the form stays `hidden` and
   everything shows.
2. **`[hidden]` needs putting back on top.** Any `display` rule on the items, such as a grid, beats
   the attribute, so `filter.css` restores it with `!important`. Items that won't hide have lost
   that rule.
3. **Chips and items are labelled by the same function.** `FilterOptions(…, components.Animal)`
   and `FilterAttrs` both call `Animal`, so the values match exactly. A new group needs a key
   function used in both places, and a `Key` of one lowercase word, since `dataset` turns
   `data-sky-code` into `skyCode`.
4. **Cap the chips.** Eighty-nine countries is a wall; `12` keeps the common ones and the search box
   finds the rest. `0` keeps all.
5. **Nothing is written as HTML.** templ escapes the attributes, the script compares strings and
   sets the count with `textContent`. Keep it that way.
6. **It reads the items once.** Thousands of `dataset` reads per keystroke would drag, so
   `filter.js` copies them when it starts, and items added later aren't seen.
7. **Tables work too.** `Target` is the `<tbody>`'s id and each `<tr>` takes `FilterAttrs`.

## Common edits

| Request | Change |
|---------|--------|
| Decades rather than years | a `Decade` key, `strconv.Itoa(s.HappenedAt.Year()/10*10) + "s"`, used in a group and in `FilterAttrs` |
| One at a time | `type="radio"`; the script doesn't mind |
| Remember the filter in the URL | in `apply`, `history.replaceState(null, '', '?' + new URLSearchParams(new FormData(form)))`, and set the fields from `location.search` on load |
| "No sightings match" | a `<p hidden>` after the list, shown in `apply` when `shown === 0` |
| Filter the map as well | dispatch a `CustomEvent('filter:change', { detail: shownAnchors })` from `apply`, and have `map.js` listen and `setStyle` its markers |
| Filter by weather | `"data-sky": sky(s.Weather.Condition.Code)` in `FilterAttrs`, from the weather block |
