// Filename: main.go

package main

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./templates/home.html")
}

func weekHandler(w http.ResponseWriter, r *http.Request) {
	// Extract the week number from the URL path
	week := strings.TrimPrefix(r.URL.Path, "/week")
	if week == "" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "./templates/week"+week+".html")
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
	// Register the handlers
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/week1", weekHandler)
	mux.HandleFunc("/week2", weekHandler)
	mux.HandleFunc("/week3", weekHandler)
	mux.HandleFunc("/week4", weekHandler)
	mux.HandleFunc("/week5", weekHandler)
	// Start listening for requests (start the web server)
	log.Println("Starting server on :4000")
	err := http.ListenAndServe(":4000", loggingMiddleware(mux))
	// Log error message if server quits unexpectedly
	if err != nil {
		log.Fatal(err)
	}
}
