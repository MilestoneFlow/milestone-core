-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE game_engine.user_wallet
(
    id              UUID PRIMARY KEY      DEFAULT uuid_generate_v4(),
    workspace_id    UUID         NOT NULL REFERENCES identity.workspace (id),
    user_id         varchar(255) NOT NULL,
    current_balance INT          NOT NULL DEFAULT 0,
    created_at      TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TYPE game_engine.user_wallet_transaction_type AS ENUM ('deposit', 'withdrawal');
CREATE TABLE game_engine.user_wallet_transaction
(
    id               UUID PRIMARY KEY                                  DEFAULT uuid_generate_v4(),
    user_wallet_id   UUID                                     NOT NULL REFERENCES game_engine.user_wallet (id),
    amount           INT                                      NOT NULL,
    transaction_type game_engine.user_wallet_transaction_type NOT NULL,
    track_data       JSONB                                    NOT NULL DEFAULT '{}'::JSONB,
    created_at       TIMESTAMP                                NOT NULL DEFAULT NOW()
);
CREATE INDEX user_wallet_transaction_user_wallet_id_idx ON game_engine.user_wallet_transaction (user_wallet_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX game_engine.user_wallet_transaction_user_wallet_id_idx;
DROP TABLE game_engine.user_wallet_transaction;
DROP TYPE game_engine.user_wallet_transaction_type;
DROP TABLE game_engine.user_wallet;
-- +goose StatementEnd
