# 03 — Onion

> Organised by: Rings. A domain in the middle, everything else around it.
> Domain-driven: Getting there. It has our word, but no rules.
> Copies of the sighting: Two — GBIF's, and ours, with a transform between them.

## What it is

We now have a `domain` package.

The layers of 02 bent into rings. `domain` sits in the middle and imports nothing, everything
else points at it. `gbif`, `store` and `handlers` are on the outside and may import `domain`,
`domain` may not import them.

Two things change, and they are the whole point of this step.

First, the app gets its own word. GBIF sends us an **Occurrence**; we store a **Sighting**.
`services/ingest.go` holds the transform between them, so `decimalLatitude` stops one folder in
from the page.

Second, `domain` declares what it needs. `SightingStore` and `SpeciesStore` are written in
`domain`, next to the types they carry, by the code that wants them, and `store` implements them.
The arrow points the other way to 02: storage no longer sits above the model, it answers to it.

`SightingStore` has a name in the literature: it's a **repository**. Two maps behind it today,
SQLite by the end of the day, and `domain` never finds out which.

```
03-onion/
├── domain/
│   ├── sighting.go     Sighting — our shape — and SightingStore
│   └── species.go      Species, and SpeciesStore
├── services/
│   └── ingest.go       GBIF → store, and the two transforms
├── gbif/
│   ├── gbif.go         the API client, and BaseURL
│   └── models.go       Occurrence, Species — GBIF's shapes
├── store/
│   ├── sightings.go    in memory, implements domain.SightingStore
│   └── species.go      in memory, implements domain.SpeciesStore
├── handlers/
│   ├── handler.go      GET /, joins each sighting to its species
│   └── sightings.html  the page, embedded into handler.go
└── main.go             wiring and the server
```

Note that `ingest.go` has not moved. It sat in `services` in 02 and it sits there still — but the
folder now earns its keep, because ingest has something to do beyond passing records along: it is
where GBIF's Occurrence becomes our Sighting.

## When does it suit?

- An app with real rules, where the rules outlive the database you picked for them.
- Anywhere the API you read from isn't the shape you want to keep.
- A codebase that needs to be tested without a network or a database.

## Advantages

- The domain is testable on its own. No HTTP, no JSON, no store, nothing to stand up.
- GBIF can't reach the page. Rename `decimalLatitude` and the change stops at the transform.
- The store is replaceable. Swap the map for SQLite and `domain` doesn't know.
- The interface is written where it is used, not where it is implemented, and so it stays small.
- `var _ domain.SightingStore = (*Sighting)(nil)` — the compiler checks the ring for you.

## Disadvantages

- Look at the folders: *domain, services, store, handlers*. Still what the code **is**, not what
  it's about. We've named the rings instead of the layers.
- `domain` is one package holding two nouns, so everything needs a prefix to tell them apart —
  `SightingStore`, `SpeciesStore`. The package name is doing no work.
- `store.Sighting` and `domain.Sighting` are the same word for a map and a record.
- The domain has no behaviour. Two structs and two interfaces, and every rule still lives
  somewhere else: the display name in `handlers`, the transform in `services`, and the rule that
  a sighting needs coordinates nowhere at all.
- `domain.Sighting` still carries `json` tags. Who is it dressed for, if it imports nothing?
- Adding weather still means opening four folders, they're just better arranged.
- `services` collects anything that isn't a struct. It's the new `main.go`.

Two of those deserve an answer rather than a shrug.

**The `json` tags.** `domain.Sighting` imports nothing and still carries them, which means it is
dressed for a caller it claims never to have heard of. The rule it's breaking is *a model per
layer*, described in [the tour's README](../README.md#a-model-per-layer): tags at the edges, no tags
in the middle. Ours are a leftover from 02, where the struct really was GBIF's wire format, and
nobody deleted them when it moved inwards. That's the honest way this happens.

**`services` as a dumping ground.** A layer earns its name when you can say what it's for in one
sentence and then act surprised when something else turns up in it. *Store the data* can refuse
work. *Talk to GBIF* can refuse work. `services` can refuse nothing, because it has no job
description, only a location. That is exactly why `ingest.go` landed there, and why the next
homeless file will too.

## Questions

- Which package owns the rule that a sighting needs coordinates? Where would you put it now?
- `handlers` looks up a species for every sighting. Whose job is that join?
- The transforms sit in `services`. Do they belong to GBIF, or to the domain?
- If `domain` imports nothing, why does it know the word `json`?
- Is `domain` a domain, or a folder we agreed to point the arrows at?

## What can we improve?

- Give each noun its own package, so `SightingStore` can just be `sighting.Store`.
- Move the rules in with the shape they're about, so the domain does something.
- Let the store live beside the domain it serves, and name the driver (SQLite, for example).
