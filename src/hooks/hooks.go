package hooks

import (
	"baby-blog/middleware"
	"baby-blog/types"
	"log/slog"
	"net/http"
	"os"
)

type HooksConnector struct {
	Logger *slog.Logger
}

// Hooks processes the page data before it's sent to the template renderer.
// It initializes a custom logger for hooks operations, sets up a hooks connector,
// and enhances the page data with application information and any data from hook functions.
//
// Parameters:
//   - pageData: A map containing the initial data to be passed to the template
//   - dbModels: Database models for data access operations
//   - r: The HTTP request object containing client information and request details
//   - w: The HTTP response writer for modifying the response
//
// Returns:
//   - An enhanced map with additional data from hooks and application defaults
func Hooks(pageData map[string]interface{}, dbModels *types.Models, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	// Create a handler that prepends "[HOOKS]" to log messages
	opts := &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: false,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)

	// Create the hooks connector with the custom logger
	hooks := HooksConnector{
		Logger: slog.New(handler).With(slog.String("component", "Hooks")),
	}

	// Add authentication information to page data
	userID := middleware.GetUserID(r)
	if userID != "" {
		user, err := dbModels.Users.GetByID(userID)
		if err == nil && user != nil {
			pageData["IsAuthenticated"] = true
			pageData["User"] = user
		} else {
			pageData["IsAuthenticated"] = false
		}
	} else {
		pageData["IsAuthenticated"] = false
	}

	// Default hook: Render PageData for the current page (if any)
	pageData = hooks.PageLoad(pageData, dbModels, r, w)

	pageData["AppName"] = "Moody"
	pageData["AppVersion"] = "0.1.0"

	return pageData
}
