package pages

func FeedbackGallery(pageData map[string]interface{}) map[string]interface{} {
	// This function is called when the feedback gallery page is loaded
	// You can add your logic here

	// Example: Add a title to the page data
	pageData["feedbacks"] = []map[string]string{
		{
			"name":    "John Doe",
			"date":    "2023-10-15",
			"message": "Great service! I really enjoyed working with your team.",
		},
		{
			"name":    "Jane Smith",
			"date":    "2023-10-10",
			"message": "The product exceeded my expectations. Would recommend!",
		},
		{
			"name":    "Mike Johnson",
			"date":    "2023-10-05",
			"message": "",
		},
		{
			"name":    "Sarah Williams",
			"date":    "2023-09-28",
			"message": "Quick delivery and excellent customer support.",
		},
	}

	return pageData
}
