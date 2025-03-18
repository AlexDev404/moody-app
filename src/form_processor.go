package main

import (
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
	case "/contact":
		name := r.FormValue("name")
		email := r.FormValue("email")
		message := r.FormValue("message")

		if name == "" || email == "" || message == "" {
			app.Logger.Warn("Invalid contact form submission", "name", name, "email", email)
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Log the contact form submission
		app.Logger.Info("Contact form submission received",
			"name", name,
			"email", email,
			"message_length", len(message))

		// Here you would typically send an email or store the contact form
		// For now, just return a success message
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Thank you for your message!"))

	default:
		app.Logger.Warn("Unknown POST route accessed", "path", r.URL.Path)
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
