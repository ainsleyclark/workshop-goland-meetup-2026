# Blocks

Small, ready-made pieces for the site. Copy a block's files into `web/` and they become the
attendee's own code, free to restyle. Nothing links back here. Read a block's `.md` before using
it: its gotchas are the mistakes it's easy to make.

| Block | What it's for | Files | Needs |
|-------|---------------|-------|-------|
| [`favicon`](favicon/favicon.md) | The site's SVG logo, used in the header and as the favicon | `favicon.md` | Nothing |
| [`swiper`](swiper/swiper.md) | A short row of cards: swipes on phones, a grid on desktop | `swiper.md`, `.css`, `.js` | Swiper 11 (CDN) |
| [`lightbox`](lightbox/lightbox.md) | Click a photo to see it full size | `lightbox.md`, `.css`, `.js` | Nothing |
| [`reveal`](reveal/reveal.md) | Sections fade in as you scroll to them | `reveal.md`, `.css`, `.js` | Nothing |
| [`threejs`](threejs/threejs.md) | The plumbing for a 3D scene: you write the scene | `threejs.md`, `.css`, `.js`, `example.js` | three.js r186 (CDN) |
| [`weather`](weather/weather.md) | The conditions each sighting was seen in: an icon, the temperature, the wind | `weather.md`, `.css` | Nothing |
| [`map`](map/map.md) | The sightings on a map, as clusters or heat | `map.md`, `.css`, `.js` | Leaflet 1.9 (CDN); markercluster or heat (CDN) if used |
| [`charts`](charts/charts.md) | A sparkline of sightings per year, and bars for the top few of anything | `charts.md`, `.css` | Nothing |
| [`filter`](filter/filter.md) | Chips and a search box that narrow the list without a reload | `filter.md`, `.css`, `.js` | Nothing |
| [`transitions`](transitions/transitions.md) | A cross-fade between pages, the header staying put | `transitions.md`, `.css` | Nothing |

Each `.md` holds the templ in a code fence, never as a `.templ` file. That keeps `templ generate`
from picking these up here.

`swiper` and `favicon` are ported from needypeanut's blocks, which are Hugo partials for plain
HTML sites. Its contact form isn't here: it sends mail through Resend on Vercel, and this app has
neither.
