package http

import (
	"net/http"
	"strings"
	"user-svc/application/query"
)

// StudentHandler handles incoming HTTP traffic.
type StudentHandler struct {
	getStudent *query.GetStudentHandler
}

func NewStudentHandler(getStudent *query.GetStudentHandler) *StudentHandler {
	return &StudentHandler{
		getStudent: getStudent,
	}
}

// ServeHTTP implements the standard Go http.Handler interface.
/*
	type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
	}
*/
func (h *StudentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Log the incoming request
	println("Infrastructure - StudentHandler Received request for /student")

	path := strings.TrimPrefix(r.URL.Path, "/students")
	path = strings.Trim(path, "/")
	parts := strings.Split(path, "/")

	switch {
	case path == "":
		h.handleStudents(w, r)
	case len(parts) == 1:
		h.handleStudent(w, r, parts[0])
	case len(parts) == 2 && parts[1] == "courses":
		h.handleStudentCourses(w, r, parts[0])
	default:
		http.NotFound(w, r)
	}
}

func (h *StudentHandler) handleStudents(w http.ResponseWriter, r *http.Request) {
	result, err := h.getStudent.Handle(
		r.Context(),
		query.GetStudentQuery{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(result))
}
