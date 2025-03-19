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
	case "/feedback":
		pageErrors := map[string]interface{}{
			"Errors": map[string]string{
				"fullname": "Bruh",
				"email":    "Lmao",
				"subject":  "Subject is required",
				"message":  "Bad message",
			},
			"Message": "Good message",
		}
		name := r.FormValue("fullname")
		email := r.FormValue("email")
		subject := r.FormValue("subject")
		message := r.FormValue("message")

		if name == "" || email == "" || subject == "" || message == "" {
			app.Logger.Warn("Invalid contact form submission", "name", name, "email", email)
			http.Error(w, "Missing required fields", http.StatusBadRequest)
			return
		}

		// Log the contact form submission
		app.Logger.Info("Contact form submission received",
			"name", name,
			"email", email,
			"subject", subject,
			"message", message,
			"message_length", len(message))

		// Here you would typically send an email or store the contact form
		// For now, just render the page
		app.Render(w, r, app.templates, pageErrors)

	default:
		app.Logger.Warn("Unknown POST route accessed", "path", r.URL.Path)
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
