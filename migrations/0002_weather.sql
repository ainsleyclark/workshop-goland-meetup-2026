-- +goose Up
-- +goose StatementBegin
CREATE TABLE weather (
	id                   UUID PRIMARY KEY,
	created_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at           DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	observed_at          DATETIME NOT NULL,
	temperature          REAL NOT NULL,
	apparent_temperature REAL NOT NULL,
	humidity             INTEGER NOT NULL CHECK (humidity BETWEEN 0 AND 100),
	precipitation        REAL NOT NULL CHECK (precipitation >= 0),
	rain                 REAL NOT NULL CHECK (rain >= 0),
	snowfall             REAL NOT NULL CHECK (snowfall >= 0),
	cloud_cover          INTEGER NOT NULL CHECK (cloud_cover BETWEEN 0 AND 100),
	wind_speed           REAL NOT NULL,
	wind_gusts           REAL NOT NULL,
	wind_direction       INTEGER NOT NULL CHECK (wind_direction BETWEEN 0 AND 360),
	condition_code       INTEGER NOT NULL,
	condition            TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS weather;
-- +goose StatementEnd
