package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"milestone_core/flow"
	"milestone_core/server"
	"net/http"
)

type BranchingResource struct {
	BranchingService flow.BranchingService
}

type BranchingCtx struct {
	id string
}

// Routes creates a REST router for the todos resource
func (rs BranchingResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)
	r.Post("/", rs.Create)

	r.Route("/{id}", func(r chi.Router) {
		//r.Use(rs.FlowCtx)     // lets have a users map, and lets actually load/manipulate
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
	})

	return r
}

func (rs BranchingResource) List(w http.ResponseWriter, r *http.Request) {
	workspace := r.Context().Value("workspace").(string)

	branches, err := rs.BranchingService.List(workspace)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, branches)
}

func (rs BranchingResource) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	// Iterate through the cursor and decode documents into User structs
	workspace := r.Context().Value("workspace").(string)
	branch, err := rs.BranchingService.Get(workspace, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, branch)
}

func (rs BranchingResource) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	var updateInput flow.Branching
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	workspace := r.Context().Value("workspace").(string)
	err = rs.BranchingService.Update(workspace, idParam, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated branch with id: "+idParam)
}

func (rs BranchingResource) Create(w http.ResponseWriter, r *http.Request) {
	var branch flow.Branching
	err := json.NewDecoder(r.Body).Decode(&branch)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	workspace := r.Context().Value("workspace").(string)
	res, err := rs.BranchingService.Create(workspace, branch)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, res)
}
