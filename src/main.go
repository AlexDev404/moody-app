// Filename: main.go

package main

import (
	"bytes"
	"flag"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"

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
	path := r.URL.Path[1:]
	log.Println("Path_RAW: ", path)
	disallowed_routes := []string{"context", "head", "header", "footer", "current_ctx", "index"}
	log.Println("Path: ", path)
	// Remove any trailing slashes
	if path == "" {
		path = "index"
	} else {
		path = strings.TrimSuffix(path, "/")
		// Check if path is in disallowed_routes
		for _, route := range disallowed_routes {
			if path == route {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}
		}
	}

	// TemplateData is a struct that holds the title, body, and data for the template
	// First try to get the template by path if it's not root
	var templateContent string
	var tmpl *template.Template = t.Lookup(path)
	if path != "/" {
		if tmpl != nil {
			// log.Println(tmpl.Tree.Root.String())
			// tmpl.New(path[1:])
			var buf bytes.Buffer
			tmpl.Execute(&buf, nil)
			templateContent = buf.String()
			// templateContent = "Template found: " + path[1:]
			// tmpl.Execute(w, nil)
		} else {
			http.ServeFile(w, r, "./static/errors/404.html")
			return
		}
	}

	data := &types.TemplateData{
		Data: map[string]interface{}{
			"Path": path,
			"HTML": template.HTML(templateContent),
		},
	}

	// Section: Render Layouts
	// First: Let's check if there's a layout for the path
	// Remove the leading text following the last / in the string
	path = strings.TrimSuffix(path, "/"+path[strings.LastIndex(path, "/")+1:])
	log.Println("LAYOUT_Path: ", path)
	layout := t.Lookup("layout/" + path)
	var err error
	if layout == nil {
		err = t.ExecuteTemplate(w, "layout/app", data)
	} else {
		layout.Execute(w, data)
	}
	if err != nil {
		app.Logger.Error("Error rendering template", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
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
	err = filepath.Walk("templates/routes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".mustache") {
			var tmpl *template.Template
			tmpl, err = templates.ParseFiles(path)
			if err != nil {
				return err
			}
			templates = tmpl
		}
		return nil
	})
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
