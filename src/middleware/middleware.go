package middleware

import (
	"baby-blog/auth"
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Application struct{}

// User context key for storing user ID
type userCtxKey string

// UserContextKey is the key used to store the user ID in the request context
const UserContextKey userCtxKey = "user"

// PublicPaths contains the paths that don't require authentication
var PublicPaths = []string{
	"/login",
	"/register",
	"/static/",
}

// A basic middleware
func (app *Application) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		slog.Info("request", "method", r.Method, "url", r.URL.Path, "time", time.Since(start).String())
		next.ServeHTTP(w, r)
	})
}

// AuthMiddleware checks for valid JWT in cookies and adds userID to context
func (app *Application) AuthMiddleware(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the path is in the public paths list
			for _, path := range PublicPaths {
				if strings.HasPrefix(r.URL.Path, path) {
					next.ServeHTTP(w, r)
					return
				}
			}

			// Get the session cookie
			cookie, err := r.Cookie("session")
			if err != nil {
				// No cookie found, redirect to login page with "next" parameter
				redirectToLogin(w, r)
				return
			}

			// Validate the token
			claims, err := jwtManager.ValidateToken(cookie.Value)
			if err != nil {
				// Invalid token, redirect to login page
				redirectToLogin(w, r)
				return
			}

			// Add user ID to the request context
			ctx := context.WithValue(r.Context(), UserContextKey, claims.UserID)
			// Call the next handler with our new context
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// GetUserID extracts the user ID from the request context
func GetUserID(r *http.Request) string {
	userID, ok := r.Context().Value(UserContextKey).(string)
	if !ok {
		return ""
	}
	return userID
}

// RequireAuthentication ensures that only authenticated users can access certain handlers
func (app *Application) RequireAuthentication(jwtManager *auth.JWTManager, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := GetUserID(r)
		if userID == "" {
			redirectToLogin(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// redirectToLogin redirects to the login page with the current URL as the "next" parameter
func redirectToLogin(w http.ResponseWriter, r *http.Request) {
	// Build the "next" URL parameter
	next := url.QueryEscape(r.URL.Path)
	http.Redirect(w, r, "/login?next="+next, http.StatusSeeOther)
}
