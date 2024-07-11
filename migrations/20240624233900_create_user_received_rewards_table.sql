-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE game_engine.user_received_rewards
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id    varchar(255) NOT NULL,
    reward_id  UUID         NOT NULL REFERENCES game_engine.reward (id),
    created_at TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

DROP TABLE game_engine.user_received_rewards;
-- +goose StatementEnd
