-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE SCHEMA IF NOT EXISTS identity;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP SCHEMA IF EXISTS identity;
-- +goose StatementEnd
