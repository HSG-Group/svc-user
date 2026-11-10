package http

import (
	"fmt"
	"log"
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
	return token, nil
}
