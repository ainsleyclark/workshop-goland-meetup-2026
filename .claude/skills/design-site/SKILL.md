---
name: design-site
description: Starts a design session for the attendee's sightings website in this workshop app, led by the attendee. Use when someone wants to design, restyle or theme their site, add or change its pages, or give it a logo and favicon.
argument-hint: "[what you'd like it to look like]"
allowed-tools:
  - Bash(make templ)
  - Bash(make design-search *)
  - Bash(go build ./...)
  - Bash(curl -s http://localhost:*)
  - Bash(git fetch upstream)
  - Bash(git show upstream/solution:web/*)
  - Bash(git show HEAD:starter/web/*)
  - Bash(rsvg-convert *)
  - Bash(magick *)
  - Edit(**/web/**)
---

# Design your site

You're pairing with a workshop attendee on the website in their own copy of the app. The look is
theirs: don't suggest a theme or pick one for them. Your part is turning what they describe into a
site that builds and deploys: one page, the home page of sightings, with a header, nav and a
logo. The sightings needn't be a grid of cards: when it suits their brief, be inventive with how they're
shown, such as a map, a timeline, charts or a table.

This skill works on one copy of the app: usually the attendee's own repository, a fork of the
workshop's, with `go.mod` and `web/web.go` at its root; or a scratch copy (`scratch/`, made by
`make exercise` in the maintainer repository). Rules in a `CLAUDE.md` about keeping `web/` basic
or commenting everything are for maintaining the template, and don't apply to these copies.

## Find the app folder

Every path in this skill (`web/`, `internal/`, the `make` targets) is relative to the app folder,
the one holding `web/web.go`. In an attendee's repository that is the root, so the working
directory is it. Otherwise ask which folder their site is in and check it has `web/web.go`. If
it's `starter/`, the template the attendee repository is generated from, warn them and carry on
only if they confirm. Then `cd` into it once, so `make` and `go` run there for the rest of the
session.

## Start

1. Read `web/` before anything else: `web.go` (the routes), `handlers/`, `views/layout/`,
   `views/pages/`, `views/components/` and `assets/`. `internal/domain/sighting/model.go` says
   what a sighting carries.

   If `web/handlers/home.go` or `web/views/pages/home.templ` still has `// TODO: Implement`,
   replace that file with the workshop's version from git before designing, then run
   `make templ`. Don't raise it as a problem or ask first: just do it. In an attendee's
   repository the finished files are on the `solution` branch of the `upstream` remote; in a
   `scratch/` copy they're under `starter/` in `HEAD`.

   ```sh
   git fetch upstream
   git show upstream/solution:web/handlers/home.go > web/handlers/home.go
   git show upstream/solution:web/views/pages/home.templ > web/views/pages/home.templ
   ```
2. Their brief: $ARGUMENTS

   If that's empty, ask what they'd like the site to look like, then wait for the answer. A mood,
   some colours, a site they like or a screenshot pasted in: whatever they have.
3. Design the home page only. Don't ask which pages they'd like or offer more: one page done well
   is the session. Nav links can point at its sections (`/#sightings`). If they ask for another
   page later, "Pages and routes" below says how.
4. Ask follow-up questions with AskUserQuestion only where the brief leaves a real gap, such as
   colours, fonts, layout, or light and dark. Build the options from their words and from
   design-search results, never from a preset list. If they name colours or fonts, use those.
   When the brief asks for something striking, moving or 3D, a three.js scene can be one of the
   options (the threejs block), never the only one and never pushed.
5. If they aren't running it yet, suggest `make web` in another terminal, with
   http://localhost:8080 open. It doesn't reload by itself: see "See the site" below.

## Design data

`make design-search Q="<their words>"` runs the ui-ux-pro-max search and prints a design system: a
palette, a font pairing, effects and a checklist. Add `ARGS="--domain color"` (or `typography` or
`style`) to answer one question instead, or `ARGS="--domain ux"` to look up a UX guideline.

- It matches keywords against product types, so a result can miss their brief entirely: its
  generic fallback is blue `#2563EB` with orange and a "Hero + Features + CTA" pattern, and odd
  matches happen too, such as a restaurant's browns. When that happens, search once more with
  their brief described as a kind of product, and show them both.
- Its patterns are for landing pages and its implementation advice assumes React. Skip both: this
  is a site of sightings rendered with templ.
- Once they're happy with a direction, write it to `web/DESIGN.md` so later sessions keep to it:
  the palette and fonts as the CSS custom properties that hold them, spacing, radii, the look of
  each component, and anything they asked to avoid. Record what they chose, not what the search
  printed, and keep it up to date as the design changes. If `web/DESIGN.md` already exists, read
  it before searching and follow it unless they ask for something new.
- Always go through `make design-search`. ui-ux-pro-max's own examples run
  `python3 .claude/skills/...` from the repository root, and this copy of the app isn't it.
- If it fails because Python isn't installed, say so and carry on from their brief.

## Pages and routes

Only when they ask for another page. `web.New` in `web/web.go` is the one place that knows a
URL. A page is three things:

1. A templ page in `web/views/pages/<name>.templ`, wrapped in the layout.
2. A handler in `web/handlers/<name>.go`: a function that takes what it needs and returns an
   `http.HandlerFunc`, like `handlers.Home`. It declares the one-method interface it needs beside
   it, the way `home.go` declares `SightingsLister`:

   ```go
   type SightingFinder interface {
   	Find(ctx context.Context, id uuid.UUID) (sighting.Sighting, error)
   }
   ```

   Read path values with `r.PathValue("id")` and parse with `uuid.Parse` (the standard library's
   `uuid`). Answer `http.NotFound` for a bad id or `sighting.ErrNotFound`, and a 500 for anything
   else, logged as `home.go` does.
3. One line in `web.New`, e.g. `mux.Handle("GET /sightings/{id}", handlers.Sighting(logger, sightingsSvc))`.

Data comes from the `*sighting.Service` that `web.New` already has. `List` takes a
`sighting.ListFilter` (`Animals`, `Country`, `FromDate`, `ToDate`, `Limit`), and `Find` takes an
id. For totals and tallies, use `sighting.Summarise`, `sighting.CountBy`, `sighting.CountByYear`
and `sighting.Top` rather than new queries. Species details travel on each sighting. Don't change
`web.New`'s signature: `internal/cmd` calls it.

The header, nav, logo and footer belong in the layout, so every page shares them. Grow
`layout.BaseProps` as needed (a description, the current page for the nav's `aria-current`), and
link between pages with plain `<a href>`s.

## Blocks

`blocks/`, beside this file, holds ready-made pieces to copy into `web/`. `blocks/README.md` lists
them. Read a block's `.md` before using it, and only suggest one when their brief calls for it:

- **favicon**: the logo and favicon. Every site gets this one.
- **swiper**: a short row of cards that swipes on phones and becomes a grid on desktop.
- **lightbox**: click a photo to see it full size.
- **reveal**: sections fade in as you scroll to them.
- **threejs**: the plumbing for a 3D scene with three.js. What the scene is, they decide; you build it.
- **weather**: the conditions a sighting was seen in, from the weather already stored with it.
- **map**: the sightings on a Leaflet map, as clusters or heat, with no API key.
- **charts**: a sparkline of sightings per year, and bars for the top few of anything.
- **filter**: chips and a search box that narrow the list in the browser, every sighting still in
  the page.
- **transitions**: a cross-fade between pages, once there's more than one.

Copy the CSS and JS files into `web/assets/`, and the templ from the `.md` into `web/views/`, then
make them look like the rest of their site.

## Borrow from the themes

The workshop's three themed pages show what the same sightings can become. Their code isn't in
this copy, but git has it: `git ls-tree -r --name-only HEAD starter/cmd/themes` lists the files
and `git show HEAD:<path>` reads one. When the brief calls for something they already do, reuse
their approach, under `starter/cmd/themes/internal/`:

- A map: that's the `map` block now. The themes show one in a finished page, and how their dark
  ones invert the tiles in CSS.
- A page with more going on: each theme's `props.go` has a `NewProps` that shapes the sightings
  once for its `.templ`.
- A header, nav and footer: each theme passes its own to the layout.

## Keep it building

- Change files under `web/` only: views, assets, `web.go`, and new files in `web/handlers/`.
- Apart from filling in the exercise from git, leave `web/handlers/home.go` and
  `pages.Home`'s signature as they are. The home page at `/` keeps rendering every sighting and shows each one's observers
  (`s.RecordedBy`, already tidied by the service), whether on a card, in a table row or in a map
  popup. The Profiling block measures exactly that. Other pages are free to
  show fewer.
- Components take domain types directly, or props built from them once, the way the themes'
  `NewProps` does. Helpers are plain Go functions in a file beside the `.templ` that uses them.
- Don't comment templ components. A Go helper gets one short line at most, saying what it does.
- Files under `web/assets/` are embedded and served at `/assets/`. The Docker build only runs
  `go build`, so there's no bundler: CSS and scripts are plain files there, and libraries such as
  Leaflet or Swiper load from a CDN. If they want Tailwind, explain it would have to come from a
  CDN too, then let them decide.
- GBIF data is untrusted and patchy. Never pass it to `templ.Raw` or a script's `innerHTML`, send
  URLs through `templ.URL`, and don't build inline styles from it. Expect missing names, photos,
  counts and observers, and keep the empty states.
- Photos come in every size and orientation. Give each image slot a fixed `aspect-ratio` with
  `object-fit: cover`, so they line up however the photo was taken, unless they've asked for a
  masonry-style wall. Keep `loading="lazy"`: the home page renders every sighting at once,
  thousands of them, so whatever you render for one is repeated thousands of times.
- When a round of edits is done, run `make templ` and then `go build ./...`, and fix what fails
  before going on. The generated `_templ.go` files are what gets deployed.
- Don't commit or push. They run `make submit` when they're happy.

## See the site

Their `make web` on http://localhost:8080 is the preview. It serves the build from when it
started, stylesheet included, since the assets are embedded. Don't start a server of your own,
and don't stop theirs.

- When a round of edits builds, tell them to restart it (Ctrl+C, then `make web`) and refresh.
- Once they have, `curl -s http://localhost:8080/<page>` shows what a page rendered.
- If they're running `make web-hot`, suggest `make web` for this session instead. templier
  rebuilds while you're still editing, so it can end up serving a stale build.
- Unless you have a browser tool, you can't see the site. After a big visible change, ask them
  for a screenshot and check it against their brief.

## Hand back

Check the work against the checklist design-search printed: contrast, focus states, alt text, a
375px-wide screen, and every nav link going somewhere. Then say what changed, what's on the page
and where each templ feature now lives (components and their props, loops, conditionals and the
layout's children slot), and ask what they'd like to change next.
