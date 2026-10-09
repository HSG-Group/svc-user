package persistence

import (
	"context"
)

// MemoryHelloRepo is an adapter that implements ports.HelloRepo.
// In reality, this might connect to a PostgreSQL database via pgx.
type MemoryHelloRepo struct{}

func NewMemoryHelloRepo() *MemoryHelloRepo {
	//show meessage that request go here
	println("Infrastructure - HelloRepo request received")
	return &MemoryHelloRepo{}
}

// GetStaticName is the actual implementation of the interface.
func (m *MemoryHelloRepo) GetStaticName(ctx context.Context) (string, error) {
	// Hardcoded response simulating a database lookup
	println("Infrastructure - HelloRepo Fetching static name...")
	return "World", nil
}
