package hello

import (
	"errors"
	"strings"
	"time"
)

// Message is a domain value object representing a custom greeting.
type Message struct {
	Text      string
	CreatedAt time.Time
}

// NewMessage validates the name and applies domain formatting rules.
// It uses only Go standard library — no external dependencies allowed in domain/.
func NewMessage(name string) (*Message, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("name cannot be empty")
	}

	return &Message{
		Text:      "Hello, " + name + " from the Domain layer!",
		CreatedAt: time.Now(),
	}, nil
}
