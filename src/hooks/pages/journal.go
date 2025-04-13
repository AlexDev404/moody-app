// Example of a form
package pages

import (
	"baby-blog/database/models"
	"baby-blog/forms"
	"baby-blog/forms/validator"
	"baby-blog/types"
	"log"
	"net/http"
)

func Journal(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	if r.Method != http.MethodPost {
		return pageData
	} else {
		if err := r.ParseForm(); err != nil {
			log.Fatal(forms.FormHandlerErrorMessage, "error", err)
			http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
			pageData["Failure"] = "✗ Failed to submit journal entry. Please try again later."
			return pageData
		}
		return JournalForm(pageData, db, r, w)
	}
}

func JournalForm(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {

	validator := validator.NewValidator()
	formData, formErrors := forms.JournalForm(w, r, validator)

	for key, value := range formErrors {
		formData[key] = value
	}

	if formErrors == nil {
		journalData := &models.Journal{
			Title:   formData["title"].(string),
			Content: formData["content"].(string),
		}
		err := db.Journal.Insert(journalData)
		if err != nil {
			formData["Failure"] = "✗ Failed to submit journal entry. Please try again later."
		} else {
			formData["Message"] = "✓ Your journal entry has been submitted. Thank you!"
		}
	}

	// Merge formData into pageData
	for key, value := range formData {
		pageData[key] = value
	}
	return pageData
}
