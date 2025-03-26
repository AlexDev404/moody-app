package forms

import (
	"baby-blog/forms/validator"
	"net/http"
)

func TodoForm(w http.ResponseWriter, r *http.Request, v *validator.Validator) (map[string]interface{}, map[string]interface{}) {
	var formErrors map[string]interface{}
	formData := map[string]interface{}{
		"task": r.FormValue("task"),
	}

	// Validate form data
	errors := validateTodoFormData(formData, v)

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

func validateTodoFormData(formData map[string]interface{}, v *validator.Validator) map[string]string {
	errors := v.Errors

	// Check if the task is valid
	v.Check(validator.NotBlank(formData["task"].(string)), "task", "Task is required")
	v.Check(validator.MaxLength(formData["task"].(string), 200), "task", "Task must be less than 200 characters")

	return errors
}
