package forms

import (
	"baby-blog/forms/validator"
	"net/http"
)

func MoodForm(w http.ResponseWriter, r *http.Request, v *validator.Validator) (map[string]interface{}, map[string]interface{}) {
	var formErrors map[string]interface{}
	formData := map[string]interface{}{
		"mood_id": r.FormValue("mood_id"),
	}

	// Validate form data
	errors := validateMoodFormData(formData, v)

	// Check if any errors occurred
	if len(errors) > 0 {
		formErrors = map[string]interface{}{
			"Errors":  errors,
			"Failure": "✗ Please check your errors and try again.",
		}
		logger.Warn("Invalid form submission", "errors", errors)
		return formData, formErrors
	}

	// No errors
	return formData, nil
}

func validateMoodFormData(formData map[string]interface{}, v *validator.Validator) map[string]string {
	errors := v.Errors

	// Check if the mood is valid
	v.Check(validator.NotBlank(formData["mood_id"].(string)), "mood", "Mood is required")

	return errors
}
