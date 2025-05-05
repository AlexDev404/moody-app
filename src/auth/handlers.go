package auth

import (
	"baby-blog/database/models"
	"html/template"
	"net/http"
	"strings"
	"time"
)

// ModelInterface defines an interface for accessing user models
type ModelInterface interface {
	Authenticate(email, password string) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	Insert(user *models.User, password string) error
}

// AuthHandlers encapsulates the authentication handlers
type AuthHandlers struct {
	Templates  *template.Template
	UserModel  ModelInterface
	JWTManager *JWTManager
}

// LoginHandler handles user login requests
func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Process the form if it's a POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Could not process form", http.StatusInternalServerError)
			return
		}

		// Extract the email and password from the form
		email := r.PostForm.Get("email")
		password := r.PostForm.Get("password")

		// Get the "next" parameter if it exists
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}

		// Authenticate the user
		user, err := h.UserModel.Authenticate(email, password)
		if err != nil {
			// Render the login page with an error message
			data := struct {
				Error   string
				Email   string
				NextURL string
			}{
				Error:   "Invalid email or password",
				Email:   email,
				NextURL: next,
			}
			h.Templates.ExecuteTemplate(w, "login.tmpl", data)
			return
		}

		// Generate a JWT token
		token, err := h.JWTManager.GenerateToken(user.ID)
		if err != nil {
			http.Error(w, "Could not generate authentication token", http.StatusInternalServerError)
			return
		}

		// Set the token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: http.SameSiteLaxMode,
		})

		// Redirect to the next URL or home page
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	// For GET requests, render the login page
	next := r.URL.Query().Get("next")
	data := struct {
		Error   string
		Email   string
		NextURL string
	}{
		NextURL: next,
	}
	h.Templates.ExecuteTemplate(w, "login.tmpl", data)
}

// RegisterHandler handles user registration
func (h *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Process the form if it's a POST request
	if r.Method == http.MethodPost {
		// Parse the form data
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Could not process form", http.StatusInternalServerError)
			return
		}

		// Extract the email, password, and password confirmation
		email := strings.TrimSpace(r.PostForm.Get("email"))
		password := r.PostForm.Get("password")
		passwordConfirm := r.PostForm.Get("password_confirm")

		// Get the "next" parameter if it exists
		next := r.URL.Query().Get("next")
		if next == "" {
			next = "/"
		}

		// Basic validation
		errors := make(map[string]string)

		// Validate email
		if email == "" {
			errors["email"] = "Email is required"
		}

		// Validate password
		if password == "" {
			errors["password"] = "Password is required"
		} else if len(password) < 8 {
			errors["password"] = "Password must be at least 8 characters long"
		}

		// Check if passwords match
		if password != passwordConfirm {
			errors["password_confirm"] = "Passwords do not match"
		}

		// Check if email is already in use
		existingUser, err := h.UserModel.GetByEmail(email)
		if err != nil {
			http.Error(w, "Could not check user existence", http.StatusInternalServerError)
			return
		}

		if existingUser != nil {
			errors["email"] = "Email is already in use"
		}

		// If there are validation errors, re-render the form
		if len(errors) > 0 {
			data := struct {
				Email   string
				Errors  map[string]string
				NextURL string
			}{
				Email:   email,
				Errors:  errors,
				NextURL: next,
			}
			h.Templates.ExecuteTemplate(w, "register.tmpl", data)
			return
		}

		// Create the user
		user := &models.User{
			Email: email,
		}

		err = h.UserModel.Insert(user, password)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		// Generate a JWT token
		token, err := h.JWTManager.GenerateToken(user.ID)
		if err != nil {
			http.Error(w, "Could not generate authentication token", http.StatusInternalServerError)
			return
		}

		// Set the token as a cookie
		http.SetCookie(w, &http.Cookie{
			Name:     "session",
			Value:    token,
			HttpOnly: true,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			SameSite: http.SameSiteLaxMode,
		})

		// Redirect to the next URL or home page
		http.Redirect(w, r, next, http.StatusSeeOther)
		return
	}

	// For GET requests, render the registration page
	next := r.URL.Query().Get("next")
	data := struct {
		Email   string
		Errors  map[string]string
		NextURL string
	}{
		NextURL: next,
	}
	h.Templates.ExecuteTemplate(w, "register.tmpl", data)
}

// LogoutHandler handles user logout
func (h *AuthHandlers) LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Clear the session cookie
	http.SetCookie(w, &http.Cookie{
		Name:     "session",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		MaxAge:   -1, // Expire immediately
		SameSite: http.SameSiteLaxMode,
	})

	// Redirect to the login page
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
