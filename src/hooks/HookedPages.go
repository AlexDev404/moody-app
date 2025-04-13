package hooks

import (
	"baby-blog/hooks/pages"
	"baby-blog/types"
	"net/http"
)

// Define the map of pages that need hooks, with associated functions
var hookedPages = map[string]func(pageData map[string]interface{}, db *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{}{
	"index": pages.Home,
}
