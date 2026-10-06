# Mock API

A local stand-in for the APIs the workshop calls, for when they're down, slow, or rate limiting a
full room. GBIF is served under `/gbif/v1`, and a generated Open-Meteo mock under `/weather/v1`.

```bash
make mock                    # from the starter (or an attendee folder), listens on :8082
go run ./cmd/mockapi         # the same, or add serve --port 8082 to pick the port
```

Then point the app at it from `starter/.env` (or an attendee folder's `.env`):

```bash
GBIF_BASE_URL=http://localhost:8082/gbif/v1
OPEN_METEO_BASE_URL=http://localhost:8082/weather/v1
```

Or skip editing `.env` and pass `--local` to the app's CLI instead, e.g.
`workshop --local sightings ingest -c KE -n 100`.

## GBIF

The mock serves the three endpoints the starter uses: `/occurrence/search`, `/species/{key}` and
`/species/{key}/media`. Results are GBIF's own JSON, stored untouched, so they decode exactly as
they would from the real API.

Search supports `taxonKey` (any rank, so `212` finds every bird), `country`,
`decimalLatitude`/`decimalLongitude`, `geo_distance`, `eventDate`, `limit` and `offset`, following
GBIF's paging rules. One difference from GBIF: `count` is the number of matches in the mock's
data, not GBIF's global total.

### Data

The data is embedded from `internal/gbif/data/gbif.sql.gz`. To refresh it from the live API, run:

```bash
make mock-scrape             # or go run ./cmd/mockapi scrape
```

This searches every animal the starter knows, plus broad groups (mammals, birds, reptiles,
amphibians, fish, insects, plants, fungi) and the theme countries. Then it fetches each species
found and its media. The seed list lives in `internal/gbif/scrape.go` and takes the animals from
`internal/domain/species`. Tune `perSeed` there to keep the file comfortably under GitHub's 100MB
limit.

## Open-Meteo

The mock serves the one endpoint the starter uses: `/archive`, taking `latitude`, `longitude`,
`start_hour`, `end_hour` and `hourly`. Unlike GBIF's fixed set of occurrences, a sighting's
coordinates and event date can be almost anything, so there's no fixed dataset to record ahead of
time. Instead, each reading is generated deterministically from the request: the same place and
hour always return the same weather, so demos and tests stay reproducible, but nothing is scraped
or embedded.

## Layout

| Path | Purpose |
| --- | --- |
| `main.go` | Calls `cmd.Run()`. |
| `internal/common` (the app's) | The HTTP server and request logger are the same ones `workshop web` uses. |
| `internal/cmd` | The `serve` (default) and `scrape` commands, built with urfave/cli like the app's. |
| `internal/gbif` | The GBIF handlers, search filters, embedded data and scraper. |
| `internal/weather` | The Open-Meteo archive handler and the deterministic weather generator. |

Each mocked API lives in its own package, mounted at its own prefix in `internal/cmd/serve.go`.
