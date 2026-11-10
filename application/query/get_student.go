package query

import (
	"context"
	"fmt"

	"user-svc/application/ports"
	"user-svc/domain/student"
)

// ---------------------------------------------------------------------------
// List all students
// ---------------------------------------------------------------------------

// ListStudentsQuery defines the input payload for listing students.
type ListStudentsQuery struct{}

// ListStudentsHandler orchestrates the list-students use-case.
type ListStudentsHandler struct {
	studentRepo ports.StudentRepo
}

func NewListStudentsHandler(repo ports.StudentRepo) *ListStudentsHandler {
	return &ListStudentsHandler{studentRepo: repo}
}

// Handle returns every student from the repository.
func (h *ListStudentsHandler) Handle(ctx context.Context, _ ListStudentsQuery) ([]student.StudentRecord, error) {
	fmt.Println("Application - Query - Handling ListStudentsQuery...")

	students, err := h.studentRepo.GetStudents(ctx)
	if err != nil {
		return nil, fmt.Errorf("list students: %w", err)
	}

	return students, nil
}

// ---------------------------------------------------------------------------
// Get a single student by ID
// ---------------------------------------------------------------------------

// GetStudentByIDQuery defines the input payload.
type GetStudentByIDQuery struct {
	ID string
}

// GetStudentByIDHandler orchestrates the get-student-by-id use-case.
type GetStudentByIDHandler struct {
	studentRepo ports.StudentRepo
}

func NewGetStudentByIDHandler(repo ports.StudentRepo) *GetStudentByIDHandler {
	return &GetStudentByIDHandler{studentRepo: repo}
}

// Handle returns a single student identified by ID.
func (h *GetStudentByIDHandler) Handle(ctx context.Context, q GetStudentByIDQuery) (student.StudentRecord, error) {
	fmt.Printf("Application - Query - Handling GetStudentByIDQuery (id=%s)...\n", q.ID)

	record, err := h.studentRepo.GetStudentByID(ctx, q.ID)
	if err != nil {
		return student.StudentRecord{}, fmt.Errorf("get student %s: %w", q.ID, err)
	}

	return record, nil
}

// ---------------------------------------------------------------------------
// List courses for a student
// ---------------------------------------------------------------------------

// ListStudentCoursesQuery defines the input payload.
type ListStudentCoursesQuery struct {
	StudentID string
}

// ListStudentCoursesHandler orchestrates the list-courses use-case.
type ListStudentCoursesHandler struct {
	studentRepo ports.StudentRepo
}

func NewListStudentCoursesHandler(repo ports.StudentRepo) *ListStudentCoursesHandler {
	return &ListStudentCoursesHandler{studentRepo: repo}
}

// Handle returns every course for the given student.
func (h *ListStudentCoursesHandler) Handle(ctx context.Context, q ListStudentCoursesQuery) ([]student.Course, error) {
	fmt.Printf("Application - Query - Handling ListStudentCoursesQuery (studentId=%s)...\n", q.StudentID)

	courses, err := h.studentRepo.GetStudentCourses(ctx, q.StudentID)
	if err != nil {
		return nil, fmt.Errorf("list courses for student %s: %w", q.StudentID, err)
	}

	return courses, nil
}
