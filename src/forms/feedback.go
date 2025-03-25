package forms

import (
	"baby-blog/forms/validator"
	"net/http"
)

func FeedbackForm(w http.ResponseWriter, r *http.Request, v *validator.Validator) (map[string]interface{}, map[string]interface{}) {
	var formErrors map[string]interface{}
	formData := map[string]interface{}{
		"fullname": r.FormValue("fullname"),
		"email":    r.FormValue("email"),
		"subject":  r.FormValue("subject"),
		"message":  r.FormValue("message"),
	}

	// Validate form data
	errors := validateFormData(formData, v)

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

func validateFormData(formData map[string]interface{}, v *validator.Validator) map[string]string {
	errors := v.Errors

	// Check if the fullname is valid
	v.Check(validator.NotBlank(formData["fullname"].(string)), "fullname", "Full name is required")

	// Check if the email is valid
	v.Check(validator.NotBlank(formData["email"].(string)), "email", "Email is required")
	v.Check(validator.IsValidEmail(formData["email"].(string)), "email", "Email is not valid")

	// Check if the subject is valid
	v.Check(validator.NotBlank(formData["subject"].(string)), "subject", "Subject is required")
	v.Check(validator.MinLength(formData["subject"].(string), 5), "subject", "Subject must be at least 5 characters")

	// Check if the message is valid
	v.Check(validator.NotBlank(formData["message"].(string)), "message", "Message is required")
	v.Check(validator.MaxLength(formData["message"].(string), 500), "message", "Message must be less than 500 characters")
	v.Check(validator.MinLength(formData["message"].(string), 10), "message", "Message must be at least 10 characters")

	return errors
}
