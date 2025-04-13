package pages

import (
	"baby-blog/database/models"
	"baby-blog/forms"
	"baby-blog/forms/validator"
	"baby-blog/types"
	"log"
	"net/http"
)

func Home(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	// Get today's mood entry if it exists
	allMood, err := db.Moods.GetAll()
	if err != nil {
		pageData["Failure"] = "Failed to load mood data"
		return pageData
	}

	if allMood != nil {
		pageData["emotions_felt"] = len(allMood)
		pageData["has_mood_today"] = true
	} else {
		pageData["emotions_felt"] = 0
		pageData["has_mood_today"] = false
	}

	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			log.Fatal(forms.FormHandlerErrorMessage, "error", err)
			http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
			pageData["Failure"] = "✗ Failed to submit journal entry. Please try again later."
			return pageData
		}
		return HomeForm(pageData, db, r, w)
	}
	return pageData
}

func HomeForm(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	validator := validator.NewValidator()
	formData, formErrors := forms.HomeForm(w, r, validator)

	if formErrors != nil {
		pageData["Failure"] = formErrors["Failure"]
		pageData["Errors"] = formErrors["Errors"]
		return pageData
	}

	mood := formData["mood"].(string)
	newMood := &models.MoodEntry{
		MoodText: mood,
	}
	if err := db.Moods.Insert(newMood); err != nil {
		log.Fatal(forms.FormHandlerErrorMessage, "error", err)
		http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
		pageData["Failure"] = "✗ Failed to submit mood. Please try again later."
		return pageData
	}

	pageData["Message"] = "✓ Mood submitted successfully!"
	return pageData
}
