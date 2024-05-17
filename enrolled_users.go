package main

import (
	"math"
	"milestone_core/server"
	"milestone_core/users"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type usersResource struct {
	usersService users.Service
}

func (rs usersResource) Routes() chi.Router {
	r := chi.NewRouter()
	r.Get("/", rs.List)

	return r
}

func (rs usersResource) List(w http.ResponseWriter, r *http.Request) {
	pageNumberStr := r.URL.Query().Get("page")
	rowsPerPageStr := r.URL.Query().Get("rows")

	pageNumber, _ := strconv.Atoi(pageNumberStr)
	rowsPerPage, _ := strconv.Atoi(rowsPerPageStr)

	workspaceId := server.GetWorkspaceIdFromContext(r.Context())

	enrolledUsers, err := rs.usersService.List(workspaceId)
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
		EnrolledUsers []*users.EnrolledUser `json:"enrolledUsers"`
		CurrentPage   int                   `json:"currentPage"`
		Rows          int                   `json:"rowsCount"`
		TotalPages    int                   `json:"totalPages"`
		TotalRows     int                   `json:"totalRows"`
	}{
		paginatedEnrolledUsers, pageNumber, len(paginatedEnrolledUsers), totalPagesNumber, len(enrolledUsers),
	}

	server.SendJson(w, response)
}
