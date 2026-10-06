# Swiper

A touch-friendly row of cards. Phones show one card and a bit of the next, tablets show two, and
from 1024px it becomes an ordinary grid. It uses [Swiper](https://swiperjs.com) 11 from a CDN.

Good for a short row of highlights: recent photos, the most-seen animals, one card per country.
Not for the full list of sightings: a slider of thousands is unusable, and each slide adds a dot.

| File | Copy to |
|------|---------|
| `swiper.css` | `web/assets/css/swiper.css` |
| `swiper.js` | `web/assets/js/swiper.js` |
| the component below | `web/views/components/swiper.templ` |

## Component

```templ
package components

import (
	"cmp"
	"strconv"
)

type SwiperProps struct {
	DesktopSlides  int // Columns at 1024px and up; 3 when zero.
	Dark           bool
	HidePagination bool
	Class          string
}

templ Swiper(props SwiperProps) {
	<div
		class={ "swiper", "carousel", templ.KV("carousel--dark", props.Dark), props.Class }
		data-swiper-desktop={ strconv.Itoa(cmp.Or(props.DesktopSlides, 3)) }
	>
		<div class="swiper-wrapper">
			{ children... }
		</div>
		if !props.HidePagination {
			<div class="swiper-pagination"></div>
		}
	</div>
}
```

## Using it

Each child must be a `swiper-slide`. Hand it a short slice, never the full list:

```templ
<section class="featured">
	<div class="container">
		<h2>Latest photos</h2>
	</div>
	@components.Swiper(components.SwiperProps{DesktopSlides: 4}) {
		for _, s := range withPhotos(sightings, 12) {
			<div class="swiper-slide">
				@components.SightingCard(s)
			</div>
		}
	}
</section>
```

`withPhotos` is a plain Go helper beside the `.templ` that keeps the first `n` sightings with a
`StillImage`. On a page with its own handler, `ListFilter{Limit: 12}` does the capping in SQL.

In the layout's `<head>`, in this order. `defer` keeps them in order, so Swiper is loaded before
`swiper.js` runs:

```templ
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/swiper@11.2.10/swiper-bundle.min.css" integrity="sha384-gAPqlBuTCdtVcYt9ocMOYWrnBZ4XSL6q+4eXqwNycOr4iFczhNKtnYhF3NEXJM51" crossorigin="anonymous"/>
<link rel="stylesheet" href="/assets/css/swiper.css"/>
<script src="https://cdn.jsdelivr.net/npm/swiper@11.2.10/swiper-bundle.min.js" integrity="sha384-2UI1PfnXFjVMQ7/ZDEF70CR943oH3v6uZrFQGGqJYlvhh4g6z6uVktxYbOlAczav" crossorigin="anonymous" defer></script>
<script src="/assets/js/swiper.js" defer></script>
```

Then set its tokens in the site's CSS: `--swiper-gutter` (the page's side padding, **in px**),
`--swiper-max` (the content's max width) and `--swiper-theme-color` (the active dot).

## Gotchas

1. **Full bleed needs room.** The slider should be a sibling of the section's `.container`, not
   inside it, so the next card can peek in from the screen edge. If the layout's `<main>` limits
   the width, move that limit onto a `.container` inside each section. If you'd rather keep it
   inside a container, set `--swiper-gutter: 0px` on it and the slides stop at the edge.
2. **No `overflow: hidden` on any ancestor.** It clips slides mid-drag. `swiper.css` already sets
   `overflow-x: clip` on the carousel's parent, which stops sideways page scroll without breaking
   dragging.
3. **`--swiper-gutter` must match the page's side padding**, or the first card sits too far in
   or hangs off the left edge.
4. **Never put a `reveal` class inside a slide.** Off-screen slides never scroll into view, so they
   stay invisible until you swipe to them and then fade in mid-swipe. Reveal the section instead.
5. **Dark section?** Pass `Dark: true`, or the inactive dots are invisible.
6. Without the CDN (offline, blocked), `swiper.js` adds `carousel--static` and the row falls back
   to native scroll-snap. The slider still works, just more plainly.

## Common edits

| Request | Change |
|---------|--------|
| Two columns on desktop, not three | `DesktopSlides: 2` |
| Keep it a slider on desktop too | `swiper.js`: `1024: { slidesPerView: 3, spaceBetween: 24 }`, and delete the `@media (min-width: 1024px)` block in `swiper.css` |
| One whole card on phones | `swiper.js`: `slidesPerView: 1` (loses the peek that shows it swipes) |
| Hide the dots | `HidePagination: true` |
| Loop | `swiper.js`: `loop: true` (makes the dots decorative) |
| Autoplay | `swiper.js`: `autoplay: { delay: 5000, pauseOnMouseEnter: true }`; skip it when `matchMedia('(prefers-reduced-motion: reduce)').matches` |

Ported from needypeanut's `blocks/swiper` (a Hugo partial).
