-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
CREATE TYPE reward_type AS ENUM ('points', 'level', 'badge', 'custom');
CREATE TABLE game_engine.reward
(
    id           UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    workspace_id UUID         NOT NULL REFERENCES identity.workspace (id) ON DELETE CASCADE,
    key          VARCHAR(255) NOT NULL,
    name         VARCHAR(255) NOT NULL,
    type         reward_type  NOT NULL,
    metadata     JSONB,
    options      JSONB        NOT NULL DEFAULT '{}'::jsonb,
    deleted_at   TIMESTAMP,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX reward_workspace_id_idx ON game_engine.reward (workspace_id);

CREATE FUNCTION game_engine.delete_reward_record(obj_id UUID) RETURNS VOID AS
$$
BEGIN
    UPDATE game_engine.reward
    SET deleted_at = CURRENT_TIMESTAMP
    WHERE id = obj_id;
END;
$$ LANGUAGE plpgsql;

CREATE TABLE game_engine.reward_rule
(
    id         UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    reward_id  UUID         NOT NULL REFERENCES game_engine.reward (id) ON DELETE CASCADE,
    condition  VARCHAR(255) NOT NULL,
    value      JSONB        NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX reward_rule_reward_id_idx ON game_engine.reward_rule (reward_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX game_engine.reward_rule_reward_id_idx;
DROP TABLE game_engine.reward_rule;
DROP FUNCTION game_engine.delete_reward_record(UUID);
DROP INDEX game_engine.reward_workspace_id_idx;
DROP TABLE game_engine.reward;
DROP TYPE reward_type;
-- +goose StatementEnd
