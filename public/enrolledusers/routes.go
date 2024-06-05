package enrolledusers

import (
	"math"
	"milestone_core/server"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type UsersResource struct {
	UsersService Service
}

func (rs UsersResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", rs.List)

	r.Route("/{id}", func(r chi.Router) {
		r.Delete("/", rs.Delete)
		r.Post("/reset", rs.ResetState)
	})

	return r
}

func (rs UsersResource) List(w http.ResponseWriter, r *http.Request) {
	pageNumberStr := r.URL.Query().Get("page")
	rowsPerPageStr := r.URL.Query().Get("rows")

	pageNumber, _ := strconv.Atoi(pageNumberStr)
	rowsPerPage, _ := strconv.Atoi(rowsPerPageStr)

	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	enrolledUsers, err := rs.UsersService.List(workspaceId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	totalPagesNumber := int(math.Ceil(float64(float32(len(enrolledUsers)) / float32(rowsPerPage))))

	usersOffset := min(rowsPerPage*(pageNumber-1), len(enrolledUsers))
	usersMax := min(usersOffset+rowsPerPage, len(enrolledUsers))
	if usersOffset < 0 {
		usersOffset = 0
		usersMax = 0
	}
	paginatedEnrolledUsers := enrolledUsers[usersOffset:usersMax]
	var response = struct {
		EnrolledUsers []*EnrolledUser `json:"enrolledUsers"`
		CurrentPage   int             `json:"currentPage"`
		Rows          int             `json:"rowsCount"`
		TotalPages    int             `json:"totalPages"`
		TotalRows     int             `json:"totalRows"`
	}{
		paginatedEnrolledUsers, pageNumber, len(paginatedEnrolledUsers), totalPagesNumber, len(enrolledUsers),
	}

	server.SendJson(w, response)
}

func (rs UsersResource) Delete(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	userId := chi.URLParam(r, "id")

	err := rs.UsersService.Delete(workspaceId, userId)
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, nil)

}

func (rs UsersResource) ResetState(w http.ResponseWriter, r *http.Request) {
	workspaceId := server.GetWorkspaceIdFromContext(r.Context())
	userId := chi.URLParam(r, "id")

	err := rs.UsersService.PutState(workspaceId, userId, UserState{
		FlowsData: FlowsData{
			CompletedFlowsIds:          nil,
			SkippedFlowsIds:            nil,
			CurrentFlowID:              "",
			LastSubmittedFlowID:        "",
			LastSubmittedFlowTimestamp: 0,
		},
	})
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	server.SendJson(w, nil)
}
