package pages

import (
	"net/http"

	"baby-blog/database/models"
	"baby-blog/forms"

	"github.com/AlexDev404/folio"
)

func (p *Pages) Home(c *folio.Ctx) error {
	allMoods, err := p.Models.Moods.GetAll()
	if err != nil {
		c.Logger.Error("loading moods", "error", err)
		c.Fail("Failed to load mood data", nil)
		return nil
	}

	c.Set("emotions_felt", len(allMoods))
	c.Set("has_mood_today", len(allMoods) > 0)
	return nil
}

func (p *Pages) SubmitMood(c *folio.Ctx) error {
	form := forms.NewMoodForm(c.Request)
	if errs := form.Validate(); errs != nil {
		c.Logger.Warn("invalid mood submission", "errors", errs)
		c.Status(http.StatusUnprocessableEntity)
		c.Fail("✗ Please check your errors and try again.", errs)
		return nil
	}

	if err := p.Models.Moods.Insert(&models.MoodEntry{MoodText: form.Mood}); err != nil {
		c.Logger.Error("inserting mood", "error", err)
		c.Status(http.StatusInternalServerError)
		c.Fail("✗ Failed to submit mood. Please try again later.", nil)
		return nil
	}

	// Redirect so refreshing the page doesn't submit the mood again.
	c.Flash("✓ Mood submitted successfully!")
	c.Redirect("/", http.StatusSeeOther)
	return nil
}
