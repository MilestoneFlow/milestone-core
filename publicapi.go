package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"milestone_core/publicapi"
	"milestone_core/server"
	"milestone_core/users"
	"net/http"
)

type PublicApiResource struct {
	Service          publicapi.Service
	Tracker          publicapi.Tracker
	UserStateService publicapi.UserStateService
}

func (rs PublicApiResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/helpers", rs.GetHelpers)
	r.Get("/flows/{id}", rs.Get)

	r.Post("/enroll", rs.Enroll)
	r.Get("/{externalUserId}/state", rs.GetUserState)
	r.Post("/{externalUserId}/state", rs.UpdateUserState)
	r.Post("/track", rs.Track)

	return r
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

func (rs PublicApiResource) UpdateUserState(w http.ResponseWriter, r *http.Request) {
	externalUserId := chi.URLParam(r, "externalUserId")
	token := server.GetTokenFromPublicApiClientContext(r.Context())
	var userState users.UserState
	err := json.NewDecoder(r.Body).Decode(&userState)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.UserStateService.PutState(token, externalUserId, userState)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "User state updated successfully")
}

func (rs PublicApiResource) Track(w http.ResponseWriter, r *http.Request) {
	var body TrackEventsRequest
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

type TrackEventsRequest struct {
	Events         []publicapi.EventTrack `json:"data"`
	ExternalUserID string                 `json:"externalUserId"`
}
