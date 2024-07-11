-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE SCHEMA webhooks;

CREATE TABLE webhooks.webhook_endpoints
(
    id           UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    workspace_id UUID         NOT NULL REFERENCES identity.workspace (id),
    name         VARCHAR(255) NOT NULL,
    url          TEXT         NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at   TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TYPE webhooks.message_status AS ENUM ('pending', 'sent', 'failed');

CREATE TABLE webhooks.webhook_messages
(
    id           UUID PRIMARY KEY                 DEFAULT uuid_generate_v4(),
    workspace_id UUID                    NOT NULL REFERENCES identity.workspace (id),
    endpoint_id  UUID                    NOT NULL REFERENCES webhooks.webhook_endpoints (id),
    payload      JSONB                   NOT NULL,
    created_at   TIMESTAMP               NOT NULL DEFAULT NOW(),
    sent_at      TIMESTAMP                        DEFAULT NULL,
    status       webhooks.message_status NOT NULL DEFAULT 'pending'
);

CREATE INDEX idx_webhook_messages_status ON webhooks.webhook_messages (status, created_at);

CREATE TABLE webhooks.webhook_delivery_logs
(
    id               UUID PRIMARY KEY   DEFAULT uuid_generate_v4(),
    message_id       UUID      NOT NULL REFERENCES webhooks.webhook_messages (id),
    delivery_attempt INT       NOT NULL,
    delivered_at     TIMESTAMP NOT NULL DEFAULT NOW(),
    response_status  INT,
    response_body    TEXT,
    success          BOOLEAN   NOT NULL DEFAULT FALSE,
    created_at       TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE webhooks.webhook_delivery_logs;
DROP INDEX webhooks.idx_webhook_messages_status;
DROP TABLE webhooks.webhook_messages;
DROP TYPE webhooks.message_status;
DROP TABLE webhooks.webhook_endpoints;
DROP SCHEMA webhooks;
-- +goose StatementEnd
