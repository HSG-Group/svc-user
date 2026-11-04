package query

import (
	"context"
	"fmt"
	"time"

	"user-svc/application/ports"
	"user-svc/domain/hello"
)

// GetStudentQuery defines the input payload for the use-case.
type GetStudentQuery struct{}

// GetStudentHandler orchestrates the use-case.
type GetStudentHandler struct {
	studentRepo ports.StudentRepo
}

// NewGetStudentHandler acts as the constructor.
func NewGetStudentHandler(repoParam ports.StudentRepo) *GetStudentHandler {
	return &GetStudentHandler{studentRepo: repoParam}
}

// Handle executes the Application logic.
// It fetches data from Infra via Ports, and passes it to the Domain layer.
/**
* @param h *GetStudentHandler - The handler instance. (receiver)
* @param ctx context.Context - The context for the request.
* @param query GetStudentQuery - The query object containing the input parameters.
* @return string - The resulting message text.
* @return error - An error if any occurred during processing.
 */
func (h *GetStudentHandler) Handle(ctx context.Context, query GetStudentQuery) (string, error) {
	fmt.Println("Application - Query - Handling GetStudentQuery...")

	// 1. Whole context chain (type + nested WithValue/WithCancel layers)
	fmt.Printf("Application - Query - Handle - ctx: %v\n", ctx)
	fmt.Printf("Application - Query - Handle - ctx type: %T\n", ctx)

	// 2. Deadline / cancellation state
	if deadline, ok := ctx.Deadline(); ok {
		fmt.Printf("Application - Query - Handle - deadline: %v (in %v)\n", deadline, time.Until(deadline))
	} else {
		fmt.Println("Application - Query - Handle - deadline: none")
	}
	fmt.Printf("Application - Query - Handle - ctx.Err(): %v\n", ctx.Err())

	// 3. Specific values you put in the context
	// fmt.Printf("Application - Query - Handle - requestID: %v\n", ctx.Value(requestIDKey))

	name, err := h.studentRepo.GetStudentName(ctx)
	if err != nil {
		return "", err
	}

	msg, err := hello.NewMessage(name)
	if err != nil {
		return "", err
	}

	return msg.Text, nil
}
