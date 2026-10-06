# 02 — By layer

> Organised by: What the code does. One package per layer.
> Domain-driven: No.
> Copies of the sighting: One, GBIF-shaped, used everywhere.

## What it is

The same files as 01, turned into packages. Each file described the kind of code it held, so each
kind of code becomes a folder: a package for the models, one for the GBIF client, one for storage,
one for the handlers.

Nothing else changes. Same types, same store, same handler, same wiring — only the grouping.

You have seen this cut before. *Models, views, controllers* is the same idea, and it's what most
frameworks in most languages ship, because it fits every application ever written. That
universality is the appeal, and it's also the problem: a scheme that fits everything tells you
nothing about this one.

```
02-layered/
├── models/
│   ├── occurrence.go   Occurrence — GBIF's shape
│   └── species.go      Species — GBIF's taxonomy record
├── gbif/
│   └── gbif.go         the API client, and BaseURL
├── store/
│   ├── store.go        the Storage interface
│   └── mem.go          occurrences and species, two maps in memory
├── services/
│   └── ingest.go       GBIF → store (the species lookup)
├── handlers/
│   ├── handler.go      GET /, joining each occurrence to its species
│   └── sightings.html  the page, embedded into handler.go
└── main.go             wiring and the server
```

Note where `ingest.go` ended up. It isn't a model, a client, a store or a handler, it's the thing
that uses all four, so none of the four layers will take it — and the scheme grows a fifth folder,
`services`, named after nothing in particular, to catch it. The one file that describes what this
app actually does is the one file the scheme had no word for.

## When does it suit?

- A CRUD API over one table.
- A team who need a rule for where a new file goes, and need it today.
- A flat package that has outgrown one folder.

## Advantages

- It's easy to decide where things go.
- It makes you start thinking about layers, and about what depends on what.
- GBIF is a package now, with a boundary. It can be tested and replaced on its own.
- Things can finally be private. `MemStore`'s two maps are `store`'s business and nobody else's.
- The dependencies are visible. You can read the import block and see what depends on each other.

## Disadvantages

- Looking at the folders, can you tell what this app does? *models, storage, handlers* describes
  every Go application ever written. Every one of these packages is named for what it
  **contains**. A package name starts earning its keep when it says what it **provides**.
- The naming wobbles as soon as you add a second layer. `models` is plural and `store` is
  singular; `store.Storage` stutters at every call site. And a layer scheme always grows a
  `util` or a `common` in the end, because there's a folder for every kind of code except the
  leftovers, and then everything unrelated collects in it.
- Where do the shared constants go?
- Every layer imports `models`, so `models` can import nothing.
- `gbif` imports `models` too, so the client decodes GBIF's JSON straight into our struct.
  There is one shape here and it belongs to both of them.
- If you want to add weather to a sighting and you open four packages: a model, a client, a 
  store, a handler.
- One feature is spread across four folders, and one folder holds part of every feature.
- The shape is still GBIF's. `models` is the JSON they sent us, renamed.

That third one, *every layer imports `models`, so `models` can import nothing*, is load-bearing,
and it's the classic way this cut fails. Put the types in the middle, let the layers above and
below both reach for them, and the day `models` needs to save something you get this:

```go
// package models
func (o Occurrence) Save() error {
    return store.SaveOccurrence(o) // import cycle: models → store → models
}
```

Not a design smell you can live with for a while. It stops compiling. The only reason this one
builds is that `models` is kept deliberately stupid: no methods, no behaviour, no imports. Which is
a strange thing to demand of the package that is supposed to hold your domain.

## Questions

- If we added weather to each sighting, how many packages would you open?
- Which package owns the rule that a sighting needs coordinates?
- `store` imports `models`. Could `models` ever import `store`?
- Is `models` a domain, or a folder of structs GBIF's JSON happens to fit?
- `gbif` decodes into `models.Occurrence`. Whose type is that — ours, or theirs?
- What else is going to end up in `services`, once it exists?

## What can we improve?

- Group by the thing, not the kind of thing, so one feature lives in one place.
- Let the code that needs something describe what it needs.
- Give the app our own word for a sighting, so GBIF's field names stop reaching the page.

On the `util` question: the app's `internal/common/README.md` is this repository arguing with
itself about where that line falls, and it's worth reading once you've seen where the day ends
up. A dumping ground with a README explaining what may go in it is a different animal from a
dumping ground.
