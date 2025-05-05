package pages

import (
	"baby-blog/database/models"
	"baby-blog/forms"
	"baby-blog/forms/validator"
	"baby-blog/middleware"
	"baby-blog/types"
	"log"
	"net/http"
)

func Home(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	// Get user ID from context
	userID := middleware.GetUserID(r)
	
	// If authenticated, show user-specific data
	if userID != "" {
		// Get user's mood entries
		allMood, err := db.Moods.GetAllForUser(userID)
		if err != nil {
			pageData["Failure"] = "Failed to load mood data"
			return pageData
		}

		if len(allMood) > 0 {
			pageData["emotions_felt"] = len(allMood)
			pageData["has_mood_today"] = true
		} else {
			pageData["emotions_felt"] = 0
			pageData["has_mood_today"] = false
		}
	} else {
		// Not authenticated, show default data
		pageData["emotions_felt"] = 0
		pageData["has_mood_today"] = false
	}

	if r.Method == http.MethodPost {
		// Check if user is authenticated before allowing mood submission
		if userID == "" {
			pageData["Failure"] = "You must be logged in to submit a mood."
			return pageData
		}
		
		if err := r.ParseForm(); err != nil {
			log.Println(forms.FormHandlerErrorMessage, "error", err)
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

	// Get user ID from context
	userID := middleware.GetUserID(r)
	if userID == "" {
		pageData["Failure"] = "You must be logged in to submit a mood."
		return pageData
	}

	mood := formData["mood"].(string)
	newMood := &models.MoodEntry{
		MoodText: mood,
		UserID:   userID, // Associate the mood with the current user
	}
	
	if err := db.Moods.Insert(newMood); err != nil {
		log.Println(forms.FormHandlerErrorMessage, "error", err)
		http.Error(w, forms.FormHandlerBadRequestMessage, http.StatusBadRequest)
		pageData["Failure"] = "✗ Failed to submit mood. Please try again later."
		return pageData
	}

	pageData["Message"] = "✓ Mood submitted successfully!"
	pageData["mood_id"] = newMood.ID
	return pageData
}
