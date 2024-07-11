package authorization

import (
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	internalsql "milestone_core/shared/sql"
)

type WorkspaceBase struct {
	ID      string `json:"id" db:"id"`
	Default bool   `json:"default" db:"is_default"`
}

func GetWorkspaceIDByUserIdentifier(dbConnection *sqlx.DB, cognitoId string) (string, error) {
	var workspace WorkspaceBase
	err := dbConnection.Get(&workspace, "SELECT workspace_id as id, is_default FROM identity.workspace_user WHERE workspace_user.user_id = $1 AND workspace_user.is_default = true LIMIT 1", cognitoId)
	if errors.Is(err, sql.ErrNoRows) {
		return "", nil
	}
	if err != nil {
		return "", err
	}

	return workspace.ID, nil
}

//func GetWorkspaceIDByPublicApiToken(dbConnection *sqlx.DB, token string) (string, error) {
//	var workspaceId string
//	err := dbConnection.Get(&workspaceId, "SELECT workspace_id FROM identity.api_client WHERE identity.api_client.token = $1 LIMIT 1", token)
//	if err != nil {
//		return "", err
//	}
//
//	return workspaceId, nil
//}

func GetWorkspaceIDByPublicApiToken(dbConnection *sqlx.DB, token string) (*string, error) {
	return internalsql.FetchOne[string](dbConnection, "SELECT workspace_id FROM identity.api_client WHERE identity.api_client.token = $1 LIMIT 1", token)
}

func UserHasAccessToWorkspace(dbConnection *sqlx.DB, cognitoId string, workspaceId string) (bool, error) {
	var workspace WorkspaceBase
	err := dbConnection.Get(&workspace, "SELECT workspace_id as id, is_default FROM identity.workspace_user WHERE workspace_user.user_id = $1 AND workspace_user.workspace_id = $2 LIMIT 1", cognitoId, workspaceId)
	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateDefaultWorkspaceForUser(cognitoId string, dbConnection *sqlx.DB) (string, error) {
	query := "INSERT INTO identity.platform_user (id) VALUES ($1)"
	_, err := dbConnection.Queryx(query, cognitoId)
	if err != nil {
		return "", err
	}

	var workspaceId string
	query = "INSERT INTO identity.workspace (name, base_url) VALUES ('Default Workspace', 'https://default.workspace') RETURNING id"
	err = dbConnection.QueryRowx(query).Scan(&workspaceId)
	if err != nil {
		return "", err
	}

	_, err = dbConnection.Exec("INSERT INTO identity.workspace_user (workspace_id, user_id, is_default, role) VALUES ($1, $2, true, 'admin')", workspaceId, cognitoId)
	if err != nil {
		return "", err
	}

	return workspaceId, nil
}
