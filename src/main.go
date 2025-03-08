// Filename: main.go

package main

import (
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"

	"baby-blog/types"
	"html/template"
)

// Application is a wrapper for types.Application
type Application struct {
	*types.Application
	templates *template.Template
}

func (app *Application) ViewTemplate(w http.ResponseWriter, r *http.Request, t *template.Template) {
	// Get the URL path
	path := r.URL.Path

	// TemplateData is a struct that holds the title, body, and data for the template
	data := &types.TemplateData{
		Title: "Baby Blog",
		Body:  template.HTML("<h1>Welcome to Baby Blog</h1>"),
		Data: map[string]interface{}{
			"Path": path,
		},
	}

	// If the path is the root, serve the index template
	// Otherwise, serve the template that corresponds to the path
	// if path == "/" {
	// template, err := template.ParseFiles("./templates/index.mustache")
	err := t.ExecuteTemplate(w, "app", data)
	if err != nil {
		app.Logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	// } else {

	// http.ServeFile(w, r, "./templates/404.mustache")
	// }
}

func getTemplates() (*template.Template, error) {
	// Parse the initial templates
	log.Println("Parsing 'initial' templates...")
	templates, err := template.ParseGlob("templates/*.mustache")
	if err != nil {
		return nil, err
	}
	// Add the partials to the templates
	log.Println("Parsing 'partial' templates...")
	templates, err = templates.ParseGlob("templates/partials/*.mustache")
	if err != nil {
		return nil, err
	}
	// Add the routes to the templates
	log.Println("Parsing 'route' templates...")
	templates, err = templates.ParseGlob("templates/routes/*.mustache")
	if err != nil {
		return nil, err
	}
	return templates, nil
}

// func init() {
// 	Application.templates, _ := getTemplates()

// }

func main() {
	addr := flag.String("addr", "4000", "HTTP network address")
	flag.Parse()
	mux := http.NewServeMux()
	// Serve the static files
	fileServer := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	templates, t_err := getTemplates()
	if t_err != nil {
		log.Fatalf("Error parsing templates: %v", t_err)
	}
	typesApp := &types.Application{
		Logger: logger,
	}
	app := &Application{
		Application: typesApp,
		templates:   templates,
	}

	log.Printf("Templates loaded: %v", templates.DefinedTemplates())
	log.Printf("App templates: %v", app.templates.DefinedTemplates())
	// Register the handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		app.ViewTemplate(w, r, app.templates)
	})
	// Start listening for requests (start the web server)
	err := http.ListenAndServe((":" + *addr), app.Middleware.LoggingMiddleware(mux))
	// Log error message if server quits unexpectedly
	if err != nil {
		panic(err.Error())
	}
	app.Logger.Info("Now listening on port " + *addr)
}
