// infrastructure/http/router_with_mux.go
package http

import (
	"net/http"
)

type Handlers struct {
	Hello   *HelloHandler
	Student *StudentHandler
	// Profile  *ProfileHandler
	// Avatar   *AvatarHandler
	// Presence *PresenceHandler
	// Friends  *FriendsHandler
	// Search   *SearchHandler
}

func NewRouterWithMux(h Handlers) http.Handler {
	mux := http.NewServeMux()

	// ---------- Public ----------
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /hello", h.Hello.ServeHTTP)

	// ---------- Students ----------
	mux.HandleFunc("GET /students", h.Student.List)
	mux.HandleFunc("GET /students/{id}", h.Student.GetByID)
	mux.HandleFunc("GET /students/{id}/courses", h.Student.ListCourses)

	// ---------- Users (enable as handlers are built) ----------
	// mux.HandleFunc("GET /users/search",                    h.Search.Search)
	// mux.HandleFunc("GET /users/{id}",                      h.Profile.Get)
	// mux.HandleFunc("PATCH /users/{id}",                    h.Profile.Update)
	// mux.HandleFunc("POST /users/{id}/avatar",              h.Avatar.Upload)
	//
	// mux.HandleFunc("GET /users/{id}/presence",             h.Presence.Get)
	// mux.HandleFunc("PUT /users/{id}/presence",             h.Presence.Set)
	//
	// mux.HandleFunc("GET /users/{id}/friends",              h.Friends.List)
	// mux.HandleFunc("GET /users/{id}/friends/pending",      h.Friends.ListPending)
	// mux.HandleFunc("POST /users/{id}/friends",             h.Friends.Send)
	// mux.HandleFunc("PATCH /users/{id}/friends/{friendId}", h.Friends.Accept)
	// mux.HandleFunc("DELETE /users/{id}/friends/{friendId}",h.Friends.Remove)
	//
	// mux.HandleFunc("POST /users/{id}/blocks",              h.Profile.Block)
	// mux.HandleFunc("DELETE /users/{id}/blocks/{targetId}", h.Profile.Unblock)

	return mux
}
