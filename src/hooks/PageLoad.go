package hooks

import (
	"baby-blog/types"
	"net/http"
)

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
