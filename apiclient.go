package main

import (
	"github.com/go-chi/chi/v5"
	"milestone_core/apiclient"
	"milestone_core/server"
	"net/http"
)

type ApiClientResource struct {
	Service apiclient.Service
}

// Routes creates a REST router for the todos resource
func (rs ApiClientResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Post("/", rs.Create)
	r.Get("/{id}", rs.Get)
	r.Delete("/{id}", rs.Delete)

	return r
}

func (rs ApiClientResource) List(w http.ResponseWriter, r *http.Request) {
	workspace := r.Context().Value("workspace").(string)
	apiClients, err := rs.Service.List(workspace)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, apiClients)
}

func (rs ApiClientResource) Get(w http.ResponseWriter, r *http.Request) {
	workspace := r.Context().Value("workspace").(string)
	id := chi.URLParam(r, "id")
	apiClient, err := rs.Service.Get(workspace, id)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, apiClient)
}

func (rs ApiClientResource) Create(w http.ResponseWriter, r *http.Request) {
	workspace := r.Context().Value("workspace").(string)

	token, err := rs.Service.Create(workspace)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, struct {
		Token string `json:"token"`
	}{Token: token})
}

func (rs ApiClientResource) Delete(w http.ResponseWriter, r *http.Request) {
	workspace := r.Context().Value("workspace").(string)
	id := chi.URLParam(r, "id")
	err := rs.Service.Delete(workspace, id)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "Deleted")
}
