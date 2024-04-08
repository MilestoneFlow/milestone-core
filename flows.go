package main

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"milestone_core/awsinternal"
	"milestone_core/flow"
	"milestone_core/progress"
	"milestone_core/server"
	"net/http"
	"path/filepath"
	"time"
)

type FlowsResource struct {
	FlowService     flow.Service
	ProgressService progress.Service
	Ctx             FlowCtx
}

type FlowCtx struct {
	id string
}

// Routes creates a REST router for the todos resource
func (rs FlowsResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..
	//r.Use(rs.Ctx)

	r.Get("/", rs.List)

	r.Route("/{id}", func(r chi.Router) {
		r.Post("/{stepId}/media", rs.UploadMediaFile)
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
		r.Post("/continue", rs.MoveToNextStep)
		r.Put("/{stepId}", rs.UpdateStep)
		r.Post("/capture", rs.Capture)
		r.Get("/analytics", rs.GetFlowAnalytics)
	})

	return r
}

func (rs FlowsResource) List(w http.ResponseWriter, r *http.Request) {
	workspace := r.Context().Value("workspace").(string)
	flows, err := rs.FlowService.List(workspace)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, flows)
}

func (rs FlowsResource) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspace := r.Context().Value("workspace").(string)

	// Iterate through the cursor and decode documents into User structs
	resultFlow, err := rs.FlowService.Get(workspace, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, resultFlow)
}

func (rs FlowsResource) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspace := r.Context().Value("workspace").(string)

	var updateInput flow.UpdateInput
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.FlowService.Update(workspace, idParam, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated flow with id: "+idParam)
}

func (rs FlowsResource) MoveToNextStep(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspace := r.Context().Value("workspace").(string)

	step, err := rs.ProgressService.MoveToNextStep(workspace, idParam, "1", int32(time.Now().Unix()))
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, step)
	return
}

func (rs FlowsResource) UpdateStep(w http.ResponseWriter, r *http.Request) {
	flowId := chi.URLParam(r, "id")
	stepId := chi.URLParam(r, "stepId")
	workspace := r.Context().Value("workspace").(string)

	inputFlow, err := rs.FlowService.Get(workspace, flowId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	var updateInput flow.Step
	err = json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.FlowService.UpdateStep(workspace, inputFlow, stepId, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated flow with id: "+flowId)
}

func (rs FlowsResource) Capture(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspace := r.Context().Value("workspace").(string)

	var updateInput flow.UpdateInput
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	newId, err := rs.FlowService.Capture(workspace, idParam, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, newId)
}

func (rs FlowsResource) UploadMediaFile(w http.ResponseWriter, r *http.Request) {
	flowId := chi.URLParam(r, "id")
	stepId := chi.URLParam(r, "stepId")
	workspace := r.Context().Value("workspace").(string)

	resFlow, err := rs.FlowService.Get(workspace, flowId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}
	found := false
	for _, step := range resFlow.Steps {
		if step.StepID == stepId {
			found = true
			break
		}
	}
	if !found {
		server.SendBadRequestErrorJson(w, errors.New("step not found"))
		return
	}

	err = r.ParseMultipartForm(10 << 20) // Max upload size ~10MB
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	file, headers, err := r.FormFile("uploadedFile")
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}
	defer file.Close()

	filename := uuid.New().String() + filepath.Ext(headers.Filename)
	err = awsinternal.UploadToS3(context.TODO(), filename, file)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, struct {
		FileName string `json:"fileName"`
	}{
		FileName: "https://milestone-uploaded-flows-media.s3.amazonaws.com/step_media/" + filename,
	})
}

func (rs FlowsResource) GetFlowAnalytics(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspace := r.Context().Value("workspace").(string)

	analytics, err := rs.FlowService.GetFlowAnalytics(workspace, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, analytics)
}
