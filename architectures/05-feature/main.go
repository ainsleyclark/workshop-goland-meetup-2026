package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"05-feature/clients/gbif"
	"05-feature/domain/sighting"
	"05-feature/domain/sighting/stores/sightingmem"
	"05-feature/domain/species"
	"05-feature/domain/species/stores/speciesmem"
	"05-feature/web"
)

func main() {
	ctx := context.Background()

	sightingStore := sightingmem.New()
	speciesStore := speciesmem.New()
	client := gbif.New()

	speciesSvc := species.NewService(speciesStore)
	sightingSvc := sighting.NewService(sightingStore, speciesSvc, client)

	err := sightingSvc.Ingest(ctx)
	if err != nil {
		log.Fatalf("ingesting data: %v", err)
	}

	slog.InfoContext(ctx, "Ingested data")

	router := web.New(sightingSvc)

	slog.InfoContext(ctx, "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
