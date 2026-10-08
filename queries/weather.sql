-- name: WeatherFind :one
SELECT * FROM weather WHERE id = ?;

-- name: WeatherCreate :one
INSERT INTO weather (
	id,
	observed_at,
	temperature,
	apparent_temperature
) VALUES (?, ?, ?, ?) RETURNING *;
