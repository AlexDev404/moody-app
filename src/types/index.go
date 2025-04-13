package types

import (
	"baby-blog/database/models"
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

type Models struct {
	Mood *models.MoodModel
}
