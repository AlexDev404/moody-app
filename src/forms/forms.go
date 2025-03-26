package forms

import (
	"log/slog"
	"os"
)

// Create a new logger instance
var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))
