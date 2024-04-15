package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"milestone_core/flow"
	"milestone_core/server"
	"milestone_core/template"
	"net/http"
)

type TemplateResource struct {
	TemplateService template.Service
}

// Routes creates a REST router for the todos resource
func (rs TemplateResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Post("/", rs.CreateFromTemplate)
	})

	return r
}

func (rs TemplateResource) List(w http.ResponseWriter, r *http.Request) {
	flows, err := rs.TemplateService.List()
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, flows)
}

func (rs TemplateResource) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	resultFlow, err := rs.TemplateService.Get(idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, resultFlow)
}

func (rs TemplateResource) CreateFromTemplate(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	override := flow.Flow{}
	err := json.NewDecoder(r.Body).Decode(&override)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	id, err := rs.TemplateService.CreateFromTemplate(workspaceId, idParam, override)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, id)
}
