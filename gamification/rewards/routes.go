package rewards

import (
	"encoding/json"
	"errors"
	"github.com/go-chi/chi/v5"
	"milestone_core/shared/rest"
	"milestone_core/shared/server"
	"net/http"
	"reflect"
)

type Resource struct {
	Service Service
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

func (rs Resource) List(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	rewards, err := rs.Service.GetRewards(workspaceId)
	if err != nil {
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendResponse(w, rewards, http.StatusOK)
}

func (rs Resource) Create(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	var reward Reward
	err := json.NewDecoder(r.Body).Decode(&reward)
	if err != nil {
		rest.SendErrorResponse(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	err = rs.Service.CreateReward(workspaceId, reward)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&RewardError{}) {
			rest.SendErrorResponse(w, err, err.(*RewardError).HttpCode)
			return
		}
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendMessageResponse(w, "reward created", http.StatusCreated)
}

func (rs Resource) Get(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	rewardId := chi.URLParam(r, "id")
	reward, err := rs.Service.GetRewardById(workspaceId, rewardId)
	if err != nil {
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}
	if reward == nil {
		rest.SendErrorResponse(w, errors.New("reward not found"), http.StatusNotFound)
		return
	}

	rest.SendResponse(w, reward, http.StatusOK)
}

func (rs Resource) Update(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	rewardId := chi.URLParam(r, "id")
	var reward Reward
	err := json.NewDecoder(r.Body).Decode(&reward)
	if err != nil {
		rest.SendErrorResponse(w, errors.New("invalid request body"), http.StatusBadRequest)
		return
	}

	err = rs.Service.UpdateReward(workspaceId, rewardId, &reward)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&RewardError{}) {
			rest.SendErrorResponse(w, err, err.(*RewardError).HttpCode)
			return
		}
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendMessageResponse(w, "reward updated", http.StatusOK)
}

func (rs Resource) Delete(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	rewardId := chi.URLParam(r, "id")
	err := rs.Service.DeleteReward(workspaceId, rewardId)
	if err != nil {
		if reflect.TypeOf(err) == reflect.TypeOf(&RewardError{}) {
			rest.SendErrorResponse(w, err, err.(*RewardError).HttpCode)
			return
		}
		rest.SendErrorResponse(w, err, http.StatusInternalServerError)
		return
	}

	rest.SendMessageResponse(w, "reward deleted", http.StatusOK)
}
