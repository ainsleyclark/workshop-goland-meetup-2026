# TODO

This is everything you'll build today, in order. Each task says which file to open, what to do,
and how to check it worked. The spot in the code is marked `// TODO`. Tick tasks off as you go.
`make todo` lists the markers you still have.

If something breaks and you can't see why, reset the database and go again:
`make run` → `db` → `reset`. Nothing in here is precious.

If GBIF or Open-Meteo stop answering or there's no internet, use `make run-local`.

Stretch goals are optional. Skip them if you're behind. Come back to them if you're ahead.

---

## Block 1

### Architectures

Nothing to write, but we'll read `architectures/` together, in order. Keep these three 
questions in mind for each folder:

- Can you tell what it does?
- Where does a new feature go?
- Where do the rules live?

### The weather client

Nothing to write, but it's worth having a look before you build on it.

- Read `internal/clients/openmeteo`. Note what `Reading` takes in and what it returns.
- Run `make run` → `weather` to see a reading for a place and time.
- Read `internal/domain/weather/model.go`. This is what you're about to store.

### Querying with sqlc

`internal/infra/db/sqlc` is generated from `migrations/` and `queries/`. Your copy knows nothing
about weather yet.

**Migration**

- [ ] Create `migrations/0002_weather.sql`.
- [ ] In it, create the `weather` table. `sightings.weather_id` already points at it.
- [ ] Copy the goose markers and column style from `0001_init.sql`.
- [ ] One column per value in `model.go`. No JSON.
- [ ] Temperature is enough to start. Add the other columns once everything works end to end.
- [ ] Check: `make run` → `db` → `migrate` → `status`, then `up`.
- [ ] If it doesn't apply: fix the SQL, `reset`, try again.

**Query**

- [ ] Create `queries/weather.sql`.
- [ ] Add two queries: `WeatherFind` and `WeatherCreate`.
- [ ] Copy the annotations and the `RETURNING` from `queries/species.sql`.
- [ ] Run `make sqlc`. It regenerates, then builds, so anything broken shows up now.
- [ ] Check: `internal/infra/db/sqlc/weather.sql.go` exists, and `models.go` has a `Weather`.

**Store**

- [ ] Open `internal/domain/weather/stores/weathersqlite/store.go`.
- [ ] Implement `Find`. Call the generated query. 
- [ ] Implement `Create`. Call the generated query.
- [ ] Both return a `weather.Weather`, so map the generated row onto it.
- [ ] Use `speciessqlite/store.go` as your reference. It does the same for species.

**Join**

When a sighting is read back, it should come with its weather.

- [ ] Open `queries/sightings.sql`.
- [ ] In `SightingFind`, `SightingFindByGbifKey` and `SightingList`, add `sqlc.embed(weather)`
      to the `SELECT`.
- [ ] In the same three queries, add a `JOIN` from `weather.id` to `sightings.weather_id`.
- [ ] Run `make sqlc`. The sighting rows now have a `Weather` field.
- [ ] Open `internal/domain/sighting/stores/sightingsqlite/transform.go`.
- [ ] Map `row.Weather` onto `sighting.Weather`. The marker shows where.
- [ ] Open `store.go` in the same folder. In `Create`, pass `transform` a `db.Weather` with just
      its `ID` set, the same way it already passes `db.Species`. An `INSERT` can't join.
- [ ] `sightings` → `list` still shows no weather. That's right: nothing records it yet.

**Stretch**

- [ ] Write `weathersqlite/store_test.go`. Test against a real database, like
      `speciessqlite/store_test.go`. `internal/infra/db/dbtest` gives you a test client.

### The weather service

- [ ] Open `internal/domain/weather/weather.go`.
- [ ] Implement `Find`. It passes straight through to the store.
- [ ] Implement `Create`. Make a new ID with `uuid.New()`, then call the store with it.

**Stretch**

- [ ] Write `validate.go`: a `Validate` method on `CreateParams`.
- [ ] Return every broken rule, each wrapping `ErrInvalid`. Think about ranges (humidity, cloud
      cover, wind direction) and values that can't be negative.
- [ ] Call `Validate` from `Create` before anything reaches the store.
- [ ] Test the service in `weather_test.go` and the rules in `validate_test.go`. The species
      package has examples of both.

---

## Block 2

### Weather ingestion

`Ingest` in `internal/domain/sighting/ingest.go` already calls `weatherPerSighting` for each
occurrence and saves the ID it returns on the sighting. The function itself is yours.

- [ ] Open `internal/domain/sighting/ingest_weather.go`.
- [ ] Call the Open-Meteo client: `s.openmeteo.Reading`. Give it the occurrence's coordinates and
      event date.
- [ ] Call the weather service: `s.weather.Create`. Build its `weather.CreateParams` from the
      reading.
- [ ] Return the ID of the weather you created.
- [ ] Check: `make run` → `db` → `reset`.
- [ ] Then `sightings` → `ingest`, then `sightings` → `list`. Every row should now have weather.

### templ

**The page**

- [ ] Open `web/views/pages/home.templ`.
- [ ] Wrap the page in `layout.Base` and give it a title.
  - [ ] If there are no sightings, say so.
  - [ ] If there are, show the total and render `components.SightingCard` for each one.
- [ ] Run `make templ` to regenerate `home_templ.go`. (`make web-hot` does this on every save.)

**The handler**

- [ ] Open `web/handlers/home.go`.
- [ ] List sightings with `lister.List`.
- [ ] If that fails, show an error page. `error.go` next door has `renderError` for this.
- [ ] Set the `Content-Type` header to `text/html`.
- [ ] Render `pages.Home` with the sightings.
- [ ] Check: `make web`, then open http://localhost:8080.

---

## Block 3

### Theming

The site is yours from here. Claude Code has a skill, `/design-site`, that knows how this app is
put together and how to add a page without breaking the build.

- [ ] Give Claude your API key.
  - GoLand: Air's account settings → more providers → Anthropic.
  - Everyone else: `make key`.
- [ ] Run `claude` in this folder, then type `/design-site`.
- [ ] Tell it what you want: a colour, a mood, a place, an animal. Paste screenshots of sites you
      like.
- [ ] After run `make web` to see your changes.

**Ideas**

- A page per sighting (`/sightings/{id}`).
- A page per animal or per country.
- The weather on each card.
- A statistics page.
- The sightings on a map.
- A photo gallery.
- A custom 404 page.

### Profiling

The home page has three performance problems built in - one in the handler, one in the service,
one in the store. Find them before they're shown to you.

- [ ] Open two terminals.
- [ ] First terminal: `make web`.
- [ ] Second terminal: `make load-pprof`. It loads the site for 20 seconds and saves a CPU
      profile and an allocation profile.
- [ ] Read the `load-pprof` target in the Makefile. Note the `:6060` endpoint it hits, and why
      you'd never expose it in production.
- [ ] Open the profiles with the commands it prints. Start with allocations, then CPU.
- [ ] Ask: where is the time going? What is SQLite being asked to do on every request?
- [ ] GoLand: run the project's run configuration with the profiler. Start a CPU recording, run
      `make load`, stop the recording, read the flame graph.
- [ ] Fix one problem at a time. Run `make load` after each fix to see the difference.

---

## Stuck?

Try it first. But if you're behind and the next block needs this one finished, the finished
version of every file is on the `solution` branch of the workshop repository, at the same path as
yours. `make setup` added it as `upstream`:

```sh
git fetch upstream
git show upstream/solution:internal/domain/weather/weather.go
```

Swap the path for the file you're on. The paths match yours one for one.
