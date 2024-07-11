package rest

import "net/http"

type AppRequest struct {
	*http.Request
	WorkspaceID string
}

func NewAppRequest(r *http.Request, workspaceId string) *AppRequest {
	return &AppRequest{
		Request:     r,
		WorkspaceID: workspaceId,
	}
}
