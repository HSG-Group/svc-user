package ports

import (
	"context"
	"user-svc/domain/student"
)

type StudentRepo interface {
	GetStudentName(ctx context.Context) (string, error)
	GetStudentByID(ctx context.Context, id string) (student.StudentRecord, error)
	GetStudents(ctx context.Context) ([]student.StudentRecord, error)
}
