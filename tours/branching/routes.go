package branching

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"milestone_core/shared/server"
	"milestone_core/tours/flows"
	"net/http"
)

type BranchingResource struct {
	BranchingService flows.BranchingService
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
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	branches, err := rs.BranchingService.List(workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, branches)
}

func (rs BranchingResource) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	// Iterate through the cursor and decode documents into User structs
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	branch, err := rs.BranchingService.Get(workspaceId, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, branch)
}

func (rs BranchingResource) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	var updateInput flows.Branching
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	err = rs.BranchingService.Update(workspaceId, idParam, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated branch with id: "+idParam)
}

func (rs BranchingResource) Create(w http.ResponseWriter, r *http.Request) {
	var branch flows.Branching
	err := json.NewDecoder(r.Body).Decode(&branch)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	res, err := rs.BranchingService.Create(workspaceId, branch)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, res)
}
