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
	Service publicapi.Service
}

func (rs PublicApiResource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/{id}", rs.Get)
	r.Post("/enroll", rs.Enroll)
	r.Get("/{externalUserId}/state", rs.GetUserState)
	r.Post("/{externalUserId}/state", rs.UpdateUserState)
	r.Post("/{externalUserId}/enroll", rs.EnrollInNextFlow)

	return r
}

func (rs PublicApiResource) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	token := r.Context().Value("token").(string)
	resFlow, err := rs.Service.GetFlow(token, id)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, resFlow)
}

func (rs PublicApiResource) Enroll(w http.ResponseWriter, r *http.Request) {
	token := r.Context().Value("token").(string)
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
	token := r.Context().Value("token").(string)
	userState, err := rs.Service.GetUserState(token, externalUserId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, userState)
}

func (rs PublicApiResource) UpdateUserState(w http.ResponseWriter, r *http.Request) {
	externalUserId := chi.URLParam(r, "externalUserId")
	token := r.Context().Value("token").(string)
	var userState users.UserState
	err := json.NewDecoder(r.Body).Decode(&userState)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	userState.UserId = externalUserId
	err = rs.Service.UpdateUserState(token, userState)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "User state updated successfully")
}

func (rs PublicApiResource) EnrollInNextFlow(w http.ResponseWriter, r *http.Request) {
	externalUserId := chi.URLParam(r, "externalUserId")
	token := r.Context().Value("token").(string)
	userState, err := rs.Service.EnrollInNextFlow(token, externalUserId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, userState)
}
