package workspace

type Workspace struct {
	ID          string `json:"id"  db:"id"`
	Name        string `json:"name" db:"name"`
	BaseURL     string `json:"baseUrl"  db:"base_url"`
	InviteToken string `json:"inviteToken"  db:"invite_token"`
}
