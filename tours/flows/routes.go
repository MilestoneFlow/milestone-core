package flows

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"milestone_core/shared/awsinternal"
	"milestone_core/shared/server"
	"net/http"
	"path/filepath"
)

type FlowsResource struct {
	FlowService Service
	Analytics   Analytics
	Ctx         FlowCtx
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
		r.Delete("/", rs.Archive)
		r.Put("/{stepId}", rs.UpdateStep)
		r.Post("/capture", rs.Capture)
		r.Get("/analytics", rs.GetFlowAnalytics)
		r.Post("/publish", rs.Publish)
		r.Post("/unpublish", rs.Unpublish)
		r.Get("/possible-depends-on-list", rs.GetPossibleDependsOnListForFlow)
	})

	r.Route("/archive", func(r chi.Router) {
		r.Get("/", rs.ListArchive)
		r.Post("/{id}", rs.RestoreFlow)
	})

	return r
}

func (rs FlowsResource) List(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	flows, err := rs.FlowService.List(workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, flows)
}

func (rs FlowsResource) Get(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	// Iterate through the cursor and decode documents into User structs
	resultFlow, err := rs.FlowService.Get(workspaceId, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, resultFlow)
}

func (rs FlowsResource) Update(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	var updateInput UpdateInput
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.FlowService.Update(workspaceId, idParam, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated flow with id: "+idParam)
}

func (rs FlowsResource) UpdateStep(w http.ResponseWriter, r *http.Request) {
	flowId := chi.URLParam(r, "id")
	stepId := chi.URLParam(r, "stepId")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	inputFlow, err := rs.FlowService.Get(workspaceId, flowId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	var updateInput Step
	err = json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.FlowService.UpdateStep(workspaceId, inputFlow, stepId, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "updated flow with id: "+flowId)
}

func (rs FlowsResource) Capture(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	var updateInput UpdateInput
	err := json.NewDecoder(r.Body).Decode(&updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	newId, err := rs.FlowService.Capture(workspaceId, idParam, updateInput)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, newId)
}

func (rs FlowsResource) UploadMediaFile(w http.ResponseWriter, r *http.Request) {
	flowId := chi.URLParam(r, "id")
	stepId := chi.URLParam(r, "stepId")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	resFlow, err := rs.FlowService.Get(workspaceId, flowId)
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
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	flow, err := rs.FlowService.Get(workspaceId, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	analytics, err := rs.Analytics.GetFlowAnalytics(flow)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, analytics)
}

func (rs FlowsResource) Publish(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	err := rs.FlowService.Publish(workspaceId, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "published flow with id: "+idParam)
}

func (rs FlowsResource) Unpublish(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	err := rs.FlowService.UnPublish(workspaceId, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "unpublished flow with id: "+idParam)
}

func (rs FlowsResource) Archive(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	err := rs.FlowService.Archive(workspaceId, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "archived flow with id: "+idParam)
}

func (rs FlowsResource) GetPossibleDependsOnListForFlow(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	possibleDependsOnList, err := rs.FlowService.GetPossibleDependsOnListForFlow(workspaceId, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, possibleDependsOnList)
}

func (rs FlowsResource) ListArchive(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	flows, err := rs.FlowService.ListArchivedFlows(workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, flows)
}

func (rs FlowsResource) RestoreFlow(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	err := rs.FlowService.RestoreFlow(workspaceId, idParam)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "restored flow with id: "+idParam)
}
