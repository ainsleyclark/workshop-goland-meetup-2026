package sqlite

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Empty URL", func(t *testing.T) {
		t.Parallel()

		db, err := New(t.Context(), "")

		require.EqualError(t, err, "sqlite: url is empty")
		assert.Nil(t, db)
	})

	t.Run("Opens a usable database", func(t *testing.T) {
		t.Parallel()

		tt := map[string]string{
			"In memory": ":memory:",
			"File":      filepath.Join(t.TempDir(), "test.db"),
		}

		for name, url := range tt {
			t.Run(name, func(t *testing.T) {
				t.Parallel()

				db, err := New(t.Context(), url)
				require.NoError(t, err)
				t.Cleanup(func() { require.NoError(t, db.Close()) })

				var got int
				err = db.QueryRowContext(t.Context(), "SELECT 1").Scan(&got)
				require.NoError(t, err)
				assert.Equal(t, 1, got)
			})
		}
	})

	t.Run("Ping failure", func(t *testing.T) {
		t.Parallel()

		ctx, cancel := context.WithCancel(t.Context())
		cancel()

		db, err := New(ctx, ":memory:")

		require.Error(t, err)
		assert.ErrorContains(t, err, "sqlite: db ping failed")
		assert.ErrorContains(t, err, context.Canceled.Error())
		assert.Nil(t, db)
	})
}
