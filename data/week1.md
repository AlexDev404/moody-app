# Week 1: Introduction to Go Web Apps

<p class="mb-4 text-lg">
  This week, I learned how to set up a basic web server in Go using the
  `net/http` package.
</p>

<pre class="bg-gray-100 p-4 rounded border border-gray-300 overflow-x-auto">
    <code class="font-mono text-sm leading-relaxed">
      package main

      import (
          "log"
          "net/http"
      )

      func home(w http.ResponseWriter, r *http.Request) {
          w.Write([]byte("Hello from my first Go web app!"))
      }

      func main() {
          mux := http.NewServeMux()
          mux.HandleFunc("/", home)

          log.Println("Server starting on :4000")
          err := http.ListenAndServe(":4000", mux)
          if err != nil {
              log.Fatal(err)
          }
      }
  </code>
</pre>
<p class="mb-4 text-lg">
  I also learned about routing, serving static files, and the importance of a
  structured project layout.
</p>
