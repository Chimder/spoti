-- +goose Up
-- +goose StatementBegin
CREATE DATABASE IF NOT EXISTS spotify;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP DATABASE IF EXISTS spotify
-- +goose StatementEnd