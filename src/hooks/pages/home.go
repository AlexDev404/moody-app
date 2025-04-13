package pages

import (
	"baby-blog/types"
	"net/http"
)

func Home(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	// Get today's mood entry if it exists
	todaysMood, err := db.Moods.GetToday()
	if err != nil {
		pageData["error"] = "Failed to load mood data"
		return pageData
	}

	if todaysMood != nil {
		pageData["emotions_felt"] = 1
		pageData["current_mood"] = todaysMood.MoodText
		pageData["has_mood_today"] = true
	} else {
		pageData["emotions_felt"] = 0
		pageData["has_mood_today"] = false
	}

	return pageData
}
