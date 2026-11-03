package query

import (
	"context"

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
	// Log the incoming request
	println("Application - Query - Handling GetHelloQuery...")

	// Call to Infrastructure (via port interface)
	name, err := h.helloRepo.GetStaticName(ctx)
	if err != nil {
		return "", err
	}

	// Let Domain layer apply business rules/formatting (may return validation error)
	msg, err := hello.NewMessage(name)
	if err != nil {
		return "", err
	}

	return msg.Text, nil
}
