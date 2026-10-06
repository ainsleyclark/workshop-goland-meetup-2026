# 01 — Flat

> Organised by: Nothing. One package.
> Domain-driven: No.
> Copies of the sighting: One, GBIF-shaped, used everywhere.

## What it is

A flat, non-hierarchical package all grouped by the kind of code it is, each file describes what 
it does, i.e. the kind of code. A file for types, a file for storage, a file for http handlers 
and so forth.

There's no packages, everything is within scope.

```
01-flat/
├── models.go      Species, Occurrence — GBIF's shapes
├── gbif.go        the API client, and BaseURL
├── ingest.go      GBIF → store (the species lookup)
├── store.go       occurrences and species, two JSON files
├── handlers.go    GET /, and the display-name rule
├── sightings.html the page, embedded into handlers.go
└── main.go        wiring and the server
```

## When does it suit?

- A proof of concept, or quick script.
- A one-file tool.
- A service with one use case.

Most projects should start like this. Start anywhere else and you might abstract too early, and 
make decisions you may regret later down the line.

## Advantages

- It's quick.
- It's simple and not overly complicated.
- There's nothing to navigate, everything is in one place.
- The code that you read, is the code that runs.
- No premature abstraction, no interface written for one implementation.
- No circular dependencies.
- Adding a field is one edit, not 10.

## Disadvantages

- Looking at the files, can you tell what this app does?
- It may get hard to navigate as the project grows.
- Everything can call anything, nothing prevents it.
- There's no way to make anything private.
- There's only one shape, but good design is anticipating structures will be in different forms.
- One shape can't guide its caller. `Occurrence` carries a `Key` that the page needs on the way 
  out and that nobody should supply on the way in, and nothing in the struct says which is which.
- `init()` is the obvious place to seed sample data, and it's a trap: you can't turn it off, so it 
  runs on every `go test` too.

## Questions

- If we added weather to each sighting, where would it go?
- How would you enforce the rule that a sighting needs coordinates?
- How would you test ingest without calling GBIF?
- If GBIF renamed a field, how far does the change reach?

## What can we improve?

- Turn each file into a package, so everything can't call each other.
- Put an interface alongside GBIF, so we can test without the network.
- Split GBIF's shape from ours, so a renamed field can't reach the page.

That first one is worth dwelling on, because it's the whole reason the next four folders exist. Ask
a C programmer what the language lacks and one answer you'll get is a way to wall off one part of a
program from another. You can simulate it, but nothing checks that you did. A Go package is that
wall, and the compiler checks it. Every folder from here on is a boundary somebody has to ask
permission to cross.
