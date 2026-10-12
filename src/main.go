package main

import (
	"flag"
	"log/slog"
	"os"

	"baby-blog/database"
	"baby-blog/database/models"
	"baby-blog/pages"

	"github.com/AlexDev404/folio"
)

func main() {
	addr := flag.String("addr", "4000", "HTTP network address")
	dsn := flag.String("dsn", "postgresql://postgres:postgres@localhost:5432/moody?sslmode=disable", "PostgreSQL DSN")
	dev := flag.Bool("dev", false, "Reload templates and the browser on changes")
	secret := flag.String("secret", os.Getenv("MOODY_SECRET"), "Session secret, 32+ characters (default $MOODY_SECRET)")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

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

	app, err := folio.New(folio.Config{
		Logger:    logger,
		Reload:    *dev,
		SecretKey: []byte(*secret),
		Globals: map[string]any{
			"AppName":    "Moody",
			"AppVersion": "0.1.0",
		},
	})
	if err != nil {
		logger.Error("parsing templates", "error", err)
		os.Exit(1)
	}

	app.Use(folio.RequestLogger(logger), folio.Recoverer(logger), folio.Compress())
	pages.Register(app, models.New(db))

	if err := app.ListenAndServe(":" + *addr); err != nil {
		logger.Error("server stopped", "error", err)
		os.Exit(1)
	}
}
