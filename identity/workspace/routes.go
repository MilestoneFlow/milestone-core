package workspace

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"milestone_core/identity/authorization"
	"milestone_core/server"
	"net/http"
)

type Resource struct {
	Service Service
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.Get)
	r.Post("/", rs.Create)
	r.Put("/", rs.Update)
	r.Post("/refresh-link", rs.RefreshLink)
	r.Post("/invite-members", rs.InviteMembers)
	r.Post("/remove-member", rs.RemoveMember)

	return r
}

func (rs Resource) Get(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	workspace, err := rs.Service.Get(workspaceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if workspace == nil {
		http.Error(w, "Workspace not found", http.StatusNotFound)
		return
	}

	server.SendJson(w, workspace)
}

func (rs Resource) Create(w http.ResponseWriter, r *http.Request) {
	var workspace Workspace
	err := json.NewDecoder(r.Body).Decode(&workspace)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	userData := r.Context().Value("user").(authorization.UserData)
	currentWorkspace, err := rs.Service.GetByUserIdentifier(userData.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if currentWorkspace != nil {
		server.SendBadRequestErrorJson(w, errors.New("User already has a workspace"))
		return
	}

	err = rs.Service.CreateForUser(workspace, userData.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server.SendMessageJson(w, "Workspace created")
}

func (rs Resource) Update(w http.ResponseWriter, r *http.Request) {
	var workspace Workspace
	err := json.NewDecoder(r.Body).Decode(&workspace)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	err = rs.Service.Update(workspaceId, workspace)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server.SendMessageJson(w, "Workspace updated")
}

func (rs Resource) InviteMembers(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	var userIdentifier struct {
		Emails []string `json:"emails"`
	}
	err := json.NewDecoder(r.Body).Decode(&userIdentifier)

	err = rs.Service.AddUserIdentifiers(workspaceId, userIdentifier.Emails)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server.SendMessageJson(w, "Users invited to workspace")
}

func (rs Resource) RemoveMember(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	var userIdentifier struct {
		Email string `json:"email"`
	}
	err := json.NewDecoder(r.Body).Decode(&userIdentifier)

	err = rs.Service.RemoveUserIdentifier(workspaceId, userIdentifier.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	server.SendMessageJson(w, "User removed from workspace")
}

func (rs Resource) RefreshLink(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	token, err := rs.Service.RefreshInviteToken(workspaceId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	baseUrl := "https://dashboard.milestoneflow.io/invite/"
	server.SendJson(w, map[string]string{"link": baseUrl + token})
}
