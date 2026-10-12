// Package pages holds the Load and Action hooks for each route.
package pages

import (
	"context"
	"regexp"

	"baby-blog/database/models"
	"baby-blog/playlist"

	"github.com/AlexDev404/folio"
)

// userKey is the session value holding the signed-in user's ID.
const userKey = "user_id"

// PlaylistGenerator suggests songs for a mood.
type PlaylistGenerator interface {
	Generate(ctx context.Context, mood string) ([]playlist.Song, error)
}

type Pages struct {
	Models    *models.Models
	Playlists PlaylistGenerator
}

// Register wires every page into app.
func Register(app *folio.App, p *Pages) {
	// Everything except the login and registration pages needs a login.
	app.Use(app.RequireLogin(folio.LoginConfig{
		SessionKey: userKey,
		LoginURL:   "/login",
		// /logout is public so a signed-out visit doesn't become ?next=/logout,
		// which would sign the user out again right after logging in.
		Public: []string{"/register", "/logout"},
	}))
	app.LoadAll(p.CurrentUser)

	app.Action("login", p.Login)
	app.Load("login", p.LoginPage)
	app.Action("register", p.Register)
	app.Load("register", p.RegisterPage)
	app.HandleFunc("/logout", Logout)

	app.Action("index", p.SubmitMood)
	app.Load("index", p.Home)
	app.Load("tools", p.Tools)
	app.Action("mood/[id]", p.GeneratePlaylist)
	app.Load("mood/[id]", p.Mood)
}

// CurrentUser makes the signed-in user available to every page as .Page.User.
func (p *Pages) CurrentUser(c *folio.Ctx) error {
	c.Set("IsAuthenticated", false)
	id := userID(c)
	if id == "" {
		return nil
	}
	user, err := p.Models.Users.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		// The account is gone; drop the stale session.
		c.Session().Clear()
		return nil
	}
	c.Set("User", user)
	c.Set("IsAuthenticated", true)
	return nil
}

func userID(c *folio.Ctx) string { return c.Session().Get(userKey) }

var uuidPattern = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
