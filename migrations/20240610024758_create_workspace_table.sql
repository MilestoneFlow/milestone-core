-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE identity.workspace
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    old_mongo_id VARCHAR(255),
    name         VARCHAR(255) NOT NULL,
    base_url     VARCHAR(255),
    invite_token VARCHAR(255),
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE identity.workspace_user
(
    workspace_id UUID          NOT NULL REFERENCES identity.workspace (id) ON DELETE CASCADE,
    user_id      VARCHAR(255) NOT NULL REFERENCES identity.platform_user (id) ON DELETE CASCADE,
    role         VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE identity.workspace_user
    ADD PRIMARY KEY (workspace_id, user_id);

CREATE TABLE identity.workspace_user_invite
(
    workspace_id UUID          NOT NULL REFERENCES identity.workspace (id) ON DELETE CASCADE,
    email        VARCHAR(255) NOT NULL,
    token        VARCHAR(255) NOT NULL,
    created_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP
);
ALTER TABLE identity.workspace_user_invite
    ADD PRIMARY KEY (workspace_id, email);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP TABLE identity.workspace_user;
DROP TABLE identity.workspace_user_invite;
DROP TABLE identity.workspace;
-- +goose StatementEnd
