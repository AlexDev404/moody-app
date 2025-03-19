package hooks

import (
	"baby-blog/hooks/pages"
	"baby-blog/types"
)

func (hooks *HooksConnector) PageLoad(pageData map[string]interface{}, dbModels *types.Models) map[string]interface{} {
	// This function is called when the page is loaded
	// You can add your logic here
	pageData = pages.FeedbackGallery(pageData, dbModels)

	return pageData
}
