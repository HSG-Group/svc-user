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

func NewRouterWithMux(h Handlers, auth TokenValidator) http.Handler {
	mux := http.NewServeMux()

	// Per-route chains
	authed := func(hf http.HandlerFunc) http.Handler {
		return Chain(hf, Auth(auth))
	}
	// ownOnly := func(hf http.HandlerFunc) http.Handler {
	// 	return Chain(hf, Auth(auth), OwnProfileOnly()) // Auth runs first, then ownership
	// }

	// ---------- Public (no auth) ----------
	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("GET /hello", h.Hello.ServeHTTP)

	// ---------- Students ----------
	mux.Handle("GET /students", authed(h.Student.List))
	mux.Handle("GET /students/{id}", authed(h.Student.GetByID))
	mux.Handle("GET /students/{id}/courses", authed(h.Student.ListCourses))

	// ---------- Users (matches the README table) ----------
	// mux.Handle("GET /users/search",               authed(h.Search.Search))
	// mux.Handle("GET /users/{id}",                 authed(h.Profile.Get))
	// mux.Handle("PATCH /users/{id}",               ownOnly(h.Profile.Update))
	// mux.Handle("POST /users/{id}/avatar",         ownOnly(h.Avatar.Upload))
	// mux.Handle("PUT /users/{id}/presence",        ownOnly(h.Presence.Set))
	// mux.Handle("GET /users/{id}/friends/pending", ownOnly(h.Friends.ListPending))
	// mux.Handle("POST /users/{id}/blocks",         ownOnly(h.Profile.Block))
	// mux.Handle("DELETE /users/{id}/blocks/{targetId}", ownOnly(h.Profile.Unblock))

	// Global middleware: wraps EVERYTHING (including /health and 404s)
	return Chain(mux,
		Recovery(), // outermost: catches panics from everything below
		Logging(),
		CORS("*"),
	)
}
