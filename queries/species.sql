-- name: SpeciesFind :one
SELECT * FROM species WHERE id = ?;

-- name: SpeciesFindByGbifKey :one
SELECT * FROM species WHERE gbif_key = ?;

-- name: SpeciesList :many
SELECT * FROM species;

-- name: SpeciesCreate :one
INSERT INTO species (id, "gbif_key", "kingdom", "phylum", "class", "order", "family", "genus", "species", "rank", "scientific_name", "canonical_name", "vernacular_name")
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) RETURNING *;
