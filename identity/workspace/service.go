package workspace

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"milestone_core/identity/users"
)

type Service struct {
	DbConnection *sqlx.DB
	UsersService users.Service
}

func (s Service) Get(workspaceId string) (*Workspace, error) {
	var workspace Workspace
	err := s.DbConnection.Get(&workspace, "SELECT w.id as id, w.name as name, w.base_url as base_url, coalesce(w.invite_token, '') as invite_token FROM identity.workspace w WHERE w.id = $1", workspaceId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (s Service) FetchAllForUser(userId string) ([]Workspace, error) {
	var workspaces []Workspace
	err := s.DbConnection.Select(&workspaces, `
		SELECT w.id as id, w.name as name, w.base_url as base_url, coalesce(w.invite_token, '') as invite_token
		FROM identity.workspace w
		JOIN identity.workspace_user wu ON w.id = wu.workspace_id 
		WHERE wu.user_id = $1
		ORDER BY wu.is_default DESC, w.created_at DESC
		`, userId)

	return workspaces, err
}

func (s Service) GetUsers(workspaceId string) (*WorkspaceUsers, error) {
	workspaceActiveUsers, err := s.UsersService.GetWorkspaceUsers(workspaceId)
	if err != nil {
		return nil, err
	}
	if workspaceActiveUsers == nil {
		workspaceActiveUsers = make([]users.User, 0)
	}

	invitedUsers, err := s.UsersService.GetWorkspaceInvitedUsers(workspaceId)
	if err != nil {
		return nil, err
	}
	if invitedUsers == nil {
		invitedUsers = make([]users.InvitedUser, 0)
	}

	return &WorkspaceUsers{
		ActiveUsers:  workspaceActiveUsers,
		InvitedUsers: invitedUsers,
	}, nil
}

func (s Service) GetByUserId(userId string) (*Workspace, error) {
	var workspace Workspace
	err := s.DbConnection.Get(&workspace, "SELECT w.id as id, w.name as name, w.base_url as base_url, coalesce(w.invite_token, '') as invite_token FROM identity.workspace w JOIN identity.workspace_user wu ON w.id = wu.workspace_id WHERE wu.user_id = $1", userId)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &workspace, nil
}

func (s Service) CreateForUser(workspace Workspace, userId string) error {
	inviteToken, err := s.generateInviteToken()
	if err != nil {
		return err
	}

	workspace.InviteToken = inviteToken
	query := "INSERT INTO identity.workspace (name, base_url, invite_token) VALUES ($1, $2, $3) RETURNING id"
	row := s.DbConnection.QueryRow(query)

	var workspaceId string
	if err = row.Scan(&workspaceId); err != nil {
		return err
	}

	_, err = s.DbConnection.Exec("INSERT INTO identity.workspace_user (workspace_id, user_id, is_default) VALUES ($1, $2, true)", workspaceId, userId)
	return err
}

func (s Service) Update(id string, workspace Workspace) error {
	_, err := s.DbConnection.Exec("UPDATE identity.workspace SET name = $1, base_url = $2 WHERE id = $3", workspace.Name, workspace.BaseURL, id)
	return err
}

func (s Service) InviteUsers(workspaceId string, userEmails []string) error {
	workspace, err := s.Get(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return errors.New("workspace not found")
	}

	tx, err := s.DbConnection.Begin()
	if err != nil {
		return err
	}

	for _, email := range userEmails {
		inviteToken, err := s.generateInviteToken()
		if err != nil {
			err := tx.Rollback()
			return err
		}

		_, err = tx.Exec("INSERT INTO identity.workspace_user_invite (workspace_id, email, token) VALUES ($1, $2, $3)", workspaceId, email, inviteToken)
		if err != nil {
			err := tx.Rollback()
			return err
		}
	}

	return tx.Commit()
}

func (s Service) RemoveUser(workspaceId string, userId string) error {
	workspace, err := s.Get(workspaceId)
	if err != nil {
		return err
	}
	if workspace == nil {
		return errors.New("workspace not found")
	}

	_, err = s.DbConnection.Query("DELETE FROM identity.workspace_user WHERE workspace_id = $1 AND user_id = $2", workspaceId, userId)
	if err != nil {
		return err
	}

	return nil
}

func (s Service) RefreshInviteToken(workspaceId string) (string, error) {
	inviteURL, err := s.generateInviteToken()
	if err != nil {
		return "", err
	}

	_, err = s.DbConnection.Exec("UPDATE identity.workspace SET invite_token = $1 WHERE id = $2", inviteURL, workspaceId)
	return inviteURL, err
}

func (s Service) generateInviteToken() (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, 64)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = letters[b%byte(len(letters))]
	}
	return string(bytes), nil
}

type WorkspaceUsers struct {
	ActiveUsers  []users.User        `json:"activeUsers"`
	InvitedUsers []users.InvitedUser `json:"invitedUsers"`
}
