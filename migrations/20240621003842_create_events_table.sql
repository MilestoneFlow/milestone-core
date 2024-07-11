-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TABLE game_engine.event
(
    id           UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    workspace_id UUID         NOT NULL REFERENCES identity.workspace (id) ON DELETE CASCADE,
    key          VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    deleted_at   TIMESTAMP,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX event_workspace_id_idx ON game_engine.event (workspace_id);

CREATE FUNCTION game_engine.delete_event_record(obj_id UUID) RETURNS VOID AS
$$
BEGIN
    UPDATE game_engine.event
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = obj_id;
END;
$$ LANGUAGE plpgsql;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP FUNCTION game_engine.delete_event_record(UUID);
DROP INDEX game_engine.event_workspace_id_idx;
DROP TABLE game_engine.event;
-- +goose StatementEnd
