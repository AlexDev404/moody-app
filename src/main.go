package main

import (
	"bytes"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"baby-blog/database"
	"baby-blog/database/models"
	"baby-blog/forms"
	"baby-blog/forms/validator"
	"baby-blog/hooks"
	"baby-blog/types"
	"html/template"
	"sync"
)

// Application is a wrapper for types.Application
type Application struct {
	*types.Application
	templates    *template.Template
	models       *types.Models
	bufferPool   sync.Pool
}

func (app *Application) runHooks(pageData map[string]interface{}) map[string]interface{} {
	return hooks.Hooks(pageData, app.models)
}

func (app *Application) startup() {
	addr := flag.String("addr", "4000", "HTTP network address")
	dsn := flag.String("dsn", "postgresql://postgres:postgres@localhost:5432/baby_blog?sslmode=disable", "PostgreSQL DSN")
	flag.Parse()

	fileServer := http.FileServer(http.Dir("static"))

	mux := http.NewServeMux()
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

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

	typesApp := &types.Application{
		Logger: logger,
	}

	app = &Application{
		Application: typesApp,
		templates:   templates,
		models: &types.Models{
			Feedback: &models.FeedbackModel{Database: db},
			Journal:  &models.JournalModel{Database: db},
			Todo:     &models.TodoModel{Database: db},
		},
	}

	app.bufferPool = sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}

	log.Printf("Templates loaded: %v", templates.DefinedTemplates())
	log.Printf("App templates: %v", app.templates.DefinedTemplates())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, app.templates, nil)
	})

	mux.HandleFunc("/journal", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.Render(w, r, app.templates, nil)
		case http.MethodPost:
			forms.JournalForm(w, r, validator.NewValidator())
			//app.JournalHandler(w, r)
		default:
			http.Error(w, , http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.Render(w, r, app.templates, nil)
		case http.MethodPost:
			app.TodoHandler(w, r)
		default:
			http.Error(w, MainServerMethodNotAllowedMessage, http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/feedback", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.Render(w, r, app.templates, nil)
		case http.MethodPost:
			app.FeedbackHandler(w, r)
		default:
			http.Error(w, MainServerMethodNotAllowedMessage, http.StatusMethodNotAllowed)
		}
	})

	app.Logger.Info("Now listening on port http://127.0.0.1:" + *addr)

	err := http.ListenAndServe((":" + *addr), app.Middleware.LoggingMiddleware(mux))

	if err != nil {
		panic(err.Error())
	}
}

func main() {
	app := &Application{}
	app.startup()
}
