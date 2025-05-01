package pages

import (
	"baby-blog/types"
	"log"
	"net/http"
)

func Tools(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	// Get today's mood entry if it exists
	allMoods, err := db.Moods.GetAllWithPlaylist()
	if err != nil {
		log.Printf("Error getting mood data: %v", err)
		pageData["Failure"] = "Failed to load mood data"
		return pageData
	}

	if allMoods != nil {
		log.Print("Today's mood entry found")
		pageData["moods"] = allMoods
		pageData["has_mood_today"] = true
	} else {
		log.Print("Today's mood entry not found")
		pageData["allMoods"] = nil
		pageData["has_mood_today"] = false
	}
	return pageData
}
