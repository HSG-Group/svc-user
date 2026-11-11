package http

import (
	"fmt"
	"log"
)

const (
	// FakeToken is a placeholder token for local development.
	FakeToken = "YOUR_TOKEN_HERE"
)

// FakeValidator is for local development only.
// It treats the bearer token as the user ID.
type FakeValidator struct{}

func (FakeValidator) Validate(token string) (string, error) {
	log.Println("FakeValidator: Validating token:", token)
	println("FakeValidator: Validating token:", token)
	if token == "" || token == "invalid" {
		return "", fmt.Errorf("invalid token")
	}

	// Simulate a validation error for any token that is not the expected fake token.
	if token != FakeToken {
		return "", fmt.Errorf("Your token is not the expected fake token. Please use the correct token for local development.")
	}

	return token, nil
}
