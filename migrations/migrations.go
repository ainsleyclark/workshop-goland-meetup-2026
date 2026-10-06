package migrations

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/pressly/goose/v3"
)

//go:embed *.sql
var files embed.FS

// New creates a migration provider backed by the embedded migrations
// within this folder.
func New(conn *sql.DB) (*goose.Provider, error) {
	provider, err := goose.NewProvider(goose.DialectSQLite3, conn, files)
	if err != nil {
		return nil, fmt.Errorf("migrations: create provider: %w", err)
	}
	return provider, nil
}
