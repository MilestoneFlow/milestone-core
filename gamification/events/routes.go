package events

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"milestone_core/shared/rest"
	"milestone_core/shared/server"
	"net/http"
)

type Resource struct {
	EventsService Service
}

func (rs Resource) Routes() chi.Router {
	r := chi.NewRouter()

	r.Get("/", rs.List)
	r.Post("/", rs.Create)
	r.Route("/{id}", func(r chi.Router) {
		r.Get("/", rs.Get)
		r.Put("/", rs.Update)
		r.Delete("/", rs.Delete)
	})

	return r
}

func (rs Resource) PublicRoutes() chi.Router {
	r := chi.NewRouter()

	r.Post("/{key}/track", rs.TrackEvent)

	return r
}

func (rs Resource) List(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	events, err := rs.EventsService.GetEvents(workspaceId)
	if err != nil {
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendResponse(w, events, http.StatusOK)
}

func (rs Resource) Create(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		rest.SendErrorResponse(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	err = rs.EventsService.CreateEvent(workspaceId, event)
	if err != nil {
		if errors.Is(err, Errors.KeyExistsError) {
			rest.SendErrorResponse(w, errors.New("event key already exists"), http.StatusBadRequest)
			return
		}

		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendMessageResponse(w, "event created", http.StatusCreated)
}

func (rs Resource) Get(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	eventId := chi.URLParam(r, "id")
	event, err := rs.EventsService.GetEventById(workspaceId, eventId)
	if err != nil {
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendResponse(w, event, http.StatusOK)
}

func (rs Resource) Update(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	eventId := chi.URLParam(r, "id")
	var event Event
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		rest.SendErrorResponse(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	err = rs.EventsService.UpdateEvent(workspaceId, eventId, &event)
	if err != nil {
		if errors.Is(err, Errors.KeyExistsError) || errors.Is(err, Errors.InvalidEventError) {
			rest.SendErrorResponse(w, err, http.StatusBadRequest)
			return
		}

		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendMessageResponse(w, "event updated", http.StatusOK)
}

func (rs Resource) Delete(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	eventId := chi.URLParam(r, "id")
	err := rs.EventsService.DeleteEvent(workspaceId, eventId)
	if err != nil {
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendMessageResponse(w, "event deleted", http.StatusOK)
}

func (rs Resource) TrackEvent(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	key := chi.URLParam(r, "key")
	var request TrackEventRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		rest.SendErrorResponse(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	err = rs.EventsService.CreateEventUser(workspaceId, key, request.UserId, request.Metadata)
	if err != nil {
		if errors.Is(err, Errors.EventNotFoundByKeyError) {
			rest.SendErrorResponse(w, errors.New("event not found by key"), http.StatusNotFound)
			return
		}
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendMessageResponse(w, "event tracked", http.StatusOK)
}

type TrackEventRequest struct {
	UserId   string           `json:"userId"`
	Metadata *json.RawMessage `json:"metadata"`
}
