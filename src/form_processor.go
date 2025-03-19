package main

import (
	"baby-blog/database/models"
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
		formData, formErrors := forms.FeedbackForm(w, r)

		// Merge form errors with form data (essentially, append)
		for key, value := range formErrors {
			formData[key] = value
		}

		// Handle database insertion (this part needs proper implementation)
		if formErrors == nil {
			feedbackData := &models.Feedback{
				Fullname: formData["fullname"].(string),
				Email:    formData["email"].(string),
				Subject:  formData["subject"].(string),
				Message:  formData["message"].(string),
			}
			err := app.models.Feedback.Insert(feedbackData)
			if err != nil {
				formData["Failure"] = "✗ Failed to submit feedback. Please try again later."

			} else {
				formData["Message"] = "✓ Your feedback has been submitted. Thank you!"
			}
		}

		// Render the feedback page again with the form data and any errors
		app.Render(w, r, app.templates, formData)
	default:
		app.Logger.Warn("Unknown POST route accessed", "path", r.URL.Path)
		http.Error(w, "Not Found", http.StatusNotFound)
	}
}
