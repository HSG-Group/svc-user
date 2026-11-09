// infrastructure/http/router.go
package http

import (
	"context"
	"net/http"
	"strings"
)

type route struct {
	method   string
	segments []string
	handler  http.HandlerFunc
}

type Router struct {
	routes []route
}

type paramsKey struct{}

func split(p string) []string {
	return strings.Split(strings.Trim(p, "/"), "/")
}

func (rt *Router) Handle(method, pattern string, h http.HandlerFunc) {
	rt.routes = append(rt.routes, route{method, split(pattern), h})
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	parts := split(r.URL.Path)
	pathMatched := false

	for _, rte := range rt.routes {
		params, ok := match(rte.segments, parts)
		if !ok {
			continue
		}
		pathMatched = true
		if rte.method != r.Method {
			continue
		}
		ctx := context.WithValue(r.Context(), paramsKey{}, params)
		rte.handler(w, r.WithContext(ctx))
		return
	}

	if pathMatched {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.NotFound(w, r)
}

// match compares pattern segments like ["students", ":id"] with request segments.
func match(pattern, parts []string) (map[string]string, bool) {
	if len(pattern) != len(parts) {
		return nil, false
	}
	params := map[string]string{}
	for i, seg := range pattern {
		if strings.HasPrefix(seg, ":") {
			params[seg[1:]] = parts[i]
			continue
		}
		if seg != parts[i] {
			return nil, false
		}
	}
	return params, true
}

// Param reads a path parameter, like r.PathValue in the standard mux.
func Param(r *http.Request, name string) string {
	p, _ := r.Context().Value(paramsKey{}).(map[string]string)
	return p[name]
}

func NewRouter() *Router {
	return &Router{}
}
