package main

import (
	"net/http"
	"strings"
)

func (app *Application) getPath(r *http.Request) string {
	path := r.URL.Path[1:]
	if path == "" {
		path = "index"
	} else {
		path = strings.TrimSuffix(path, "/")
	}
	return path
}

func (app *Application) isDisallowedRoute(path string) bool {
	disallowedRoutes := []string{"context", "head", "header", "footer", "current_ctx"}
	for _, route := range disallowedRoutes {
		if path == route {
			return true
		}
	}
	return strings.Contains(path, "layout/")
}
