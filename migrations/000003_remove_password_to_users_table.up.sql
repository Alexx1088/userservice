-- +goose Up
ALTER TABLE users DROP password;

-- +goose Down
ALTER TABLE users ADD password;