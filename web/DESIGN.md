# Wildlife Watch — Design Document

**Chrome leopard, hot pink.** Modern art × African wild cats × the mid-2000s glam era: pink
flip phones, oversized white shades, tiaras, script logos, lip-gloss pink, animal print and
MySpace sparkles. The site is an exhibition: a hero with a cubist wild-cat mask, a gold
ticker, a Mondrian floor plan of species "rooms", then every sighting in collapsible rooms.
One dark theme, no toggle. Channel the era's vibe without naming or quoting real people.

**Very light, on purpose.** No photos at all: every visual is CSS or inline SVG. The page is
one HTML document (~560 KB for 4,263 sightings, ~50 KB gzipped), one stylesheet, one small
script and three fonts at one or two weights each. Keep it that way: no per-sighting images, no libraries.

## Palette

Gold is the print, hot pink is the accent (it replaced the earlier acid lime). Text on gold,
pink and chrome is always ink (#0B0B0B); white on pink fails. Every pair is checked ≥ 4.5:1;
the chrome type's darkest band (5.7:1) and the pink gloss's (4.3:1) are only used for large
display text.

```css
:root {
  color-scheme: dark;
  --color-background: #0B0B0B;
  --color-surface: #141414;
  --color-surface-2: #1C1C1C;
  --color-foreground: #F2F2F2;
  --color-soft-foreground: #CFCFCF;
  --color-muted-foreground: #9A9A9A;
  --color-gold: #E0A526;     /* leopard gold */
  --color-pink: #FF2E93;     /* hot pink */
  --color-baby-pink: #FFB3D9;
  --color-ink: #0B0B0B;
  --color-border: #2A2A2A;
  --color-ring: #FF2E93;
  --pink-gloss: linear-gradient(180deg, #FFC2E2 0%, #FF4FA3 45%, #E0157A 55%, #FF7DBF 100%);
  --chrome: linear-gradient(180deg, #FFFFFF 0%, #E2E2E2 32%, #8A8A8A 49%, #F7F7F7 53%, #B4B4B4 74%, #EDEDED 100%);
  --chrome-flat: linear-gradient(180deg, #F4F4F4 0%, #CDCDCD 48%, #A9A9A9 52%, #E6E6E6 100%);
  --leopard: url("data:image/svg+xml,…");  /* 60px rosette tile, on --color-gold */
}
```

## Typography

- **Script:** Pacifico, hot pink: the site name (header, hero, footer) and the kickers above
  each title, in sentence case.
- **Display:** Unbounded 800 only, uppercase: titles, room names, counts.
- **Everything else:** Space Mono 400/700: body, labels (uppercase, tracked 0.08–0.14em), tables.
- No italic font is loaded; scientific names use the browser's slanted Space Mono.

## Spacing & radii

- Container 1200px, 20px side padding. Header 64px, ticker 44px.
- `--radius-frame` 14px for the floor plan's chrome frame and rooms; `--radius-pill` for
  buttons and badges; Mondrian tiles are square.
- Background: black with a faint 24px dot grid.

## Components

### Logo
- `assets/images/logo.svg`: the heavy W in chrome on a black disc, with a pink sparkle above.

### Header
- Logo + "Wildlife Watch" in pink Pacifico (hidden under 720px); nav pills in mono: Floor plan →
  `/#floor-plan`, Rooms → `/#sightings`. Chrome rule along the bottom. Solid, no blur.

### Hero (`components.Hero(HeroProps)`, from `components.NewHeroProps`)
- Pink script kicker "Now showing: the wild exhibition"; the title is "Wildlife" in glossy
  pink Pacifico, tilted −6° and overlapping chrome Unbounded "WATCH" below it;
  mono lede ("No photos. No Flash. Just who saw what, and where.").
- Gel buttons, 2000-style: hot pink "Enter the gallery" and chrome "Data: GBIF.org".
- Split-flap hit counter: the sighting total as seven pink digits in a chrome frame.
- `components.CatArt()`: inline SVG cubist wild-cat mask: chrome ears with pink insides, a
  pink leopard-print face, gold cheek plane, oversized white sunglasses with pink gradient
  lenses and gloss streaks, a chrome rhinestone tiara with a pink centre stone, cheetah
  tear-lines, white whiskers, on a gold sun with op-art rings and four twinkling sparkles.
  A static pink-and-gold glow sits behind it (no CSS filters on animated parts).
- `components.Phone(sighting)`: a pink flip phone tilted over the art's lower left, its blue
  LCD showing "1 new message" and the newest sighting as a lowercase text ("omg!! … xoxo"),
  clamped to six lines. Under 820px it sits below the art's bottom right instead.
- With the ticker, fills exactly one screen under the header on desktop (`100svh`, clamped).

### Ticker (`components.Ticker(sighting.Summary)`)
- Full-bleed gold strip, ink mono uppercase, ✦ between items, scrolling. Pauses on hover and
  with its Pause button; static under reduced motion.

### Floor plan (`components.FloorPlan([]Room)`)
- A Mondrian in chrome: 12-column grid, 8px chrome lines. Rooms 1–5 by rank take spans
  7×3, 5×2, 3×3, 4×2, 3×2, plus a 2×3 leopard-print art tile; dense packing fills holes,
  and rooms 6+ follow as 3×2 tiles. Colours cycle gold, lime, black, chrome, charcoal-gold.
- Colours cycle gold, hot pink, black, chrome, charcoal with baby-pink text.
- Each tile links to its room (and opens it): room number, name, big count (sized with
  container query units), "sightings", ↗.

### Rooms (`components.Rooms([]Room)`, built once by `components.NewRooms`)
- One `<details name="rooms">` per species, most seen first; only one open at a time, and
  closed rooms cost the browser nothing to lay out.
- Summary: gold room-number pill, species, scientific name, "N sightings · N countries ·
  date range", pink +/–. Open rooms get a chrome border.
- Inside: a table, newest first: Seen, Where, Observers, Count (a dash when unknown),
  sticky gold headers, faint stripes, pink hover. Under 720px rows stack: date and count,
  then place, then "by" observers.

### Footer
- Chrome rule; name in pink script, "N sightings on show. Data from GBIF.org.", "Best viewed
  at 1024×768".

### Sparkle cursor trail
- In `assets/js/app.js`: small four-point sparkles in pink, baby pink, gold and white follow
  the mouse and fall away (700ms), at most 25 alive. Mouse only, never under reduced motion.

## Avoid
- Photos, image-heavy layouts, masonry, JS libraries, `backdrop-filter`, and CSS filters on
  animated elements.
- White text on gold/pink/chrome; acid lime and orange from earlier designs.

## Accessibility
- 3px hot-pink focus ring (ink inside pink and gold tiles); skip link; 44px targets.
- Decorative ✦ and arrows are hidden from screen readers; the counter has a hidden plain number;
  the cat art has an SVG `<title>`; each room table has a caption.
- `prefers-reduced-motion` stops the ticker, the art's sparkles and the cursor trail, and hides
  the Pause button.
