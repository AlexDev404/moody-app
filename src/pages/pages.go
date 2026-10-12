// Package pages holds the Load and Action hooks for each route.
package pages

import (
	"baby-blog/database/models"

	"github.com/AlexDev404/folio"
)

type Pages struct {
	Models *models.Models
}

func Register(app *folio.App, m *models.Models) {
	p := &Pages{Models: m}

	app.Action("index", p.SubmitMood)
	app.Load("index", p.Home)
	app.Load("tools", p.Tools)
}
