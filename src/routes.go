package main

import "net/http"

func (app *Application) routes() (mux *http.ServeMux) {
	mux = http.NewServeMux()
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))

	// Root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			app.Render(w, r, app.templates, nil)
			return
		}
		app.Render(w, r, app.templates, nil)
	})

	// Auth routes
	mux.HandleFunc("/login", app.HandleLogin)
	mux.HandleFunc("/register", app.HandleRegister)
	mux.HandleFunc("/logout", app.HandleLogout)

	return mux
}
