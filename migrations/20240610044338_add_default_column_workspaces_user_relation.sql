-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TABLE identity.workspace_user ADD COLUMN is_default BOOLEAN NOT NULL DEFAULT FALSE;
CREATE UNIQUE INDEX idx_unique_is_default_per_user ON identity.workspace_user(user_id)
    WHERE is_default = TRUE;

-- Common Table Expression (CTE) to get row numbers
WITH ranked_users AS (
    SELECT
        user_id,
        workspace_id,
        ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created_at) AS rn
    FROM
        identity.workspace_user
)

-- Update statement to set is_default to TRUE for the first entry for each user_id
UPDATE identity.workspace_user wu
SET is_default = TRUE
FROM ranked_users ru
WHERE wu.user_id = ru.user_id
  AND wu.workspace_id = ru.workspace_id
  AND ru.rn = 1;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
DROP INDEX idx_unique_is_default_per_user;
ALTER TABLE identity.workspace_user DROP COLUMN is_default;
-- +goose StatementEnd
