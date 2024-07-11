package apiclient

type ApiClient struct {
	ID          string `json:"id,omitempty" db:"id"`
	WorkspaceID string `json:"workspaceId" db:"workspace_id"`
	Token       string `json:"token" db:"token"`
}
