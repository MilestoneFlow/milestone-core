package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"milestone_core/flow"
	"milestone_core/progress"
	"milestone_core/server"
	"net/http"
	"time"
)

type FlowsResource struct {
	FlowService     flow.Service
	ProgressService progress.Service
}

type FlowCtx struct {
	id string
}

// Routes creates a REST router for the todos resource
func (rs FlowsResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List)

	r.Route("/{id}", func(r chi.Router) {
		//r.Use(rs.FlowCtx)     // lets have a users map, and lets actually load/manipulate
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
		r.Post("/continue", rs.MoveToNextStep)
		r.Post("/{stepId}/start", rs.StartStep)
		r.Post("/{stepId}/complete", rs.CompleteStep)
		r.Put("/{stepId}/data", rs.UpdateStepData)
		r.Post("/capture", rs.Capture)
	})

	return r
}

func (rs FlowsResource) List(w http.ResponseWriter, r *http.Request) {
	flows, err := rs.FlowService.List()
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, flows)
}

func (rs FlowsResource) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	// Iterate through the cursor and decode documents into User structs
	resultFlow, err := rs.FlowService.Get(idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, resultFlow)
}

func (rs FlowsResource) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	var updateInput flow.UpdateInput
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.FlowService.Update(idParam, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated flow with id: "+idParam)
}

func (rs FlowsResource) MoveToNextStep(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	step, err := rs.ProgressService.MoveToNextStep(idParam, "1", int32(time.Now().Unix()))
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, step)
	return
}

func (rs FlowsResource) StartStep(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	idStep := chi.URLParam(r, "stepId")

	step, err := rs.ProgressService.StartStep(idParam, idStep, "1", uint32(time.Now().Unix()))
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, step)
	return
}

func (rs FlowsResource) CompleteStep(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	idStep := chi.URLParam(r, "stepId")
	segmentId := r.URL.Query().Get("segmentId")

	step, err := rs.ProgressService.CompleteStep(idParam, idStep, "1", uint32(time.Now().Unix()), segmentId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, step)
	return
}

func (rs FlowsResource) UpdateStepData(w http.ResponseWriter, r *http.Request) {
	flowId := chi.URLParam(r, "id")
	stepId := chi.URLParam(r, "stepId")

	var updateInput flow.StepData
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.FlowService.UpdateStepData(flowId, stepId, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated flow with id: "+flowId)
}

func (rs FlowsResource) Capture(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")

	var updateInput flow.UpdateInput
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.FlowService.Capture(idParam, updateInput.NewSteps)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated flow with id: "+idParam)
}
