-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE IF NOT EXISTS identity.api_client
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    token        VARCHAR(255) NOT NULL,
    workspace_id UUID         NOT NULL REFERENCES identity.workspace (id),
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE IF EXISTS identity.api_client;
-- +goose StatementEnd
