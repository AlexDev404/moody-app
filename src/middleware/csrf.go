package middleware

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"os"
	"sync"

	"github.com/gorilla/csrf"
)

// type for context key to avoid collisions
type contextKey string

// context key for CSRF token
const csrfTokenKey contextKey = "csrf.token"

// SessionTokens stores CSRF tokens by session ID (for dev mode only)
var (
	sessionTokens = make(map[string]string)
	tokenMutex    sync.RWMutex
)

// CSRFMiddleware creates a new CSRF protection middleware
func (app *Application) CSRFMiddleware() func(http.Handler) http.Handler {
	// Get CSRF secret key from environment or use a default for development
	csrfKey := os.Getenv("CSRF_SECRET")
	if csrfKey == "" {
		csrfKey = "32-byte-long-auth-key-for-csrf-dev" // Default for development
	}

	// Check if we're in development mode
	isDevelopment := os.Getenv("APP_ENV") != "production"

	if isDevelopment {
		// In development mode, use a relaxed CSRF middleware
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// For safe methods (GET, HEAD, OPTIONS), just generate a token
				if r.Method == http.MethodGet || r.Method == http.MethodHead || r.Method == http.MethodOptions {
					// Generate a token for the session
					token := generateToken(r)

					// Set CSRF cookie for browser-side access
					http.SetCookie(w, &http.Cookie{
						Name:     "csrf_token",
						Value:    token,
						Path:     "/",
						HttpOnly: false, // Allow JavaScript access
						SameSite: http.SameSiteLaxMode,
					})

					// Store token in the request context
					ctx := context.WithValue(r.Context(), csrfTokenKey, token)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				}

				// For unsafe methods (POST, PUT, DELETE, etc.)
				// Simply check if a token is present without validating origin
				token := extractCSRFToken(r)
				if token == "" {
					// No token provided, show error
					csrfErrorHandler().ServeHTTP(w, r)
					return
				}

				// Token exists, allow the request without checking the origin
				ctx := context.WithValue(r.Context(), csrfTokenKey, token)
				next.ServeHTTP(w, r.WithContext(ctx))
			})
		}
	}

	// In production, use standard CSRF protection with all validations
	return csrf.Protect(
		[]byte(csrfKey),
		csrf.CookieName("csrf_token"),
		csrf.Path("/"),
		csrf.Secure(true), // Requires HTTPS in production
		csrf.SameSite(csrf.SameSiteStrictMode),
		csrf.ErrorHandler(csrfErrorHandler()),
	)
}

// generateToken creates a CSRF token for the current session
func generateToken(r *http.Request) string {
	// Get the session cookie
	cookie, err := r.Cookie("session")
	sessionID := ""
	if err == nil && cookie != nil {
		sessionID = cookie.Value
	} else {
		// If no session cookie, generate a random session ID
		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		sessionID = base64.StdEncoding.EncodeToString(randomBytes)
	}

	// Check if we have a token for this session
	tokenMutex.RLock()
	token, exists := sessionTokens[sessionID]
	tokenMutex.RUnlock()

	if !exists {
		// Generate a new token
		randomBytes := make([]byte, 32)
		rand.Read(randomBytes)
		token = base64.StdEncoding.EncodeToString(randomBytes)

		// Save it
		tokenMutex.Lock()
		sessionTokens[sessionID] = token
		tokenMutex.Unlock()
	}

	return token
}

// extractCSRFToken tries to extract the CSRF token from a request
func extractCSRFToken(r *http.Request) string {
	// First try to parse the form
	r.ParseForm()

	// Check form values
	if token := r.PostForm.Get("csrf_token"); token != "" {
		return token
	}

	if token := r.PostForm.Get("gorilla.csrf.Token"); token != "" {
		return token
	}

	// Check headers
	if token := r.Header.Get("X-CSRF-Token"); token != "" {
		return token
	}

	// Try to get from cookie as fallback
	cookie, err := r.Cookie("csrf_token")
	if err == nil && cookie != nil {
		return cookie.Value
	}

	return ""
}

// csrfErrorHandler returns a custom error handler for CSRF failures
func csrfErrorHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if it's an API request
		if isAPIRequest(r) {
			http.Error(w, `{"error": "CSRF token validation failed"}`, http.StatusForbidden)
			return
		}

		// For regular requests, show a friendly error page
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte(`
			<html>
				<head>
					<title>CSRF Validation Failed</title>
					<style>
						body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, Helvetica, Arial, sans-serif; line-height: 1.6; color: #333; max-width: 650px; margin: 0 auto; padding: 20px; }
						h1 { color: #d63939; }
						.container { background: #ffefef; border: 1px solid #f8d7da; padding: 20px; border-radius: 5px; }
						a { color: #0071e3; text-decoration: none; }
						a:hover { text-decoration: underline; }
					</style>
				</head>
				<body>
					<div class="container">
						<h1>Security Error</h1>
						<p>CSRF token validation failed. This could be due to:</p>
						<ul>
							<li>Your session has expired</li>
							<li>You're attempting to submit a form from an unauthorized source</li>
							<li>Your browser cookies may be disabled</li>
						</ul>
						<p><a href="/">Return to homepage</a> and try again.</p>
						<p>Error details: CSRF token missing or invalid</p>
					</div>
				</body>
			</html>
		`))
	})
}

// isAPIRequest checks if the request is an API request based on the Accept header or URL path
func isAPIRequest(r *http.Request) bool {
	// Check Accept header for JSON
	accept := r.Header.Get("Accept")
	if accept == "application/json" {
		return true
	}

	// Check if the path starts with /api/
	if len(r.URL.Path) >= 5 && r.URL.Path[0:5] == "/api/" {
		return true
	}

	return false
}

// GetCSRFToken returns the CSRF token for the current request
func GetCSRFToken(r *http.Request) string {
	// First try the standard gorilla CSRF method
	token := csrf.Token(r)
	if token != "" {
		return token
	}

	// Fallback to our custom implementation
	if token, ok := r.Context().Value(csrfTokenKey).(string); ok {
		return token
	}

	return ""
}
