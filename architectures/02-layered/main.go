package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"

	"02-layered/gbif"
	"02-layered/handlers"
	"02-layered/services"
	"02-layered/store"
)

func main() {
	ctx := context.Background()

	memStore := store.NewMemStore()
	ingester := services.Ingest{
		Client: &gbif.Client{HTTP: http.DefaultClient},
		Store:  memStore,
	}

	err := ingester.Ingest(ctx)
	if err != nil {
		log.Fatalf("ingesting data: %v", err)
	}

	slog.InfoContext(ctx, "Ingested data")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handlers.Homepage(memStore))

	slog.InfoContext(ctx, "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
