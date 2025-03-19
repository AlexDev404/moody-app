package forms

import (
	"log/slog"
	"net/http"
	"os"
)

// Create a new logger instance
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func FeedbackForm(w http.ResponseWriter, r *http.Request) map[string]interface{} {
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
		logger.Warn("Invalid contact form submission", "name", name, "email", email)
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return nil
	}

	// Log the contact form submission
	logger.Info("Contact form submission received",
		"name", name,
		"email", email,
		"subject", subject,
		"message", message,
		"message_length", len(message))

	// Here you would typically send an email or store the contact form
	// For now, just render the page
	return pageErrors
}
