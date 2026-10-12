package pages

import (
	"io"
	"log/slog"
	"net/url"
	"os"
	"testing"

	"github.com/AlexDev404/folio"
	"github.com/AlexDev404/folio/foliotest"
)

// newTestClient serves the real templates. The tests stick to pages that
// don't touch the database, so Models stays empty.
func newTestClient(t *testing.T) *foliotest.Client {
	t.Helper()
	app, err := folio.New(folio.Config{
		Templates: os.DirFS("../templates"),
		Static:    os.DirFS("../static"),
		Funcs:     Funcs,
		Logger:    slog.New(slog.NewTextHandler(io.Discard, nil)),
	})
	if err != nil {
		t.Fatal(err)
	}
	Register(app, &Pages{})
	return foliotest.New(t, app.Handler())
}

func TestLoginRequired(t *testing.T) {
	c := newTestClient(t)
	c.Get("/").ExpectRedirect("/login?next=%2F")
	c.Get("/tools?period=yesterday").ExpectRedirect("/login?next=%2Ftools%3Fperiod%3Dyesterday")
	c.Get("/mood/00000000-0000-0000-0000-000000000000").ExpectStatus(303)

	c.Header.Set("Accept", "application/json")
	c.Get("/tools").ExpectStatus(401)
}

func TestPublicPages(t *testing.T) {
	c := newTestClient(t)
	c.Get("/login").ExpectStatus(200).ExpectContains("<h1").ExpectContains(`action="/login"`)
	c.Get("/login?next=/tools").ExpectContains(`action="/login?next=%2ftools"`)
	c.Get("/login?next=//evil.example").ExpectContains(`action="/login"`).ExpectNotContains("evil.example")
	c.Get("/register").ExpectStatus(200).ExpectContains("Create an Account")
	c.Get("/static/style.css").ExpectStatus(200)
	c.Get("/partial/ux/navigation").ExpectStatus(303) // not a page; the login check runs first
}

func TestRegisterValidation(t *testing.T) {
	c := newTestClient(t)
	c.PostForm("/register", url.Values{
		"email":            {"not-an-email"},
		"password":         {"short"},
		"password_confirm": {"different"},
	}).
		ExpectStatus(422).
		ExpectContains("Enter a valid email address").
		ExpectContains("Password must be at least 8 characters long").
		ExpectContains("Passwords do not match").
		ExpectContains(`value="not-an-email"`)
}

func TestLogout(t *testing.T) {
	c := newTestClient(t)
	c.Get("/logout").ExpectRedirect("/login")
}

func TestPeriodRange(t *testing.T) {
	now := mustTime(t, "2026-10-12T15:04:05+08:00")
	for key, want := range map[string]string{
		"today":          "2026-10-12T00:00:00+08:00",
		"yesterday":      "2026-10-11T00:00:00+08:00",
		"four-days-ago":  "2026-10-08T00:00:00+08:00",
		"something-else": "2026-10-12T00:00:00+08:00",
	} {
		start, end := periodRange(key, now)
		if start.Format(timeLayout) != want || end.Sub(start).Hours() != 24 {
			t.Errorf("periodRange(%q) = %v - %v", key, start, end)
		}
	}
}

func TestCapitalizeFirst(t *testing.T) {
	for in, want := range map[string]string{"hello world": "Hello world", "élan": "Élan", "": "", "Taylor Swift fits": "Taylor Swift fits"} {
		if got := CapitalizeFirst(in); got != want {
			t.Errorf("CapitalizeFirst(%q) = %q, want %q", in, got, want)
		}
	}
}
