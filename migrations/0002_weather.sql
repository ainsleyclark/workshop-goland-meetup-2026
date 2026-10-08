-- +goose Up
-- +goose StatementBegin
CREATE TABLE weather (
	id UUID PRIMARY KEY,
	created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
	observed_at DATETIME NOT NULL,
	temperature INTEGER,
	apparent_temperature INTEGER

);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS weather;
-- +goose StatementEnd
