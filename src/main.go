package main

import (
	"context"
	"flag"
	"log/slog"
	"os"

	"baby-blog/database"
	"baby-blog/database/models"
	"baby-blog/pages"
	"baby-blog/playlist"

	"github.com/AlexDev404/folio"
)

func main() {
	addr := flag.String("addr", "4000", "HTTP network address")
	dsn := flag.String("dsn", "postgresql://postgres:postgres@localhost:5432/moody?sslmode=disable", "PostgreSQL DSN")
	openaiKey := flag.String("openai-key", os.Getenv("OPENAI_API_KEY"), "API key for GitHub Models (default $OPENAI_API_KEY)")
	dev := flag.Bool("dev", false, "Reload templates and the browser on changes")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	if *openaiKey == "" {
		logger.Error("An API key for playlist generation is required: pass --openai-key or set OPENAI_API_KEY.")
		os.Exit(1)
	}

	db, err := database.OpenDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		logger.Error(`This error usually results when the application could not connect to the database.
Ensure that PostgreSQL is installed and running on port 5432--otherwise pass a different URL to it using the
flag --dsn=URL`)
		os.Exit(1)
	}
	// release the database resources before exiting
	defer db.Close()
	logger.Info("database connection pool established")

	m := models.New(db)
	if err := m.Sessions.Ready(context.Background()); err != nil {
		logger.Error(err.Error())
		logger.Error(`Logins are stored in the sessions table, which this database doesn't have yet.
Create it by running migrations/000003_add_sessions.up.sql against the database, e.g.
  psql "$DB_DSN" -f ../migrations/000003_add_sessions.up.sql`)
		os.Exit(1)
	}
	app, err := folio.New(folio.Config{
		Logger: logger,
		Reload: *dev,
		// Sessions live in Postgres (migration 000003), so logging out ends
		// them on the server and they survive restarts.
		SessionStore: m.Sessions,
		Funcs:        pages.Funcs,
		Globals: map[string]any{
			"AppName":    "Moody",
			"AppVersion": "0.1.0",
		},
	})
	if err != nil {
		logger.Error("starting folio", "error", err)
		os.Exit(1)
	}

	app.Use(folio.RequestLogger(logger), folio.Recoverer(logger), folio.Compress())
	pages.Register(app, &pages.Pages{Models: m, Playlists: playlist.NewGenerator(*openaiKey)})

	if err := app.ListenAndServe(":" + *addr); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
