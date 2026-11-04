package http

import (
	"net/http"

	"user-svc/application/query"
)

// StudentHandler handles incoming HTTP traffic.
type StudentHandler struct {
	getStudent *query.GetStudentHandler
}

func NewStudentHandler(getStudent *query.GetStudentHandler) *StudentHandler {
	return &StudentHandler{getStudent: getStudent}
}

// ServeHTTP implements the standard Go http.Handler interface.
/*
	type Handler interface {
		ServeHTTP(ResponseWriter, *Request)
	}
*/
func (h *StudentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// Log the incoming request
	println("Infrastructure - StudentHandler Received request for /student")

	//handle the request methods
	switch r.Method {
	case http.MethodGet:
		// Execute the application use-case
		msg, err := h.getStudent.Handle(r.Context(), query.GetStudentQuery{})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg + "\n"))
	case http.MethodPost:
		// Handle POST request
		http.Error(w, "POST method not implemented", http.StatusNotImplemented)

	case http.MethodPut:
		// Handle PUT request
		http.Error(w, "PUT method not implemented", http.StatusNotImplemented)

	case http.MethodDelete:
		// Handle DELETE request
		http.Error(w, "DELETE method not implemented", http.StatusNotImplemented)
	default:
		// Handle other methods
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
