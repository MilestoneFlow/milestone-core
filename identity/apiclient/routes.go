package apiclient

import (
	"github.com/go-chi/chi/v5"
	"milestone_core/shared/server"
	"net/http"
)

type ApiClientResource struct {
	Service Service
}

func (rs ApiClientResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Post("/", rs.Create)
	r.Get("/{id}", rs.Get)
	r.Delete("/{id}", rs.Delete)

	return r
}

func (rs ApiClientResource) List(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	apiClients, err := rs.Service.List(workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, apiClients)
}

func (rs ApiClientResource) Get(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	id := chi.URLParam(r, "id")
	apiClient, err := rs.Service.Get(workspaceId, id)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, apiClient)
}

func (rs ApiClientResource) Create(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	token, err := rs.Service.Create(workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, struct {
		Token string `json:"token"`
	}{Token: token})
}

func (rs ApiClientResource) Delete(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	id := chi.URLParam(r, "id")
	err := rs.Service.Delete(workspaceId, id)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "Deleted")
}
