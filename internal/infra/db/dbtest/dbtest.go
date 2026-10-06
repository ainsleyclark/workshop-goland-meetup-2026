package dbtest

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"workshop/internal/infra/db/sqlite"
	"workshop/migrations"

	"github.com/stretchr/testify/require"
)

// Setup creates a new SQLite *sql.DB with all migrations applied,
// ready for testing. The DB is file-backed under t.TempDir() so it
// is torn down automatically when the test ends.
func Setup(t *testing.T) (context.Context, *sql.DB, func()) {
	t.Helper()

	ctx := t.Context()

	// Create a test file placed in the test temporary directory.
	url := "file:" + filepath.Join(t.TempDir(), "dbtest.db") + "?_time_format=sqlite&_timezone=UTC"

	conn, err := sqlite.New(ctx, url)
	require.NoError(t, err, "db test: opening sqlite database")

	migrator, err := migrations.New(conn)
	require.NoError(t, err, "db test: creating migration provider")

	_, err = migrator.Up(ctx)
	require.NoError(t, err, "db test: applying migrations")

	cleanup := func() {
		require.NoError(t, conn.Close(), "db test:closing sqlite database")
		require.NoError(t, migrator.Close(), "db test:closing migrator")
	}

	return ctx, conn, cleanup
}
