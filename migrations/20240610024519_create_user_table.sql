-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS identity.platform_user (
    id VARCHAR(255) PRIMARY KEY NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS identity.platform_user;
-- +goose StatementEnd
