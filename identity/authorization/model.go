package authorization

type UserData struct {
	WorkspaceID string
	UserID      string
	Email       string
}

type PublicApiClientData struct {
	WorkspaceID string
	Token       string
}
