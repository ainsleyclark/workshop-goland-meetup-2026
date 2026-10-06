# Reveal

Sections fade and rise into place as you scroll to them, once each. Plain CSS plus a small
IntersectionObserver, no library. Adapted from the scroll reveal on needypeanut.com.

| File | Copy to |
|------|---------|
| `reveal.css` | `web/assets/css/reveal.css` |
| `reveal.js` | `web/assets/js/reveal.js` |

## Wiring it up

In the layout's `<head>`: the first line marks the page as having JS before it paints, so there's
no flash of content that then vanishes. Keep all three together, because if the class is set and
`reveal.js` never loads, everything marked `.reveal` stays hidden.

```templ
<script>document.documentElement.classList.add('js')</script>
<link rel="stylesheet" href="/assets/css/reveal.css"/>
<script src="/assets/js/reveal.js" defer></script>
```

Then add `reveal` to what should arrive:

```templ
<section class="stats reveal">…</section>
<h2 class="reveal">Latest photos</h2>
```

To stagger a few siblings, give them a fixed delay: `style="--reveal-delay: 150ms"`, then 300ms,
and so on. Keep those values fixed in the template: never build a style from GBIF data.

## Gotchas

- **Sections and headings, not every card.** The home page renders thousands of sightings, and
  thousands of fading cards is slow and tiring. Reveal the list's heading or section as one.
- **Never inside a swiper slide.** Slides off to the side never scroll into view vertically, so
  they stay invisible until swiped to and then fade in mid-swipe. Reveal the whole section instead.
- Nothing is hidden without JS, and with reduced motion switched on everything simply shows.
- Tune the feel with `--reveal-distance` (default `1.5rem`) and the transition timing in `reveal.css`.
