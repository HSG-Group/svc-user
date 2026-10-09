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
	repo ports.HelloRepo
}

// NewGetHelloHandler acts as the constructor.
func NewGetHelloHandler(repo ports.HelloRepo) *GetHelloHandler {
	return &GetHelloHandler{repo: repo}
}

// Handle executes the Application logic.
// It fetches data from Infra via Ports, and passes it to the Domain layer.
func (h *GetHelloHandler) Handle(ctx context.Context, q GetHelloQuery) (string, error) {
	// Log the incoming request
	println("Application - Query - Handling GetHelloQuery...")

	// Call to Infrastructure (via port interface)
	name, err := h.repo.GetStaticName(ctx)
	if err != nil {
		return "", err
	}

	// Let Domain layer apply business rules/formatting
	msg := hello.NewMessage(name)
	return string(msg), nil
}
