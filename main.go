// Filename: main.go

package main

import (
	"log"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./templates/home.html") 
}

func week1Handler(w http.ResponseWriter, r *http.Request) {
        http.ServeFile(w, r, "./templates/week1.html")
}

func main() {
        mux := http.NewServeMux()

        // Serve the static files
        fileServer := http.FileServer(http.Dir("./static"))
        mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
        // Register the handlers
        mux.HandleFunc("/", homeHandler) 
        mux.HandleFunc("/week1", week1Handler)
        // Start listening for requests (start the web server)
        log.Println("Starting server on :4000")
        err := http.ListenAndServe(":4000", mux)
	// Log error message if server quits unexpectedly
        if err != nil {
                log.Fatal(err)
        }
}