package migrations_test

import (
	"context"
	"database/sql"
	"path/filepath"
	"testing"
	"workshop/internal/infra/db/sqlite"
	"workshop/migrations"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Applies all migrations", func(t *testing.T) {
		t.Parallel()

		ctx, conn := setup(t)

		provider, err := migrations.New(conn)
		require.NoError(t, err)
		_, err = provider.Up(ctx)
		require.NoError(t, err)

		var got int
		err = conn.QueryRowContext(ctx, `
			SELECT COUNT(*)
			FROM sqlite_schema
			WHERE type = 'table'
				AND name IN ('goose_db_version', 'sightings', 'species')
		`).Scan(&got)
		require.NoError(t, err)
		assert.Equal(t, 3, got)
	})

	t.Run("Is idempotent", func(t *testing.T) {
		t.Parallel()

		ctx, conn := setup(t)

		provider, err := migrations.New(conn)
		require.NoError(t, err)
		_, err = provider.Up(ctx)
		require.NoError(t, err)
		_, err = provider.Up(ctx)
		assert.NoError(t, err)
	})
}

func setup(t *testing.T) (context.Context, *sql.DB) {
	t.Helper()

	ctx := t.Context()
	databaseURL := "file:" + filepath.Join(t.TempDir(), "migrations.db") +
		"?_pragma=foreign_keys(1)&_time_format=sqlite&_timezone=UTC"
	conn, err := sqlite.New(ctx, databaseURL)
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })

	return ctx, conn
}

func TestMediaDefault(t *testing.T) {
	ctx, conn := setup(t)
	provider, err := migrations.New(conn)
	require.NoError(t, err)
	_, err = provider.Up(ctx)
	require.NoError(t, err)

	_, err = conn.ExecContext(ctx, `
		INSERT INTO species (id, gbif_key, kingdom, phylum, class, "order", family, genus,
			species, rank, scientific_name, canonical_name, vernacular_name)
		VALUES ('species-1', 1, '', '', '', '', '', '', '', '', '', '', '');
		INSERT INTO weather (id, observed_at, temperature, apparent_temperature, humidity,
			precipitation, rain, snowfall, cloud_cover, wind_speed, wind_gusts,
			wind_direction, condition_code, condition)
		VALUES ('weather-1', '2020-01-02 12:00:00', 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, '');
		INSERT INTO sightings (id, gbif_key, occurrence_id, happened_at, basis_of_record,
			recorded_by, remarks, reference_url, locality, state_province, country,
			country_code, latitude, longitude, species_id, weather_id)
		VALUES ('sighting-1', 1, '', '2020-01-02 12:00:00', '', '', '', '', '', '', '', '', 0, 0, 'species-1', 'weather-1');
	`)
	require.NoError(t, err)

	var media string
	err = conn.QueryRowContext(ctx, "SELECT media FROM sightings WHERE id = 'sighting-1'").Scan(&media)
	require.NoError(t, err)
	assert.Equal(t, "[]", media)
}
