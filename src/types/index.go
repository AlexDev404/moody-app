package types

import (
	"baby-blog/middleware"
	"html/template"
	"log/slog"
)

type Application struct {
	Logger     *slog.Logger
	Middleware *middleware.Application
}

type TemplateData struct {
	Title string
	Body  template.HTML
	Data  interface{}
}
