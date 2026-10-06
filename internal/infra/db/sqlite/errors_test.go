package sqlite_test

import (
	"errors"
	"fmt"
	"testing"
	"workshop/internal/infra/db/dbtest"
	"workshop/internal/infra/db/sqlite"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIsUniqueViolation(t *testing.T) {
	t.Parallel()

	ctx, conn, teardown := dbtest.Setup(t)
	t.Cleanup(teardown)

	_, err := conn.ExecContext(ctx, "CREATE TABLE items (id INTEGER PRIMARY KEY, name TEXT NOT NULL UNIQUE, quantity INTEGER CHECK (quantity > 0))")
	require.NoError(t, err)
	_, err = conn.ExecContext(ctx, "INSERT INTO items VALUES (1, 'existing', 1)")
	require.NoError(t, err)

	constraintError := func(t *testing.T, query string) error {
		t.Helper()
		_, err := conn.ExecContext(ctx, query)
		require.Error(t, err)
		return err
	}

	uniqueErr := constraintError(t, "INSERT INTO items VALUES (2, 'existing', 1)")
	primaryKeyErr := constraintError(t, "INSERT INTO items VALUES (1, 'new', 1)")
	notNullErr := constraintError(t, "INSERT INTO items VALUES (2, NULL, 1)")
	checkErr := constraintError(t, "INSERT INTO items VALUES (2, 'new', -1)")
	sqlErr := constraintError(t, "INSERT INTO missing_table VALUES (1)")

	tt := map[string]struct {
		input error
		want  bool
	}{
		"Nil error":                        {input: nil, want: false},
		"Unrelated error":                  {input: errors.New("unrelated"), want: false},
		"Matching error message":           {input: errors.New(uniqueErr.Error()), want: false},
		"Unique constraint":                {input: uniqueErr, want: true},
		"Wrapped unique constraint":        {input: fmt.Errorf("insert failed: %w", uniqueErr), want: true},
		"Nested wrapped unique constraint": {input: fmt.Errorf("save failed: %w", fmt.Errorf("insert failed: %w", uniqueErr)), want: true},
		"Joined unique constraint":         {input: errors.Join(errors.New("unrelated"), uniqueErr), want: true},
		"Primary key constraint":           {input: primaryKeyErr, want: false},
		"Wrapped other constraint":         {input: fmt.Errorf("insert failed: %w", primaryKeyErr), want: false},
		"Not null constraint":              {input: notNullErr, want: false},
		"Check constraint":                 {input: checkErr, want: false},
		"SQL error":                        {input: sqlErr, want: false},
	}

	for name, test := range tt {
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			assert.Equal(t, test.want, sqlite.IsUniqueViolation(test.input))
		})
	}
}
