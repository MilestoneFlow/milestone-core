package users

type User struct {
	ID        string      `json:"id" db:"id"`
	CreatedAt string      `json:"createdAt" db:"created_at"`
	Details   UserDetails `json:"details"`
}

type UserDetails struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type InvitedUser struct {
	WorkspaceID string `json:"workspaceId" db:"workspace_id"`
	Email       string `json:"email" db:"email"`
	Token       string `json:"-" db:"token"`
	CreatedAt   string `json:"createdAt" db:"created_at"`
}
