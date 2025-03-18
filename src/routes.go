package main

import (
	"net/http"
)

func (app *Application) POSTProcessor(w http.ResponseWriter, r *http.Request) {
	// Place all form submission routes here
	switch r.URL.Path {
	case "/contact":
		// Handle contact form submission
		// ...
		break
	default:
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
