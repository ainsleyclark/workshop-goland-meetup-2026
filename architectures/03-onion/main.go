package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"03-onion/gbif"
	"03-onion/handlers"
	"03-onion/services"
	"03-onion/store"
)

func main() {
	ctx := context.Background()

	sightingsStore := store.NewSightingsStore()
	speciesStore := store.NewSpeciesStore()

	ingester := services.Ingest{
		Client:         &gbif.Client{HTTP: http.DefaultClient},
		SightingsStore: sightingsStore,
		SpeciesStore:   speciesStore,
	}

	err := ingester.Ingest(ctx)
	if err != nil {
		log.Fatalf("ingesting data: %v", err)
	}

	slog.InfoContext(ctx, "Ingested data")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Homepage(sightingsStore, speciesStore))

	slog.InfoContext(ctx, "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
