-- name: WeatherFind :one
SELECT * FROM weather WHERE id = ?;

-- name: WeatherCreate :one
INSERT INTO weather (
	id,
	observed_at,
	temperature,
	apparent_temperature,
	humidity,
	precipitation,
	rain,
	snowfall,
	cloud_cover,
	wind_speed,
	wind_gusts,
	wind_direction,
	condition_code,
	condition
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;
