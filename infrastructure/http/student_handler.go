package http

import (
	"encoding/json"
	"net/http"

	"user-svc/application/query"
)

// StudentHandler handles incoming HTTP traffic for /students routes.
type StudentHandler struct {
	getStudent *query.GetStudentHandler
}

func NewStudentHandler(getStudent *query.GetStudentHandler) *StudentHandler {
	return &StudentHandler{
		getStudent: getStudent,
	}
}

// List handles GET /students — returns all students.
func (h *StudentHandler) List(w http.ResponseWriter, r *http.Request) {
	println("Infrastructure - StudentHandler - List")

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

// GetByID handles GET /students/{id} — returns a single student.
func (h *StudentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	println("Infrastructure - StudentHandler - GetByID:", id)

	if id == "" {
		http.Error(w, "missing student id", http.StatusBadRequest)
		return
	}

	// TODO: wire up a dedicated GetStudentByID use-case.
	// For now, return a simple JSON placeholder using the id.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"id":      id,
		"message": "student " + id,
	})
}

// ListCourses handles GET /students/{id}/courses.
func (h *StudentHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	println("Infrastructure - StudentHandler - ListCourses for student:", id)

	if id == "" {
		http.Error(w, "missing student id", http.StatusBadRequest)
		return
	}

	// TODO: wire up a ListCoursesByStudent use-case.
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"student_id": id,
		"message":    "courses for student " + id,
	})
}
