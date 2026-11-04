package persistence

import (
	"context"
	"fmt"
	"user-svc/domain/student"
)

// MemoryStudentRepo is an adapter that implements ports.StudentRepo.
// In reality, this might connect to a PostgreSQL database via pgx.
type MemoryStudentRepo struct {
	students []student.StudentRecord
}

func NewMemoryStudentRepo() *MemoryStudentRepo {
	//show meessage that request go here
	println("Infrastructure - StudentRepo	 request received")
	return &MemoryStudentRepo{
		students: []student.StudentRecord{
			{ID: "1", Name: "Alice", Address: "123 Main St", Phone: "555-1234", Status: "active", School: "School A"},
			{ID: "2", Name: "Bob", Address: "456 Elm St", Phone: "555-5678", Status: "inactive", School: "School B"},
			{ID: "3", Name: "Charlie", Address: "789 Oak St", Phone: "555-9012", Status: "active", School: "School C"},
			{ID: "4", Name: "David", Address: "321 Pine St", Phone: "555-3456", Status: "active", School: "School D"},
			{ID: "5", Name: "Eve", Address: "654 Maple St", Phone: "555-7890", Status: "inactive", School: "School E"},
			{ID: "6", Name: "Frank", Address: "987 Cedar St", Phone: "555-2345", Status: "active", School: "School F"},
			{ID: "7", Name: "Grace", Address: "135 Birch St", Phone: "555-6789", Status: "active", School: "School G"},
			{ID: "8", Name: "Hannah", Address: "246 Spruce St", Phone: "555-0123", Status: "inactive", School: "School H"},
			{ID: "9", Name: "Ian", Address: "357 Willow St", Phone: "555-4567", Status: "active", School: "School I"},
			{ID: "10", Name: "Jack", Address: "468 Aspen St", Phone: "555-8901", Status: "active", School: "School J"},
		},
	}
}

// GetStudentName is the actual implementation of the interface.
func (m *MemoryStudentRepo) GetStudentName(ctx context.Context) (string, error) {
	// Hardcoded response simulating a database lookup
	println("Infrastructure - StudentRepo Fetching student name...")
	return "World", nil
}

func (r *MemoryStudentRepo) GetStudentByID(
	ctx context.Context,
	id string,
) (student.StudentRecord, error) {
	for _, student := range r.students {
		if student.ID == id {
			return student, nil
		}
	}

	return student.StudentRecord{}, fmt.Errorf("student %q not found", id)
}

func (r *MemoryStudentRepo) GetStudents(ctx context.Context) ([]student.StudentRecord, error) {
	// Simulate a database query by returning the in-memory slice of students
	return r.students, nil
}
