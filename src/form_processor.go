package main

import (
	"baby-blog/forms"
	"net/http"
)

// POSTProcessor handles all POST requests to the application
func (app *Application) POSTHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request has a valid form
	if err := r.ParseForm(); err != nil {
		app.Logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	// Place all form submission routes here
	switch r.URL.Path {
	case "/feedback":
		forms.FeedbackForm(w, r)
		app.Render(w, r, app.templates, nil)
	default:
		app.Logger.Warn("Unknown POST route accessed", "path", r.URL.Path)
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
