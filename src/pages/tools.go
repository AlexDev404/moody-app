package pages

import (
	"time"

	"baby-blog/database/models"

	"github.com/AlexDev404/folio"
)

// Period is one tab of the mood history filter.
type Period struct {
	Key   string
	Label string
	Title string
}

var periods = []Period{
	{"today", "Today", "Today's Moods"},
	{"yesterday", "Yesterday", "Yesterday's Moods"},
	{"day-before", "The day before", "The Day Before's Moods"},
	{"three-days-ago", "Three days ago", "Three Days Ago Moods"},
	{"four-days-ago", "Four days ago", "Four Days Ago Moods"},
}

// ToolsView is the data for the mood history page.
type ToolsView struct {
	Periods []Period
	Current Period
	Moods   []models.MoodEntry
}

// Tools lists the signed-in user's moods for the period in ?period=.
func (p *Pages) Tools(c *folio.Ctx) error {
	current := periods[0]
	for _, period := range periods {
		if period.Key == c.Request.URL.Query().Get("period") {
			current = period
		}
	}

	all, err := p.Models.Moods.GetAllWithPlaylistForUser(userID(c))
	if err != nil {
		return err
	}
	start, end := periodRange(current.Key, time.Now())
	var moods []models.MoodEntry
	for _, m := range all {
		if !m.CreatedAt.Before(start) && m.CreatedAt.Before(end) {
			moods = append(moods, m)
		}
	}

	c.SetModel(ToolsView{Periods: periods, Current: current, Moods: moods})
	return nil
}

// periodRange returns the start and end of a period key, in now's time zone.
func periodRange(key string, now time.Time) (time.Time, time.Time) {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	daysAgo := map[string]int{"yesterday": 1, "day-before": 2, "three-days-ago": 3, "four-days-ago": 4}[key]
	start := today.AddDate(0, 0, -daysAgo)
	return start, start.AddDate(0, 0, 1)
}
