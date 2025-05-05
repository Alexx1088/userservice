-- +goose Up
ALTER TABLE users ADD COLUMN surname TEXT NOT NULL DEFAULT '';

