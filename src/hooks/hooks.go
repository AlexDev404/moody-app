package hooks

import (
	"baby-blog/types"
	"log/slog"
	"os"
)

type HooksConnector struct {
	Logger *slog.Logger
}

func Hooks(pageData map[string]interface{}, dbModels *types.Models) map[string]interface{} {
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

	// Default hook: Render PageData for the current page (if any)
	pageData = hooks.PageLoad(pageData, dbModels)

	pageData["AppName"] = "Moody"
	pageData["AppVersion"] = "0.1.0"

	return pageData
}
