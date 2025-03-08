package types

import (
	"baby-blog/middleware"
	"log/slog"
)

type Application struct {
	Logger     *slog.Logger
	Middleware *middleware.Application
}

type TemplateData struct {
	Data interface{}
}
