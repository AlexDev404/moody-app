package hooks

import "baby-blog/hooks/pages"

func (hooks *HooksConnector) PageLoad(pageData map[string]interface{}) map[string]interface{} {
	// This function is called when the page is loaded
	// You can add your logic here
	pageData = pages.FeedbackGallery(pageData)

	return pageData
}

// This is a placeholder for the PageLoad function
