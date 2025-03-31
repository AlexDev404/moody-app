package main

import "net/http"

func (app *Application) routes() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	// log.Printf("Templates loaded: %v", templates.DefinedTemplates())
	// log.Printf("App templates: %v", app.templates.DefinedTemplates())

	// Handlers
	// Root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, app.templates, nil)
	})

	// GET /journals
	mux.HandleFunc("GET /journals", func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, app.templates, nil)
	})

	// Dual purpose handler for GET /journal and POST /journal
	mux.HandleFunc("/journal", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.Render(w, r, app.templates, nil)
		case http.MethodPost:
			app.JournalHandler(w, r)
		default:
			http.Error(w, MainServerMethodNotAllowedMessage, http.StatusMethodNotAllowed)
		}
	})

	// GET /todos
	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, app.templates, nil)
	})

	// Dual purpose handler for GET /todo and POST /todo
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

	// GET /feedbacks
	mux.HandleFunc("GET /feedbacks", func(w http.ResponseWriter, r *http.Request) {
		app.Render(w, r, app.templates, nil)
	})

	// Dual purpose handler for GET /feedback and POST /feedback
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

	return mux
}
