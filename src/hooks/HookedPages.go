package hooks

import (
	"baby-blog/hooks/pages"
	"baby-blog/types"
)

// Define the map of pages that need hooks, with associated functions
var hookedPages = map[string]func(pageData map[string]interface{}, db *types.Models) map[string]interface{}{
	"feedbacks": pages.FeedbackGallery,
}
