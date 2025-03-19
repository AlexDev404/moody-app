package pages

import (
	"baby-blog/types"
)

func FeedbackGallery(pageData map[string]interface{}, db *types.Models) map[string]interface{} {
	// This function is called when the feedback gallery page is loaded
	// You can add your logic here

	// Fetch all feedback from the database
	feedbacks, err := db.Feedback.GetAll()
	if err != nil {
		// Log the error
		pageData["Failure"] = "Failed to load feedback"
		return pageData
	}

	// Convert database feedbacks to the format expected by the template
	feedbacksData := make([]map[string]string, 0, len(feedbacks))
	for _, feedback := range feedbacks {
		feedbacksData = append(feedbacksData, map[string]string{
			"name":    feedback.Fullname,
			"date":    feedback.CreatedAt.Format("2006-01-02"),
			"message": feedback.Message,
		})
	}

	// Replace the sample data with actual data from the database
	pageData["feedbacks"] = feedbacksData

	return pageData
}
