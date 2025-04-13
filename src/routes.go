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
	return mux
}
