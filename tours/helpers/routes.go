package helpers

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"milestone_core/shared/server"
	"net/http"
)

type Resource struct {
	Service Service
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Post("/", rs.Create)
	r.Put("/", rs.UpdateMulti)
	r.Route("/{publicId}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
		r.Delete("/", rs.Delete)
		r.Post("/publish", rs.Publish)
		r.Post("/unpublish", rs.Unpublish)
	})

	return r
}

func (rs Resource) List(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	helpers, err := rs.Service.List(workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, helpers)
}

func (rs Resource) Get(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	publicId := chi.URLParam(r, "publicId")

	helper, err := rs.Service.Get(publicId, workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, helper)
}

func (rs Resource) Create(w http.ResponseWriter, r *http.Request) {
	var helper Helper
	err := json.NewDecoder(r.Body).Decode(&helper)

	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	createdHelper, err := rs.Service.Create(workspaceId, helper)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, createdHelper)
}

func (rs Resource) UpdateMulti(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	var helpers []map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&helpers)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}
	for _, helper := range helpers {
		err = rs.Service.Update(helper["publicId"].(string), workspaceId, helper)
		if err != nil {
			server.SendBadRequestErrorJson(w, err)
			return
		}
	}

	server.SendJson(w, map[string]string{"message": "updated helpers"})
}

func (rs Resource) Update(w http.ResponseWriter, r *http.Request) {
	publicId := chi.URLParam(r, "publicId")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	var helper map[string]interface{}
	err := json.NewDecoder(r.Body).Decode(&helper)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	err = rs.Service.Update(publicId, workspaceId, helper)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, map[string]string{"message": "updated helper with publicId: " + publicId})
}

func (rs Resource) Delete(w http.ResponseWriter, r *http.Request) {
	publicId := chi.URLParam(r, "publicId")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	err := rs.Service.Delete(publicId, workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "deleted helper with publicId: "+publicId)
}

func (rs Resource) Publish(w http.ResponseWriter, r *http.Request) {
	publicId := chi.URLParam(r, "publicId")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	err := rs.Service.Publish(publicId, workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "published helper with publicId: "+publicId)
}

func (rs Resource) Unpublish(w http.ResponseWriter, r *http.Request) {
	publicId := chi.URLParam(r, "publicId")
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	err := rs.Service.Unpublish(publicId, workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, "unpublished helper with publicId: "+publicId)
}
