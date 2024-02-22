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

// Routes creates a REST router for the todos resource
func (rs usersResource) Routes() chi.Router {
	r := chi.NewRouter()
	// r.Use() // some middleware..

	r.Get("/", rs.List) // GET /users - read a list of users

	return r
}

func (rs usersResource) List(w http.ResponseWriter, r *http.Request) {
	pageNumberStr := r.URL.Query().Get("page")
	rowsPerPageStr := r.URL.Query().Get("rows")

	pageNumber, _ := strconv.Atoi(pageNumberStr)
	rowsPerPage, _ := strconv.Atoi(rowsPerPageStr)

	enrolledUsers, err := rs.usersService.List()
	if err != nil {
		server.SendBadRequestErrorJson(w, err)
		return
	}

	totalPagesNumber := int(math.Ceil(float64(float32(len(enrolledUsers)) / float32(rowsPerPage))))

	usersOffset := rowsPerPage * (pageNumber - 1)
	usersMax := min(usersOffset+rowsPerPage, len(enrolledUsers))
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
