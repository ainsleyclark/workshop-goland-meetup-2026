# Full-Stack Go: Domain-Driven Applications with sqlc and templ

**Level:** Intermediate

## Description

Build a complete full-stack application in Go, from data ingestion through persistence to a
server-rendered UI, purely in Go. This workshop teaches you how to design maintainable
full-stack applications that age well.

You'll start from a working application built on a real public API, walk through five different
architectures for the same feature to see what domain-driven design actually buys you. Then you'll
extend the application with a new domain, a type-safe data layer with sqlc and SQLite, and a 
front-end written with templ.

This workshop is intended to be hands-on and interactive, by making the architecture decisions 
yourself, with guidance at each step. You'll leave with a working application written by hand 
and a clear model of how to structure Go code around a domain rather than a framework.

## What you'll learn

- How to structure a Go codebase around its domain so it ages well.
- How to extend an existing application with a brand-new domain, end to end.
- How to generate a type-safe data layer with sqlc.
- How to serve a user interface directly from Go with templ.

## Syllabus

### 1. Architecture & Domain-Driven Design

What makes a Go project maintainable? We'll go through five different architectures with the 
same features and work out where it's suited best.

- Package layout & structure.
- Domain-driven: What belongs where, and why.
- What pushes you from one architecture to the next.
- A brief look at the finished application.

### 2. Extending the Domain

You add a brand-new domain to an existing application, sourced from a second, independent API.

- Where does a new concept live?
- Modelling entities and operations.
- Translating API responses into domain types.

### 3. A Type-Safe Data Layer with sqlc

Once you're fetching new data, it's time to store it, with sqlc and SQLite.

- Schema design, driven by the domain.
- Writing & generating queries.
- Mapping generated types to domain types.
- Query tests against a real database.

### 4. Ingestion

Wiring the new domain into the application's existing ingestion pipeline.

- Enriching every record as it's ingested.
- Persisting and associating the result.

### 5. Building the Front-End with templ

Exposing the application through handlers that call templ components and layouts.

- templ fundamentals: components, props, layout slots.
- Handlers, routing & error mapping.
- A homepage, themed your way.

### 6. Profiling

Where the application spends its time, and how you can make it faster.

- Profiling the web handler with pprof.
- Load testing, and seeing where it falls over.
- Reducing memory usage.

## Prerequisites

- Several months writing Go.
- Familiarity with SQL and relational databases.
- Go 1.27.1 or newer on the device you bring, with `~/go/bin` on your PATH.
- macOS or Linux, or Windows with WSL2. You'll need `make` (on macOS: `xcode-select --install`).
- A GitHub account. On the day you'll fork the workshop repository and work in your fork.
- Basic knowledge of git.
- Basic HTML and CSS. No JavaScript framework experience needed.

## Recommended Preparation

- Join the dedicated Slack channel.
- Go 1.27.1 (or latest version, installed).
- Claude Code installed and up to date.
- Python 3 installed (used by the design tooling).
- GoLand preferred — we'll provide a 3-month free licence!
