-- +goose Up
-- +goose StatementBegin
CREATE TABLE sightings (
	id                                UUID PRIMARY KEY,
	created_at                        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at                        DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	gbif_key                          INTEGER NOT NULL UNIQUE,
	occurrence_id                     TEXT NOT NULL,
	happened_at                       DATETIME NOT NULL,
	basis_of_record                   TEXT NOT NULL,
	recorded_by                       TEXT NOT NULL,
	individual_count                  INTEGER,
	remarks                           TEXT NOT NULL,
	reference_url                     TEXT NOT NULL,
	locality                          TEXT NOT NULL,
	state_province                    TEXT NOT NULL,
	country                           TEXT NOT NULL,
	country_code                      TEXT NOT NULL,
	latitude                          REAL NOT NULL CHECK (latitude BETWEEN -90 AND 90),
	longitude                         REAL NOT NULL CHECK (longitude BETWEEN -180 AND 180),
	coordinate_uncertainty_in_metres  REAL,
	elevation_in_metres               REAL,
	species_id                        UUID NOT NULL,
	weather_id                        UUID NOT NULL,
	media                             TEXT NOT NULL DEFAULT '[]'
		CHECK (json_valid(media) AND json_type(media) = 'array'),
	FOREIGN KEY (species_id) REFERENCES species (id),
	-- weather is created in 0002; SQLite resolves foreign keys when rows
	-- are written, not when the table is declared.
	FOREIGN KEY (weather_id) REFERENCES weather (id)
);

CREATE TABLE species (
	id UUID PRIMARY KEY,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	"gbif_key" INTEGER NOT NULL UNIQUE,
	"kingdom" TEXT NOT NULL,
	"phylum" TEXT NOT NULL,
	"class" TEXT NOT NULL,
	"order" TEXT NOT NULL,
	"family" TEXT NOT NULL,
	"genus" TEXT NOT NULL,
	"species" TEXT NOT NULL,
	"rank" TEXT NOT NULL,
	"scientific_name" TEXT NOT NULL,
	"canonical_name" TEXT NOT NULL,
	"vernacular_name" TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS sightings;
DROP TABLE IF EXISTS species;
-- +goose StatementEnd
