package query

import (
	"context"
	"fmt"
	"time"

	"user-svc/application/ports"
	"user-svc/domain/hello"
)

// GetHelloQuery defines the input payload for the use-case.
type GetHelloQuery struct{}

// GetHelloHandler orchestrates the use-case.
type GetHelloHandler struct {
	helloRepo ports.HelloRepo
}

// NewGetHelloHandler acts as the constructor.
func NewGetHelloHandler(repoParam ports.HelloRepo) *GetHelloHandler {
	return &GetHelloHandler{helloRepo: repoParam}
}

// Handle executes the Application logic.
// It fetches data from Infra via Ports, and passes it to the Domain layer.
/**
* @param h *GetHelloHandler - The handler instance. (receiver)
* @param ctx context.Context - The context for the request.
* @param query GetHelloQuery - The query object containing the input parameters.
* @return string - The resulting message text.
* @return error - An error if any occurred during processing.
 */
func (h *GetHelloHandler) Handle(ctx context.Context, query GetHelloQuery) (string, error) {
	fmt.Println("Application - Query - Handling GetHelloQuery...")

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

	name, err := h.helloRepo.GetStaticName(ctx)
	if err != nil {
		return "", err
	}

	msg, err := hello.NewMessage(name)
	if err != nil {
		return "", err
	}

	return msg.Text, nil
}
