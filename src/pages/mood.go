package pages

import (
	"context"
	"net/http"
	"time"

	"baby-blog/database/models"

	"github.com/AlexDev404/folio"
)

// MoodView is the data for a mood's playlist page.
type MoodView struct {
	Mood     *models.MoodEntry
	Playlist *models.Playlist
	Tracks   []models.Track
}

// ownedMood loads the mood in the URL, answering 404 unless it belongs to
// the signed-in user.
func (p *Pages) ownedMood(c *folio.Ctx) (*models.MoodEntry, error) {
	id := c.Param("id")
	if !uuidPattern.MatchString(id) {
		return nil, folio.ErrNotFound
	}
	mood, err := p.Models.Moods.GetByID(id)
	if err != nil {
		return nil, err
	}
	if mood == nil || mood.UserID != userID(c) {
		return nil, folio.ErrNotFound
	}
	return mood, nil
}

// Mood shows a mood's playlist. Without one, the page creates it straight
// away by posting to GeneratePlaylist.
func (p *Pages) Mood(c *folio.Ctx) error {
	mood, err := p.ownedMood(c)
	if err != nil {
		return err
	}
	view := MoodView{Mood: mood}
	view.Playlist, err = p.Models.Playlists.GetForMoodForUser(mood.ID, mood.UserID)
	if err != nil {
		return err
	}
	if view.Playlist != nil {
		if view.Tracks, err = p.Models.Tracks.GetAllForPlaylist(view.Playlist.ID); err != nil {
			return err
		}
	}
	c.SetModel(view)
	return nil
}

// GeneratePlaylist asks the model for songs and saves them as the mood's
// playlist.
func (p *Pages) GeneratePlaylist(c *folio.Ctx) error {
	mood, err := p.ownedMood(c)
	if err != nil {
		return err
	}
	existing, err := p.Models.Playlists.GetForMoodForUser(mood.ID, mood.UserID)
	if err != nil {
		return err
	}
	if existing == nil {
		ctx, cancel := context.WithTimeout(c.Request.Context(), 60*time.Second)
		defer cancel()
		songs, err := p.Playlists.Generate(ctx, mood.MoodText)
		if err != nil {
			c.Logger.Error("generating playlist", "mood", mood.ID, "error", err)
			c.Status(http.StatusBadGateway)
			c.Fail("✗ We couldn't create a playlist right now. Please try again.", nil)
			return nil
		}

		tracks := make([]models.Track, len(songs))
		for i, s := range songs {
			tracks[i] = models.Track{Artist: s.Artist, Title: s.Title, Reason: s.Reason}
		}
		pl := &models.Playlist{MoodID: mood.ID, Name: mood.MoodText, UserID: mood.UserID}
		if err := p.Models.Playlists.InsertWithTracks(pl, tracks); err != nil {
			return err
		}
		c.Flash("✓ Your playlist is ready!")
	}
	c.Redirect("/mood/"+mood.ID, http.StatusSeeOther)
	return nil
}
