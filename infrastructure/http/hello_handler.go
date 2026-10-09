package http

import (
	"net/http"

	"user-svc/application/query"
)

// HelloHandler handles incoming HTTP traffic.
type HelloHandler struct {
	getHello *query.GetHelloHandler
}

func NewHelloHandler(getHello *query.GetHelloHandler) *HelloHandler {
	return &HelloHandler{getHello: getHello}
}

// ServeHTTP implements the standard Go http.Handler interface.
func (h *HelloHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Log the incoming request
	println("Infrastructure - HelloHandler Received request for /hello")

	// Execute the application use-case
	msg, err := h.getHello.Handle(r.Context(), query.GetHelloQuery{})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(msg + "\n"))
}
