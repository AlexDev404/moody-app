package pages

import "github.com/AlexDev404/folio"

func (p *Pages) Tools(c *folio.Ctx) error {
	allMoods, err := p.Models.Moods.GetAll()
	if err != nil {
		c.Logger.Error("loading moods", "error", err)
		c.Fail("Failed to load mood data", nil)
		return nil
	}

	c.Set("moods", allMoods)
	c.Set("has_mood_today", len(allMoods) > 0)
	return nil
}
