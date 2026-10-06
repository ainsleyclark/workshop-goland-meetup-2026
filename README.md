# Full-Stack Go: Domain-Driven Applications with sqlc and templ

**GoLand Collaboration Serbia**

The workshop app: a site for **animal sightings**. Occurrence records come from
[GBIF](https://www.gbif.org), are stored in SQLite through [sqlc](https://sqlc.dev), rendered
server-side with [templ](https://templ.guide), and enriched with
[Open-Meteo](https://open-meteo.com) weather.

This repository is yours for the day. You fork it, build in your fork, and push there.

## Before you start

You need Go (the version in `go.mod`), Git and a GitHub account. GoLand is the editor used on
the day; anything that runs `make` and `go` works.

1. **Fork this repository** on GitHub: the *Fork* button at the top of
   https://github.com/ainsleyclark/workshop-goland-meetup-2026. If you have the GitHub CLI,
   `gh repo fork ainsleyclark/workshop-goland-meetup-2026 --clone` forks and clones in one go.
2. **Clone your fork** and open it:

   ```sh
   git clone git@github.com:<your-github-user>/workshop-goland-meetup-2026.git
   cd workshop-goland-meetup-2026
   ```

   Open this folder in GoLand. The `go.mod` is at the root, so there is nothing else to point
   it at.
3. **Run setup once.** It installs the Go tools the day needs and points `upstream` at the
   workshop repository so you can pull fixes:

   ```sh
   make setup
   ```

## During the day

```sh
make run        # the interactive CLI: migrate, ingest sightings, list, weather, stats, web
make web        # serve the site on :8080
make web-hot    # the same, with templ hot reload
make templ      # regenerate templ after editing a .templ file
make sqlc       # regenerate sqlc after editing a query or migration
make test       # the tests
make submit     # commit everything and push it to your fork
make update     # pull the instructor's latest changes
make help       # everything else
```

Your database starts empty: `make run`, then `db` → `migrate` → `up`, then `sightings` →
`ingest`. If GBIF or Open-Meteo buckle under the room, `make run-local` uses the mock APIs
(`make mock` serves them).

`make submit` pushes to **your fork**. If the GitHub CLI is installed it also opens a draft pull
request on the workshop repository, so the instructor can see what you've built; it is never
merged.

## Where things are

```
main.go                     thin; hands off to internal/cmd
migrations/                 goose migrations
queries/                    sqlc queries, a file per table
internal/
├── cmd/                    CLI commands and the interactive menu
├── domain/                 the rules: sighting, species, weather — and their stores
├── clients/                outbound adapters: gbif, openmeteo
├── infra/                  config, sqlite, sqlc output
└── common/                 small generic helpers
web/                        your view layer: web.go owns the routes, handlers/, views/, assets/
architectures/              the opening tour, a ladder of five layouts; read, not typed in
.claude/, .junie/           Claude Code and Junie: settings and the /design-site skill for the Theming block
```

Each of `internal/{domain,clients,infra,common}` has a README saying where its line falls.
`SYLLABUS.md` is the course outline you were sent beforehand.
