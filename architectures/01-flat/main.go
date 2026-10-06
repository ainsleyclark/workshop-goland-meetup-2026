package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
)

func main() {
	ctx := context.Background()

	store := NewStore()
	gbif := &GBIF{HTTP: http.DefaultClient}

	err := Ingest(ctx, gbif, store)
	if err != nil {
		log.Fatalf("ingesting data: %v", err)
	}

	slog.InfoContext(ctx, "Ingested data")

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", Homepage(store))

	slog.InfoContext(ctx, "Listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
