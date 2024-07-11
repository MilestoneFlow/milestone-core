-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE SCHEMA IF NOT EXISTS game_engine;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP SCHEMA IF EXISTS game_engine;
-- +goose StatementEnd
