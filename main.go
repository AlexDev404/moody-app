// Filename: main.go

package main

import (
	"html/template"
	"log"
	"log/slog"
	"net/http"
	"os"
	"time"
)

type application struct {
	logger *slog.Logger
}

type templateData struct {
	Title    string
	Greeting string
}

func (app *application) viewTemplate(w http.ResponseWriter, r *http.Request) {
	// Get the URL path
	path := r.URL.Path

	// If the path is the root, serve the index template
	// Otherwise, serve the template that corresponds to the path
	if path == "/" {
		http.ServeFile(w, r, "./templates/index.mustache")
	} else {
		// Trim the "/week/" prefix from the path
		path = path[len("/week/"):]
		p := templateData{
			Title:    "Go Study",
			Greeting: "Hello World Study Mates",
		}

		// Parse the template file
		t, _ := template.ParseFiles("./templates/" + "week" + path + ".mustache")
		t.Execute(w, p)
	}
}

// A basic middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Print("Method: ", r.Method, " URL: ", r.URL.Path, " Time: ", time.Since(start))
		next.ServeHTTP(w, r)
	})
}

func main() {
	mux := http.NewServeMux()

	// Serve the static files
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	application := &application{logger: logger}
	// Register the handler
	mux.HandleFunc("/", application.viewTemplate)
	// Start listening for requests (start the web server)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", loggingMiddleware(mux))
	// Log error message if server quits unexpectedly
	if err != nil {
		log.Fatal(err)
	}
}
