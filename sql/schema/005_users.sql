-- +goose Up
ALTER TABLE users ADD COLUMN is_chirpy_red BOOLEAN NOT NULL DEFAULT FALSE;

-- +goose Down
ALTER TABLE IF EXISTS users DROP COLUMN is_chirpy_red;