# 05 — Feature

> Organised by: The feature. A folder per noun, and what that noun needs sits inside it.
> Domain-driven: Yes, and the folders say so all the way down.
> Copies of the sighting: Two — GBIF's, and ours.

## What it is

04 named every package after what it is about, then laid them out in a flat list at the root:
`sighting`, `sightingmem`, `sightinghttp`, `speciesmem`. Nothing held them together except that
you named them carefully, and the root grew a package every time a technology arrived.

This one moves each adapter to where it belongs, and that turns out to be two different
directions. `sightingmem` goes *in*, to `domain/sighting/stores/sightingmem`, because a store is
only ever about one domain. A feature is now a folder: open `domain/sighting` and you have the
shape, the rules, the ports it declares and — one directory down, where they cannot be mistaken
for rules — the code that fills them.

`sightinghttp` goes the other way, *out* to `web/handlers`, and the name goes with it. It
was never about sightings; it was about the front page, which happens to show them. Ask what it
imports and the direction gives it away: a store imports its domain and nothing else, while the
handler imports the domain *and* the templates. Filing it under `domain/sighting` would have left
the domain tree importing the view tree — the one arrow this shape is built to prevent.

Three smaller things change with it.

`gbif` moves to `clients/`, because it is the outbound adapter for somebody else's system and it
belongs to no single noun. Persistence is the exception to that folder rather than to the idea: a
store is only ever about one domain, so it lives with it.

`sighting` stops borrowing `species.Store`. It declares `SpeciesReadWriter`, the two methods
ingest actually needs, so `species` can grow a method without widening what a sighting depends
on. `main.go` passes the species *service*, not its store, and the species rules cannot be
stepped around on the way past.

And the join moves off the page. `Sighting` now carries the `Species` it is of, resolved during
ingest, so the handler asks one service one question. In 04 it held two services and did the join
itself; here it declares `SightingsLister`, the single method it needs, and `main.go` is the only
file that knows a `*sighting.Service` satisfies it.

```
05-feature/
├── domain/
│   ├── sighting/
│   │   ├── model.go                        Sighting — our shape, carrying its species
│   │   ├── sighting.go                     the ports, the Service, the errors
│   │   ├── ingest.go                       GBIF → store, and the two transforms
│   │   ├── validate.go                     a sighting needs coordinates and a species
│   │   └── stores/sightingmem/             implements sighting.Store
│   └── species/
│       ├── model.go                        Species, and the display-name rule
│       ├── species.go                      the Store port, the Service
│       └── stores/speciesmem/              implements species.Store
├── clients/
│   └── gbif/                               the API client — Occurrence, Species, BaseURL
├── web/
│   ├── web.go                              New() — the ServeMux, and the only URLs
│   ├── handlers/home.go                    GET /, declares the lister it needs
│   └── views/                              the templates, taking domain types directly
└── main.go                                 wiring and the server
```

`web/` is where the split earns itself. A page is usually about more than one noun — this one is
a sighting and the species it is of, and the next one wants the weather too — so neither the
template nor the handler that renders it belongs to a domain. They sit together instead.

`handlers` is one package with a file per page, and a page is a function: `Home` takes what it
needs and hands back an `http.HandlerFunc`. There is no handler struct and no `Routes` method,
because there is nothing to keep between requests. `web.New` is the only place that knows a URL,
so adding a page is a file in `handlers` and a line in `web.go`. Putting any of this inside a
domain would also bury the rules under generated template code, which matters once the templates
stop being `html/template`.

This is the structure the app uses, so the tour ends where the day's build begins.

## When one domain needs another

Vertical folders promise that a sighting and a species can be worked on separately, and this is the
first folder where that promise is load-bearing. The general rules are in
[the tour's README](../README.md#when-one-domain-needs-another). Here is where each one bites.

**`sighting` imports `species`, and that's allowed.** A sighting holds the species it is of, so the
reference is right there in the record, and `sighting.go` narrows the dependency to the two methods
ingest actually uses:

```go
type SpeciesReadWriter interface {
    Find(key int) (species.Species, error)
    Save(sp species.Species) error
}
```

Turn it around and it fails. `species.Species` holds nothing that points back at a sighting, so
`species` importing `sighting` would be a cycle, and the compiler would say so before you finished
typing. The direction wasn't a taste decision; the data settled it.

**Deleting a species is the case the data forbids.** Nothing in `species` can reach the sightings of
it, and adding a back-reference to make it possible would break the direction above. This is where
an announcement beats an import: `species` says a species was deleted, `sighting` has registered an
interest, and neither folder learns the other's name. Worth noticing that nothing in this folder
does that yet. It's the first thing you'd have to build, and the first place the vertical silos cost
you something.

**"Sightings, sorted by species name, paged" belongs to neither.** `sightingmem` can't answer it
without `species`, and `speciesmem` can't answer it at all. Reading both and sorting in Go is fine
for the 480 rows in the workshop database and not fine for the real thing. The answer isn't to let
`sightingmem` reach into species' data; it's a view, joined in the database, with its own small
read-side folder that queries the view and nothing else. Once storage is SQLite this is nearly free:
a view plus a query file is a read model, and sqlc writes the Go.

**And there is only one database.** The folders tell a story in which `sighting`, `species` and
`weather` each own their own store. Useful story, because it stops anybody writing a join across
them. But the day a write touches two of them you want a transaction, and a transaction wants one
real database. So the app keeps one SQLite file and holds the boundary by convention instead:
`sighting`'s tables, `species`' tables, `weather`'s tables, each prefixed. SQLite has no schemas to
do it for you.

## When does it suit?

- An application with several things in it, each with rules, that will keep gaining more.
- A team that wants to hand somebody a feature, not a tour of four folders.
- Anywhere you want to delete a feature by deleting a directory.

On how finely to slice: start granular. `sighting`, `species` and `weather` are three small domains
and they could defensibly have been one. Merging two folders later is an afternoon; splitting one
that has grown its own tangle of internal references is a project. Err on the side of too many, and
join them when the seams turn out not to be there.

## Advantages

- A feature is a folder. Everything about a sighting is under `domain/sighting`, and nothing
  else is.
- The root stops growing. Adding SQLite adds `stores/sightingsqlite`, not a package at the top.
- The ports are declared by the code that uses them, and named for what it needs —
  `SpeciesReadWriter` is two methods, not somebody else's four.
- Dependencies still point inwards. `sightingmem` imports `sighting`; `sighting` has never heard
  of `sightingmem`, and only `main.go` names both. Nothing under `domain/` imports `web/`.
- Depth follows importance. The rules are at the top of the folder, the technology is a directory
  down, and you can read one without the other.
- The handler is small, because the join is the domain's job and the store is where it happens.
- Adding a page is one file in `web/handlers` and one line in `web.go`. No domain changes.

## Disadvantages

- The paths are long. `domain/sighting/stores/sightingsqlite/store.go` is a lot of typing for a
  file called `store.go`.
- A feature folder can hide coupling. `sighting` imports `species` and that is fine; it is much
  less fine the day `species` wants to import `sighting` back.
- Two folders now hold things called handlers and stores, and only one rule tells them apart:
  a store serves one domain, a page serves whichever it likes. That rule is easy to state and
  easy to forget at 5pm.
- Nesting invites more nesting. There is no natural floor to `stores/`, `handlers/`,
  `internal/`, and each one buys less than the last.
- Still nine packages to render one list. For an app this small, 01 works and fits on a screen.

## Questions

- Where does weather go? It is a noun, but no page is ever only about it.
- `web/handlers` renders sightings. Where does a page joining sightings, species and weather go?
- If a store needed two domains to answer one query, whose folder would it sit in?
- What is left in `clients/` that could not have been a `stores/` folder, and why?
- `Sighting` carries a whole `Species`. What breaks when species get big?
- Which folders would you open to add a "seen in the last week" filter? Which in 01?
- `handlers` declares `SightingsLister`; `sighting` declares `Store`. Both are one consumer
  naming what it needs. Why does only one of them live inside `domain/`?
- `Home` is a function returning a handler, not a struct with methods. What would have to be true
  of a page before the struct earned its keep?

## Two notes before we build

**On mocks.** The usual advice is to keep mocks in a shared subpackage, and it's reasonable advice.
This repository does something different on purpose: the tests run against a real SQLite database
via `dbtest`, with one parameterised setup per package, and there is barely a mock in it. The
reason is that most of what a mock would prove here is that we typed the interface out twice. What
we actually want to know is whether the SQL is right. Both approaches are defensible; ours is
a choice, not an oversight, and it's worth knowing that it's a choice.

**On frameworks.** There isn't one here, and the day uses the standard library plus a handful of
small libraries. A framework hands you a structure for free and takes the decision away in the same
motion. That's a fine trade once you know enough to judge what you're getting, and a poor one while
the structure is the thing you're trying to learn. That's the only reason we build it by
hand: after today you can read a framework's layout and recognise which of these five shapes it
picked, and why.

## What can we improve?

Nothing, in this shape — the next moves are all real ones, and they are the workshop:

- Swap the map for SQLite. `stores/sightingsqlite` implements the same port, written against
  queries sqlc generated, and nothing in `domain/` changes.
- Swap `html/template` for templ, so the page is type-checked against the domain types it draws.
- Add `weather` as a third domain, and watch where the page that needs all three ends up.
- Put it all under `internal/`, so the shape is enforced by the compiler and not by good manners.
