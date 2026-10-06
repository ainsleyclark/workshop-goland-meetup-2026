package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

// New opens a *sql.DB against the given URL.
//
// Possible values for URL:
// - :memory:		- In memory DB.
// - file:"/path"	- SQL file path.
func New(ctx context.Context, url string) (*sql.DB, error) {
	if url == "" {
		return nil, fmt.Errorf("sqlite: url is empty")
	}

	db, err := sql.Open("sqlite", url)
	if err != nil {
		return nil, err
	}

	if err = db.PingContext(ctx); err != nil {
		_ = db.Close() //nolint
		return nil, fmt.Errorf("sqlite: db ping failed: %v", err)
	}

	return db, nil
}
