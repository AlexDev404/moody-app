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

	"baby-blog/database"
	"baby-blog/database/models"
	"baby-blog/hooks"
	"baby-blog/types"
	"html/template"
	"sync"
)

// Application is a wrapper for types.Application
type Application struct {
	*types.Application
	templates  *template.Template
	models     *types.Models
	bufferPool sync.Pool
	JournalModel *models.JournalModel
	TodoModel    *models.TodoModel
}

func (app *Application) Render(w http.ResponseWriter, r *http.Request, t *template.Template, pageData map[string]interface{}) {
	path := app.getPath(r)
	print(path)
	if app.isDisallowedRoute(path) {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	if pageData == nil {
		pageData = map[string]interface{}{
			"Path":    path,
			"Errors":  map[string]string{},
			"Message": nil,
			"Failure": nil,
		}
	}

	pageData = app.runHooks(pageData)

	contentData := &types.TemplateData{
		Data: map[string]interface{}{
			"Path":     path,
			"PageData": pageData,
		},
	}

	var templateContent string
	var tmpl *template.Template = t.Lookup(path)
	var contentErr error
	if path != "/" {
		if tmpl != nil {
			var buf bytes.Buffer
			contentErr = tmpl.Execute(&buf, contentData)
			templateContent = buf.String()
		} else {
			w.WriteHeader(http.StatusNotFound)
			http.ServeFile(w, r, "./static/errors/404.html")
			return
		}
	}

	data := &types.TemplateData{
		Data: map[string]interface{}{
			"Path":     path,
			"HTML":     template.HTML(templateContent),
			"PageData": pageData,
		},
	}

	layout := t.Lookup("layout/" + path)
	pageBuf := app.bufferPool.Get().(*bytes.Buffer)
	pageBuf.Reset()
	defer app.bufferPool.Put(pageBuf)

	var err error
	if contentErr == nil {
		if layout == nil {
			err = t.ExecuteTemplate(pageBuf, "layout/app", data)
		} else {
			err = layout.Execute(pageBuf, data)
		}
		if err != nil {
			app.Logger.Error("Error rendering page", "error", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}
		_, err = pageBuf.WriteTo(w)
	}
	if contentErr != nil || err != nil {
		app.Logger.Error("Failed to write template to response: ", "error", err, "content_err", contentErr)
		w.WriteHeader(http.StatusInternalServerError)
		http.ServeFile(w, r, "./static/errors/500.html")
	}
}

func (app *Application) getPath(r *http.Request) string {
	path := r.URL.Path[1:]
	if path == "" {
		path = "index"
	} else {
		path = strings.TrimSuffix(path, "/")
	}
	return path
}

func (app *Application) isDisallowedRoute(path string) bool {
	disallowedRoutes := []string{"context", "head", "header", "footer", "current_ctx"}
	for _, route := range disallowedRoutes {
		if path == route {
			return true
		}
	}
	return false
}

func (app *Application) runHooks(pageData map[string]interface{}) map[string]interface{} {
	return hooks.Hooks(pageData, app.models)
}

func getTemplates() (*template.Template, error) {
	log.Println("Parsing 'initial' templates...")
	templates, err := template.ParseGlob("templates/*.tmpl")
	if err != nil {
		return nil, err
	}
	log.Println("Parsing 'partial' templates...")
	templates, err = templates.ParseGlob("templates/partials/*.tmpl")
	if err != nil {
		return nil, err
	}
	log.Println("Parsing 'route' templates...")
	err = filepath.Walk("templates/routes", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".tmpl") {
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
		},
		JournalModel: &models.JournalModel{Database: db},
		TodoModel:    &models.TodoModel{Database: db},
	}

	app.bufferPool = sync.Pool{
		New: func() any {
			return &bytes.Buffer{}
		},
	}

	log.Printf("Templates loaded: %v", templates.DefinedTemplates())
	log.Printf("App templates: %v", app.templates.DefinedTemplates())

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.Render(w, r, app.templates, nil)
		case http.MethodPost:
			app.POSTHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/journal", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.Render(w, r, app.templates, nil)
		case http.MethodPost:
			app.JournalHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/todo", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			app.Render(w, r, app.templates, nil)
		case http.MethodPost:
			app.TodoHandler(w, r)
		default:
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
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

func (app *Application) JournalHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.Logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	validator := validator.NewValidator()
	formData, formErrors := forms.JournalForm(w, r, validator)

	for key, value := range formErrors {
		formData[key] = value
	}

	if formErrors == nil {
		journalData := &models.Journal{
			Title:   formData["title"].(string),
			Content: formData["content"].(string),
		}
		err := app.JournalModel.Insert(journalData)
		if err != nil {
			formData["Failure"] = "✗ Failed to submit journal entry. Please try again later."
		} else {
			formData["Message"] = "✓ Your journal entry has been submitted. Thank you!"
		}
	}

	app.Render(w, r, app.templates, formData)
}

func (app *Application) TodoHandler(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.Logger.Error("Failed to parse form", "error", err)
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	validator := validator.NewValidator()
	formData, formErrors := forms.TodoForm(w, r, validator)

	for key, value := range formErrors {
		formData[key] = value
	}

	if formErrors == nil {
		todoData := &models.Todo{
			Task:      formData["task"].(string),
			Completed: false,
		}
		err := app.TodoModel.Insert(todoData)
		if err != nil {
			formData["Failure"] = "✗ Failed to submit todo item. Please try again later."
		} else {
			formData["Message"] = "✓ Your todo item has been submitted. Thank you!"
		}
	}

	app.Render(w, r, app.templates, formData)
}
