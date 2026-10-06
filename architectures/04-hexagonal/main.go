package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"04-hexagonal/clients/gbif"
	"04-hexagonal/domain/sighting"
	"04-hexagonal/domain/species"
	"04-hexagonal/handlers/sightinghttp"
	"04-hexagonal/stores/sightingmem"
	"04-hexagonal/stores/speciesmem"
)

func main() {
	ctx := context.Background()

	sightingStore := sightingmem.New()
	speciesStore := speciesmem.New()
	client := gbif.New()

	speciesSvc := species.NewService(speciesStore)
	sightingSvc := sighting.NewService(sightingStore, speciesStore, client)

	err := sightingSvc.Ingest(ctx)
	if err != nil {
		log.Fatalf("ingesting data: %v", err)
	}

	slog.InfoContext(ctx, "Ingested data")

	mux := http.NewServeMux()
	sightinghttp.New(sightingSvc, speciesSvc).Routes(mux)

	slog.InfoContext(ctx, "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
