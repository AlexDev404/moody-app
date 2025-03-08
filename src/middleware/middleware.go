package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type Application struct{}

// A basic middleware
func (app *Application) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("request", "method", r.Method, "url", r.URL.Path, "time", time.Since(start).String())
		next.ServeHTTP(w, r)
	})
}
