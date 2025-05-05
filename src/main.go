package main

import (
	"bytes"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"baby-blog/auth"
	"baby-blog/database"
	"baby-blog/database/models"
	"baby-blog/hooks"
	"baby-blog/middleware"
	"baby-blog/types"
	"html/template"
	"sync"
)

// Application is a wrapper for types.Application
type Application struct {
	*types.Application
	templates  *template.Template
	models     *types.Models
	jwtManager *auth.JWTManager
	bufferPool sync.Pool
}

func (app *Application) runHooks(pageData map[string]interface{}, r *http.Request, w http.ResponseWriter) map[string]interface{} {
	// Get the current user from the request context if available
	user, _ := app.GetCurrentUser(r)

	// Add the user to the page data if authenticated
	if user != nil {
		pageData["User"] = user
		pageData["IsAuthenticated"] = true
	} else {
		pageData["IsAuthenticated"] = false
	}

	// Run the existing hooks
	return hooks.Hooks(pageData, app.models, r, w)
}

func (app *Application) startup() {
	addr := flag.String("addr", "4000", "HTTP network address")
	dsn := flag.String("dsn", "postgresql://postgres:postgres@localhost:5432/moody?sslmode=disable", "PostgreSQL DSN")
	openaiKey := flag.String("openai-key", "", "OpenAI API key")
	jwtSecret := flag.String("jwt-secret", "", "JWT secret key for authentication")
	flag.Parse()

	if *openaiKey == "" {
		log.Fatal("OpenAI API key is required. Please pass the OPENAI_API_KEY flag.")
	}

	if *jwtSecret == "" {
		log.Fatal("JWT secret key is required. Please pass the --jwt-secret flag.")
	}

	os.Setenv("OPENAI_API_KEY", *openaiKey)
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	templates, tErr := getTemplates()
	if tErr != nil {
		log.Fatal("Error parsing templates: ", tErr)
	}

	db, dbErr := database.OpenDB(*dsn)
	if dbErr != nil {
		log.Print(dbErr.Error())
		log.Fatal(`This error usually results when the application could not connect to the database.
Ensure that PostgreSQL is installed and running on port 5432--otherwise pass a different URL to it using the
flag --dsn=URL`)
		os.Exit(1)
	}
	// release the database resources before exiting
	defer db.Close()

	logger.Info("database connection pool established")

	// Create JWT manager
	jwtManager := auth.NewJWTManager(*jwtSecret)

	// Create middleware instance
	middlewareApp := &middleware.Application{}

	typesApp := &types.Application{
		Logger:     logger,
		Middleware: middlewareApp,
	}

	app = &Application{
		Application: typesApp,
		templates:   templates,
		jwtManager:  jwtManager,
		models: &types.Models{
			Moods:     &models.MoodModel{Database: db},
			Playlists: &models.PlaylistModel{Database: db},
			Tracks:    &models.TrackModel{Database: db},
			Users:     &models.UserModel{Database: db},
		},
	}

	app.bufferPool = sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}

	// Register the routes
	mux := app.routes()

	app.Logger.Info("Now listening on port http://127.0.0.1:" + *addr)
	
	// Apply middleware chain: logging -> authentication -> routes
	middlewareChain := middlewareApp.LoggingMiddleware(
		middlewareApp.AuthMiddleware(jwtManager)(mux),
	)

	err := http.ListenAndServe((":" + *addr), middlewareChain)

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	app := &Application{}
	app.startup()
}
