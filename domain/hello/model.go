package hello

// Message is a domain value object representing a custom greeting.
type Message string

// NewMessage contains our core domain formatting rules
func NewMessage(name string) Message {
	// In a real-world scenario, you might have more complex business logic here.
	println("Domain - model - Creating new greeting for: " + name)
	return Message("Hello, " + name + " from the Domain layer!")
}
