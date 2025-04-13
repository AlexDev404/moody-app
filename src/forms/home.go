package forms

import (
	"baby-blog/forms/validator"
	"net/http"
)

func HomeForm(w http.ResponseWriter, r *http.Request, v *validator.Validator) (map[string]interface{}, map[string]interface{}) {
	var formErrors map[string]interface{}
	formData := map[string]interface{}{
		"mood": r.FormValue("mood"),
	}

	// Validate form data
	errors := validateHomeFormData(formData, v)

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

func validateHomeFormData(formData map[string]interface{}, v *validator.Validator) map[string]string {
	errors := v.Errors

	// Check if the mood is valid
	v.Check(validator.NotBlank(formData["mood"].(string)), "mood", "Mood is required")
	v.Check(validator.MinLength(formData["mood"].(string), 3), "mood", "Mood must be at least 3 characters")
	v.Check(validator.MaxLength(formData["mood"].(string), 100), "mood", "Mood must not exceed 100 characters")

	return errors
}
