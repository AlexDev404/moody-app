package forms

import (
	"baby-blog/forms/validator"
	"net/http"
)

func JournalForm(w http.ResponseWriter, r *http.Request, v *validator.Validator) (map[string]interface{}, map[string]interface{}) {
	var formErrors map[string]interface{}
	formData := map[string]interface{}{
		"title":   r.FormValue("title"),
		"content": r.FormValue("content"),
	}

	// Validate form data
	errors := validateJournalFormData(formData, v)

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

func validateJournalFormData(formData map[string]interface{}, v *validator.Validator) map[string]string {
	errors := v.Errors

	// Check if the title is valid
	v.Check(validator.NotBlank(formData["title"].(string)), "title", "Title is required")
	v.Check(validator.MinLength(formData["title"].(string), 3), "title", "Title must be at least 3 characters")
	v.Check(validator.MaxLength(formData["title"].(string), 100), "title", "Title must be less than 100 characters")

	// Check if the content is valid
	v.Check(validator.NotBlank(formData["content"].(string)), "content", "Content is required")
	v.Check(validator.MaxLength(formData["content"].(string), 1000), "content", "Content must be less than 1000 characters")

	return errors
}
