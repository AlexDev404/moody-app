package pages

import (
	"net/http"
	"time"

	"baby-blog/database/models"
	"baby-blog/forms"

	"github.com/AlexDev404/folio"
)

func (p *Pages) Home(c *folio.Ctx) error {
	moods, err := p.Models.Moods.GetAllForUser(userID(c))
	if err != nil {
		return err
	}
	start, end := periodRange("today", time.Now())
	today := 0
	for _, m := range moods {
		if !m.CreatedAt.Before(start) && m.CreatedAt.Before(end) {
			today++
		}
	}
	c.Set("emotions_felt", today)
	return nil
}

// SubmitMood saves a check-in and goes to its playlist page, which creates
// the playlist.
func (p *Pages) SubmitMood(c *folio.Ctx) error {
	form := forms.NewMood(c.Request)
	if errs := form.Validate(); errs != nil {
		c.Status(http.StatusUnprocessableEntity)
		c.Fail("✗ Please check your errors and try again.", errs)
		c.Set("mood", form.Mood)
		return nil
	}

	entry := &models.MoodEntry{MoodText: form.Mood, UserID: userID(c)}
	if err := p.Models.Moods.Insert(entry); err != nil {
		return err
	}
	c.Flash("✓ Mood submitted successfully!")
	c.Redirect("/mood/"+entry.ID, http.StatusSeeOther)
	return nil
}
