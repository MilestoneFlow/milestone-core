-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE game_engine.user_events
(
    id           UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    workspace_id UUID         NOT NULL REFERENCES identity.workspace (id),
    user_id      VARCHAR(255) NOT NULL,
    event_id     UUID         NOT NULL REFERENCES game_engine.event (id),
    metadata     JSONB,
    created_at   TIMESTAMP    NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE game_engine.user_events;
-- +goose StatementEnd
