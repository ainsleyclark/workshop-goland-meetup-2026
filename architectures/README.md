# Architectures

The opening tour. Five folders, each implementing the same small feature: import occurrences from
GBIF and store them, then list the store on a page.

Nobody types this code in. We read it, in order, and ask the same three questions of each one. The
day's build starts afterwards, from the app this folder sits inside, which is the shape the last
folder arrives at.

This page holds the ideas. Each folder holds one application of them, so if a README below mentions
a value object or a bounded context and you want the definition, it's here.

## How to read it

It's a ladder, not a catalogue. Every step exists because the one below it hurts, and the hurt is
the lesson, so the early folders are not strawmen and the late ones are not the answer. `01-flat`
is the right choice for most things you will write. The ladder only pays for itself when the pain
it removes is pain you actually have.

Read each one with three questions in hand:

1. Looking at the folders, can you tell what this application is about?
2. If we added weather to each sighting, how many folders would you open?
3. Where does the rule that a sighting needs coordinates live?

| Folder           | Organised by                             | Domain-driven |
|------------------|------------------------------------------|---------------|
| `01-flat`        | nothing. One package, files by kind      | No            |
| `02-layered`     | what the code does                       | No            |
| `03-onion`       | rings, with a domain in the middle       | Getting there |
| `04-hexagonal`   | what the code is about. One per noun     | Yes           |
| `05-feature`     | the feature, all the way down            | Yes           |

Every one compiles, so you can navigate it in the IDE while we talk. Bodies are sometimes empty;
the structure is the point.

## Domain-driven design is not about where your files go

This is a tour of folders, so it is worth saying plainly before we start: DDD is a way of thinking
about the problem, not a directory layout. You can apply it to `01-flat` without moving a single
file, and you can produce `05-feature`'s folder tree while thinking about nothing but folders.

The structure is a *consequence*. What follows is five attempts at letting the structure show the
thinking.

## The words

**Ubiquitous language.** One set of words, shared by everyone who talks about the system:
developers, testers, the people who asked for it. Pick a word and use it everywhere, in the code
and in the conversation.

Watch it happen across this tour. GBIF calls a record an **Occurrence**. We call it a **Sighting**.
`01-flat` and `02-layered` use GBIF's word because they have no word of their own; from `03-onion`
the application has its own vocabulary, and GBIF's leaks no further than one transform. That shift
is deliberate, so don't read it as inconsistency.

The cost of skipping this is boring and constant: the same idea called `storage` in one package,
`database` in the next and `repository` in the third, and an API that ends up with endpoints
spelled both `check` and `cheque` because nobody agreed which word the company used. (This
repository writes British English throughout for the same reason. One spelling, everywhere.)

**Bounded context.** The boundary inside which a word means one thing. The same noun means
different things to different readers, and trying to make one model serve all of them is how models
get bloated. A sighting in this application is a good example:

- To **ingest**, it's a GBIF key, a species to resolve, and a duplicate to detect.
- To **the page**, it's a name, a place and a date, and it has never heard of a GBIF key.
- To **a theme**, it's a row that either belongs on this page or doesn't. `kenya`, `humpback` and
  `antarctica` are the same data through three lenses, which is about as literal a demonstration of
  the idea as you will get.

**The building blocks.** The formal vocabulary, in the order you tend to need it:

- **Entity.** A thing with an identity of its own, so two of them with identical fields are still
  two different things. A sighting is one, and its key is what says so.
- **Value object.** A thing you care about by its value, not by which one it is. A coordinate pair,
  a country code, a barcode. It has no identity and it doesn't need one.
- **Aggregate.** A cluster of entities treated as one unit, reached through one of them, the
  **root**. Save it together, load it together, and let the root enforce the rules.
- **Repository.** A façade over persistence, described in the domain's terms rather than the
  database's, and deliberately vague about how many databases there are or what kind.
- **Domain service.** A stateless operation that no single entity can sensibly own. Ingest isn't
  something a sighting does to itself.
- **Invariants.** The rules true of every instance, kept in one place.
- **Domain event.** Something happened that other parts of the system may care about.
- **Factory.** A way of building something complicated enough that a struct literal won't do.

You won't find all of these in every application, and there is no prize for forcing them in. The
vocabulary is for describing what you have, not a checklist to satisfy.

## A model per layer

Every layer owns its own shape. The tags belong at the edges, `json` on the way in and out of an
API and `db` at the store, and the shape in the middle has no tags at all, because it is nobody's
wire format.

A `json` tag on a domain type means the domain is dressed for a caller it claims never to have
heard of. It's also how an API's field names quietly become the database's field names and then the
domain's field names: three layers with one shape, and every rename touching all three. The cost of
the rule is a transform per boundary. That transform is also the only place a rename can hurt you.

## Parse, don't validate

Validating after the fact is the weaker of the two ways to enforce a rule. It runs once the value
already exists, so the bad value is a thing you can build anywhere, and it stays bad until somebody
remembers to ask.

The stronger way is to make it unbuildable: give it a type, and one function that is the only way
in.

```go
type Coordinates struct{ Latitude, Longitude float64 }

func ParseCoordinates(lat, lng float64) (Coordinates, error) // the only door in
```

"Latitude is between -90 and 90" stops being a rule anyone can forget and becomes a rule you had to
satisfy to hold a `Coordinates` at all. Everything downstream can stop asking.

One trap. The obvious spelling doesn't do what it looks like:

```go
type Code string             // not typed enough: `var c Code = "nonsense"` compiles
type Code struct{ v string } // nothing to convert from, so the compiler has your back
```

With a string underlying type the compiler converts untyped constants for you, so your parse
function is a door standing beside an open window. The struct form closes it, at the price of
turning constants into variables and needing `MarshalText` and `Scan` before JSON, sqlc or templ
will touch it. Most Go codebases take the first form and the open window knowingly. The app does
too, and you'll see where.

## When one domain needs another

Vertical folders are a promise that two domains can be worked on separately, and the promise gets
tested the first time they need each other. Three cases, three different answers.

**Reading across, when the data allows it.** If this domain holds that domain's identifier, it has
earned the right to ask about it. Follow the data and the question of who may import whom mostly
answers itself; turn it around, where the reference doesn't exist, and you get a cycle the compiler
will refuse. The foreign keys were telling you the dependency direction all along.

**Writing across, when the data doesn't.** Deleting one thing often has to affect another that
holds no reference back. The cheap answer that keeps the boundary is an event: one domain announces
what happened, the other registers an interest, and neither imports the other. Be clear about what
this needs and what it doesn't. A map of functions and a synchronous call through it is enough. No
goroutines, no message broker, no eventual consistency to debug at 2am. You get the decoupling and
you keep the stack trace.

**Querying across, which is neither.** "Everything, sorted by a field in the other domain, paged"
can't be answered by either domain alone, and it can't be answered by reading both and joining in
Go, because sorting a million rows in memory to show ten of them isn't a plan. Don't let one store
reach into another's tables. Let the database do the join *underneath* the architecture, in a view,
and give the view its own small read-side domain that queries it and nothing else. Reports are
where domain boundaries get quietly violated in every codebase, and this is the valve that keeps
them honest.

All three lean on something the folders pretend isn't true: there is one database. The folders tell
a story in which each domain owns its own store, and it's a useful story, because it stops anybody
writing a join across two domains' tables. But a single real database is what makes transactions
possible, and you need those the moment a write touches two domains. So you keep one database and
hold the boundary by convention instead: a table prefix per domain, or a schema each where the
engine has them. SQLite doesn't, so: prefixes.

## Two rules of thumb

**The structure should reflect how the software actually works.** The test to apply when a layout
feels wrong: does the shape on disk tell the truth about the shape in the running program? Your code
is the only documentation that is never out of date, so if the two disagree, the folders are lying
to the next person who reads them.

**The rule of five.** The one you can use on your own repository on Monday: count the folders at
each level. The research on this puts what a person can hold in their head at around five things.
More than five at any one level and nobody can keep a mental model of where anything is, which
matters most at the root, where every newcomer starts.

## Be like water

You will not get the structure right first time. Nobody does, and the language gives you no hints;
you start with a blank page.

Expect to move things. Your domain will grow, the words will change, and the structure should
follow. These five folders are five recorded attempts at the same feature. That is the shape of the
work, not a sign anybody failed.

## Further reading

The vocabulary comes from Eric Evans' *Domain-Driven Design* and Vaughn Vernon's *Implementing
Domain-Driven Design*. Both are thick books. Neither is required reading for today.
