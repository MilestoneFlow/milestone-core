package main

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"milestone_core/publicapi"
	"milestone_core/server"
	"milestone_core/tracker"
	"milestone_core/users"
	"net/http"
)

type PublicApiResource struct {
	Service          publicapi.Service
	Tracker          tracker.Tracker
	UserStateService publicapi.UserStateService
}

func (rs PublicApiResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/validate", rs.ValidateToken)
	r.Get("/helpers", rs.GetHelpers)
	r.Get("/flows/{id}", rs.Get)

	// Enroll user
	r.Post("/enroll", rs.Enroll)

	// User state
	r.Get("/{externalUserId}/state", rs.GetUserState)
	r.Post("/{externalUserId}/flows/enroll", rs.EnrollInFlow)
	r.Post("/{externalUserId}/flows/state", rs.UpdateFlowState)

	// Tracker
	r.Post("/track", rs.Track)

	return r
}

func (rs PublicApiResource) ValidateToken(w http.ResponseWriter, r *http.Request) {
	token := server.GetTokenFromPublicApiClientContext(r.Context())
	err := rs.Service.ValidateToken(token)
	if err != nil {
		server.SendBadRequestErrorJson(w, errors.New("invalid token"))
		return
	}

	server.SendJson(w, true)
}

func (rs PublicApiResource) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	token := server.GetTokenFromPublicApiClientContext(r.Context())
	resFlow, err := rs.Service.GetFlow(token, id)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, resFlow)
}

func (rs PublicApiResource) Enroll(w http.ResponseWriter, r *http.Request) {
	token := server.GetTokenFromPublicApiClientContext(r.Context())
	var newUser users.EnrolledUser
	err := json.NewDecoder(r.Body).Decode(&newUser)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.Service.EnrollUser(token, newUser)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "User enrolled successfully")
}

func (rs PublicApiResource) GetUserState(w http.ResponseWriter, r *http.Request) {
	externalUserId := chi.URLParam(r, "externalUserId")
	token := server.GetTokenFromPublicApiClientContext(r.Context())
	userState, err := rs.UserStateService.GetState(token, externalUserId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, userState)
}

func (rs PublicApiResource) Track(w http.ResponseWriter, r *http.Request) {
	var body publicapi.TrackEventsRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	workspaceId := server.GetWorkspaceIdFromPublicApiClientContext(r.Context())
	err = rs.Tracker.TrackEvents(workspaceId, body.ExternalUserID, body.Events)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "ok")
}

func (rs PublicApiResource) GetHelpers(w http.ResponseWriter, r *http.Request) {
	token := server.GetTokenFromPublicApiClientContext(r.Context())
	helpers, err := rs.Service.GetHelpers(token)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, helpers)
}

func (rs PublicApiResource) EnrollInFlow(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromPublicApiClientContext(r.Context())
	externalUserId := chi.URLParam(r, "externalUserId")
	flow, err := rs.Service.EnrollInFlow(workspaceId, externalUserId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, flow)
}

func (rs PublicApiResource) UpdateFlowState(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromPublicApiClientContext(r.Context())
	externalUserId := chi.URLParam(r, "externalUserId")
	var body publicapi.FlowStateUpdateRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.Service.UpdateFlowState(workspaceId, externalUserId, body)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "ok")
}
