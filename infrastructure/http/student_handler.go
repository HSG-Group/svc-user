package http

import (
	"encoding/json"
	"net/http"

	"user-svc/application/query"
)

// StudentHandler handles incoming HTTP traffic for /students routes.
type StudentHandler struct {
	listStudents       *query.ListStudentsHandler
	getStudentByID     *query.GetStudentByIDHandler
	listStudentCourses *query.ListStudentCoursesHandler
}

func NewStudentHandler(
	listStudents *query.ListStudentsHandler,
	getStudentByID *query.GetStudentByIDHandler,
	listStudentCourses *query.ListStudentCoursesHandler,
) *StudentHandler {
	return &StudentHandler{
		listStudents:       listStudents,
		getStudentByID:     getStudentByID,
		listStudentCourses: listStudentCourses,
	}
}

// List handles GET /students — returns all students.
func (h *StudentHandler) List(w http.ResponseWriter, r *http.Request) {
	println("Infrastructure - StudentHandler - List")

	students, err := h.listStudents.Handle(r.Context(), query.ListStudentsQuery{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(students)
}

// GetByID handles GET /students/{id} — returns a single student.
func (h *StudentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	println("Infrastructure - StudentHandler - GetByID:", id)

	if id == "" {
		http.Error(w, "missing student id", http.StatusBadRequest)
		return
	}

	student, err := h.getStudentByID.Handle(r.Context(), query.GetStudentByIDQuery{ID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(student)
}

// ListCourses handles GET /students/{id}/courses.
func (h *StudentHandler) ListCourses(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	println("Infrastructure - StudentHandler - ListCourses for student:", id)

	if id == "" {
		http.Error(w, "missing student id", http.StatusBadRequest)
		return
	}

	courses, err := h.listStudentCourses.Handle(r.Context(), query.ListStudentCoursesQuery{StudentID: id})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(courses)
}
