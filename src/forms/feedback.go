package forms

import (
	"baby-blog/forms/validator"
	"log/slog"
	"net/http"
	"os"
)

// Create a new logger instance
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

func FeedbackForm(w http.ResponseWriter, r *http.Request, v *validator.Validator) (map[string]interface{}, map[string]interface{}) {
	var formErrors map[string]interface{}
	formData := map[string]interface{}{
		"fullname": r.FormValue("fullname"),
		"email":    r.FormValue("email"),
		"subject":  r.FormValue("subject"),
		"message":  r.FormValue("message"),
	}

	// Validate form data
	errors := v.Errors
	v.Check(validator.NotBlank(formData["fullname"].(string)), "fullname", "Full name is required")
	v.Check(validator.NotBlank(formData["email"].(string)), "email", "Email is required")
	v.Check(validator.NotBlank(formData["subject"].(string)), "subject", "Subject is required")
	v.Check(validator.NotBlank(formData["message"].(string)), "message", "Message is required")

	// Check if any errors occurred
	if len(errors) > 0 {
		formErrors = map[string]interface{}{
			"Errors":  errors,
			"Failure": "✗ Please check your errors and try again.",
		}
		logger.Warn("Invalid form submission", "errors", errors)
		return formData, formErrors
	}

	// // Log the contact form submission
	// logger.Info("Feedback form submission received",
	// 	"name", formData["fullname"],
	// 	"email", formData["email"],
	// 	"subject", formData["subject"],
	// 	"message_length", len(formData["message"].(string)))

	// No errors
	return formData, nil
}
