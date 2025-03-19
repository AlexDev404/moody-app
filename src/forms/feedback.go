package forms

import (
	"log/slog"
	"net/http"
	"os"
)

// Create a new logger instance
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func FeedbackForm(w http.ResponseWriter, r *http.Request) (map[string]interface{}, map[string]interface{}) {
	var formErrors map[string]interface{}
	formData := map[string]interface{}{
		"fullname": r.FormValue("fullname"),
		"email":    r.FormValue("email"),
		"subject":  r.FormValue("subject"),
		"message":  r.FormValue("message"),
	}

	// Validate form data
	errors := make(map[string]string)
	if formData["fullname"] == "" {
		errors["fullname"] = "Full name is required"
	}
	if formData["email"] == "" {
		errors["email"] = "Email is required"
	}
	if formData["subject"] == "" {
		errors["subject"] = "Subject is required"
	}
	if formData["message"] == "" {
		errors["message"] = "Message is required"
	}

	// Check if any errors occurred
	if len(errors) > 0 {
		formErrors = map[string]interface{}{
			"Errors":  errors,
			"Message": "Please fix the errors and try again",
		}
		logger.Warn("Invalid contact form submission", "errors", errors)
		return formData, formErrors
	}

	// Log the contact form submission
	logger.Info("Contact form submission received",
		"name", formData["fullname"],
		"email", formData["email"],
		"subject", formData["subject"],
		"message_length", len(formData["message"].(string)))

	// No errors
	return formData, nil
}
