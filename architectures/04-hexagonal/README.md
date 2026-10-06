# 04 — Hexagonal

> Organised by: What the code is about. One package per noun, and one per technology that serves it.
> Domain-driven: Yes.
> Copies of the sighting: Two — GBIF's, and ours.

## What it is

03 put a domain at the centre and pointed everything at it. This one changes two things, and
nothing else. The feature is the same, the storage is still a map, the page is still one page.

First, `domain` splits into a package per noun. 03 had one `domain` package holding every type and
every interface, which meant one shared bag of ports: `SightingStore` and `SpeciesStore` sat side
by side whether or not anything used both. Now `sighting` and `species` each declare the port they
need, and the name stops stuttering — `sighting.Store`, not `domain.SightingStore`.

Second, every package outside `domain` is named for the noun it serves *and* the technology it
uses. 03 kept one `store` package with `sightings.go` and `species.go` inside it; here those become
`sightingmem` and `speciesmem`, and `handlers` becomes `sightinghttp`. The folder above is only a
drawer — the files have to be filed somewhere — but the package name is the part you type at the
call site, and it now says what the thing actually is. That is not decoration. `stores/` holds two
stores, and they cannot both be `package mem`.

```
04-hexagonal/
├── domain/                  our shapes and rules — no net/http, no GBIF JSON
│   ├── sighting/
│   │   ├── model.go         Sighting — our shape
│   │   ├── sighting.go      the Store port, the Service, the errors
│   │   ├── ingest.go        GBIF → store (the species lookup)
│   │   └── validate.go      a sighting needs coordinates
│   └── species/
│       ├── model.go         Species — our shape, and the display-name rule
│       └── species.go       the Store port, the Service
├── clients/                 what we call out to — the driven side
│   └── gbif/
│       ├── client.go        Client, and BaseURL
│       ├── occurrences.go   Occurrence, and the call that returns it
│       └── species.go       Species, and the call that returns it
├── stores/                  what fills the Store ports
│   ├── sightingmem/         implements sighting.Store
│   └── speciesmem/          implements species.Store
├── handlers/                the driving side
│   └── sightinghttp/        GET /
├── web/                     the page, and the shapes it draws
└── main.go                  wiring and the server
```

Note what happened to `ingest.go`. In 02 and 03 it landed in `services`, a folder invented to hold
whatever wasn't a model, a client, a store or a handler. Here it is obviously a sighting thing, so
it sits with the sightings. The scheme finally has a *name* for the file that says what the app
does, rather than a drawer for everything it couldn't classify.

Note also the two copies of a sighting. `gbif.Occurrence` is the JSON they send us, with
`decimalLatitude` and a date that isn't a date. `sighting.Sighting` is ours. They meet in exactly
one function, `transformOccurrence`, and that is the whole point — GBIF can rename a field and
the blast radius is one line in one file.

And note the three drawer names. `clients`, `stores` and `handlers` are the words the workshop
application uses too, so nothing here gets renamed on the way up: what is left for the next step
is a *move*, not a rewrite.

## The vocabulary, against this code

This is the first folder where the DDD words have something to point at. The definitions live in
[the tour's README](../README.md#the-words); here is where each one landed:

| Word           | Here                        | Why                                                        |
|----------------|-----------------------------|------------------------------------------------------------|
| Entity         | `sighting.Sighting`         | Two records, identical fields, different keys. Two things. |
| Value object   | `species.Species`, nested   | Inside a sighting you want what was seen, not which row.   |
| Repository     | `sighting.Store`            | Persistence in the domain's words, implemented elsewhere.  |
| Domain service | `sighting.Service`          | Ingest isn't something a `Sighting` does to itself.        |
| Invariants     | `validate.go`               | True of every sighting, in one place.                      |
| Domain event   | `ErrNotFound` is the nearest | Something happened, and it was worth naming.               |

Note what isn't here. No factories, no aggregate worth the name, because nothing in this application
is complicated enough to need them yet. `species.Species` being an entity in its own right *and* a
value object inside a sighting is not a contradiction either; it's the same data in two contexts,
which is the whole point of drawing contexts.

The aggregate is worth holding onto, though, because it explains the next folder. Open
`sightinghttp`: it holds two services and does the join between them itself, because nothing here
says a sighting and its species travel together. 05 says it, and the handler gets smaller.

## When does it suit?

- An application with more than one thing in it, where the things have rules.
- A service that will outlive the database or the API it was first written against.
- Anywhere you want to test the rules without a network or a schema.

There's a reason that list doesn't mention how big the codebase is. The complexity you're adding
here doesn't buy anything technical, because 01 renders the same page. It buys the ability to put
people in different folders without them treading on each other, and to hand the thing over to
somebody who wasn't in the room. That's the trade: extra indirection now against somebody else's
comprehension later, and it's a good trade exactly when there is a somebody else.

Which cuts the other way too. On a project that is only ever you, this is overhead. And don't buy
it because you expect a team next quarter. Build what you need today, and the day the second person
arrives, move the folders then.

## Advantages

- Looking at the folders, you can tell what this app is about. *sighting, species* is not every Go
  application ever written.
- The domain declares the port; the adapter implements it. `var _ sighting.Store = (*Store)(nil)`
  is the hexagon, and the compiler checks it.
- Dependencies point inwards. `sightingmem` imports `sighting`; `sighting` has never heard of
  `sightingmem`. Only `main.go` names both.
- Swapping storage is one line in `main.go`. Add `sightingsqlite`, delete `sightingmem`, and
  nothing in `domain` changes.
- The rules have a home. `validate.go` is the question 01 and 02 couldn't answer.
- The domain owns its error vocabulary, so `ingest` can tell "no such species" from "the store is
  broken" — which 03 could not, and quietly refetched on every error.
- Testing ingest needs no network: implement two small interfaces and pass them in.

## Disadvantages

- It's more folders and more files for the same behaviour. Seven packages to render one list.
- Two shapes for one idea means a transform, and a transform means somewhere to forget a field.
- The drawers are still sorted by kind of code. `stores` and `handlers` are 02's words, filed one
  level down: ask where *the sighting feature* lives and the answer is still four folders.
- `sighting` imports `species` for its store. That port was written for `species`' convenience,
  and `sighting` now depends on all of it to use one method.
- The indirection is real. Reading `Save` means finding which adapter was wired in `main.go`.
- Ports invite ceremony. There is one implementation of `sighting.Store` here, and the interface
  earns its place only because a second one is coming.

### Where `validate.go` stops short

`validate.go` answers the question 01 and 02 couldn't, and it's still the weaker of the two answers
available, because it runs *after* the sighting exists. The stronger one is *parse, don't validate*,
described in [the tour's README](../README.md#parse-dont-validate).

Two fields here are asking for it. `model.go` carries
`EventDate string // "2016-01-11", and not a time.Time`, a field advertising that it holds three
formats and trusting you to cope, and `Latitude`/`Longitude` are bare floats whose range is checked
somewhere else entirely. A `Coordinates` type with a `ParseCoordinates` function would move that
rule to the one place a coordinate can come into existence. The app does exactly this, and the
GBIF client there already parses those three date formats at the boundary, so nothing inland has to.

### One `main.go`, for now

Every folder here has one binary, so one pressure never appears: the day you want the server *and* a
seeder *and* an interactive CLI from the same code. Wire all three into one `main` and you get flag
soup, which is where a `cmd/` folder with a subdirectory per binary comes from, each one short, each
one just assembling services that already exist. The app's `internal/cmd` is that folder.

## Questions

- If we added weather to each sighting, which packages would you open?
- `sightinghttp` needs a sighting *and* a species. Which domain should own that page?
- `sighting` imports `species`. Could `species` ever import `sighting`?
- What does `sighting` actually need from `species` — all four methods, or one?
- How many packages know the string `"/"`? How many know `decimalLatitude`?
- `stores/sighting` would be a shorter name than `stores/sightingmem`. What would `main.go` have
  to do to import it alongside `domain/sighting`?
- `stores` and `handlers` are the folder names 02 was criticised for. Is this the same mistake?
- Why is `web/` outside `domain/`, when `sightinghttp/` could have held the template?
- `sightinghttp` imports `web`. Which way is that arrow pointing, and is it the way you wanted?

## What can we improve?

- Let `sighting` declare the *species* port it needs, too, instead of borrowing one written for
  somebody else. It needs a read and a write, not a `Store`.
- Move each drawer inside the domain it serves, so a feature is one folder rather than an entry in
  each of three. The names already fit: `stores/sightingmem` becomes `sighting/stores/sightingmem`.
- Do the same for `handlers/` and it doesn't fit. A store is only ever about one domain; a page is
  usually about several, and `sightinghttp` proves it by holding two services. Send it the other
  way, to `web/handlers/`, where it sits with the templates it renders.
