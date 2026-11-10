package http

import (
	"context"
	"log"
	"net/http"
	"runtime/debug"
	"strings"
	"time"
)

// Middleware wraps a handler with extra behaviour.
type Middleware func(http.Handler) http.Handler

// Chain applies middlewares so the FIRST one listed is the OUTERMOST.
// Chain(h, A, B, C)  =>  A(B(C(h)))  => request flows A -> B -> C -> h
func Chain(h http.Handler, mws ...Middleware) http.Handler {
	for i := len(mws) - 1; i >= 0; i-- {
		h = mws[i](h)
	}
	return h
}

// ---------- Recovery ----------
func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Printf("panic: %v\n%s", rec, debug.Stack())
					http.Error(w, "internal server error", http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// ---------- Logging ----------
type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (s *statusRecorder) WriteHeader(code int) {
	s.status = code
	s.ResponseWriter.WriteHeader(code)
}

func Logging() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rec, r)
			log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start))
		})
	}
}

// ---------- CORS ----------
func CORS(allowedOrigin string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
			w.Header().Set("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Authorization,Content-Type")
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

// ---------- Auth (JWT) ----------
type ctxKey string

const UserIDKey ctxKey = "userID"

// TokenValidator is a port-like interface so middleware doesn't depend on a concrete JWT lib.
type TokenValidator interface {
	Validate(token string) (userID string, err error)
}

func Auth(v TokenValidator) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				http.Error(w, "missing bearer token", http.StatusUnauthorized)
				return
			}
			userID, err := v.Validate(token)
			if err != nil {
				http.Error(w, "invalid token", http.StatusUnauthorized)
				return
			}
			ctx := context.WithValue(r.Context(), UserIDKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// ---------- Ownership check: {id} in path must equal caller ----------
func OwnProfileOnly() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			caller, _ := r.Context().Value(UserIDKey).(string)
			if caller == "" || caller != r.PathValue("id") {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
