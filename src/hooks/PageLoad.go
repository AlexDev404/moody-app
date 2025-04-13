package hooks

import (
	"baby-blog/types"
	"net/http"
)

// PageLoad processes the current page data when a page is loaded.
// It checks if the current page path is registered in the hookedPages map and, if so,
// calls the corresponding page handler function to preprocess the page data before rendering.
//
// Parameters:
//   - pageData: A map containing the page's data, including the "Path" key representing the current URL path
//   - dbModels: Database models required for data operations
//   - r: The HTTP request object
//   - w: The HTTP response writer
//
// Returns:
//   - The processed page data map, which may be modified by page-specific handlers
func (hooks *HooksConnector) PageLoad(pageData map[string]interface{}, dbModels *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	// This function is called when the page is loaded
	// Check if current page is in the hooked pages list
	currentPath, ok := pageData["Path"].(string)
	if ok {
		for path, pageHandler := range hookedPages {
			if currentPath == path {
				// Pre-render pageData for included pages and pass them back
				// Down to the renderer
				pageData = pageHandler(pageData, dbModels, r, w)
			}
		}
	}

	return pageData
}
