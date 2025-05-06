package main

import (
	"baby-blog/database/models"
	"baby-blog/middleware"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

// HandleLogin handles login requests
func (app *Application) HandleLogin(w http.ResponseWriter, r *http.Request) {
	pageData := map[string]interface{}{
		"NextURL": r.URL.Query().Get("next"),
	}

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

		// Authenticate the user
		user, err := app.models.Users.Authenticate(email, password)
		if err != nil {
			// Render the login page with an error message
			pageData["Error"] = "Invalid email or password"
			pageData["Email"] = email
			app.Render(w, r, app.templates, pageData)
			return
		}

		// Generate a JWT token
		token, err := app.jwtManager.GenerateToken(user.ID)
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

		// Redirect safely using the utility function
		RedirectToSafeURL(w, r, "/", app.Logger)
		return
	}

	// For GET requests, render the login page
	app.Render(w, r, app.templates, pageData)
}

// HandleRegister handles user registration
func (app *Application) HandleRegister(w http.ResponseWriter, r *http.Request) {
	pageData := map[string]interface{}{
		"NextURL": r.URL.Query().Get("next"),
		"Errors":  map[string]string{},
	}

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
		existingUser, err := app.models.Users.GetByEmail(email)
		if err != nil {
			http.Error(w, "Could not check user existence", http.StatusInternalServerError)
			return
		}

		if existingUser != nil {
			errors["email"] = "Email is already in use"
		}

		// If there are validation errors, re-render the form
		if len(errors) > 0 {
			pageData["Email"] = email
			pageData["Errors"] = errors
			app.Render(w, r, app.templates, pageData)
			return
		}

		// Create the user
		user := &models.User{
			Email: email,
		}

		err = app.models.Users.Insert(user, password)
		if err != nil {
			http.Error(w, "Could not create user", http.StatusInternalServerError)
			return
		}

		// Generate a JWT token
		token, err := app.jwtManager.GenerateToken(user.ID)
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

		// Redirect safely using the utility function
		RedirectToSafeURL(w, r, "/", app.Logger)
		return
	}

	// For GET requests, render the registration page
	app.Render(w, r, app.templates, pageData)
}

// HandleLogout handles user logout
func (app *Application) HandleLogout(w http.ResponseWriter, r *http.Request) {
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

// GetCurrentUser returns the user from the current request context
func (app *Application) GetCurrentUser(r *http.Request) (*models.User, error) {
	userID := middleware.GetUserID(r)
	if userID == "" {
		return nil, nil
	}
	return app.models.Users.GetByID(userID)
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// RespondJSON sends a JSON response
func RespondJSON(w http.ResponseWriter, status int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(response)
}

// HandleMoodFilter handles the mood filter AJAX requests
func (app *Application) HandleMoodFilter(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Check if user is authenticated
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the time filter from the query string
	timeFilter := r.URL.Query().Get("timeFilter")
	if timeFilter == "" {
		timeFilter = "today" // Default to today
	}

	// Get the mood entries based on the time filter
	entries, err := app.getMoodEntriesByTimeFilter(userID, timeFilter)
	if err != nil {
		app.respondWithError(w, err)
		return
	}

	// Create response data structure
	type FilterResponse struct {
		Success bool        `json:"success"`
		Data    interface{} `json:"data"`
	}

	// If no entries found, return empty response
	if len(entries) == 0 {
		response := FilterResponse{
			Success: true,
			Data: map[string]interface{}{
				"playlist": map[string]string{
					"name": app.getTimeFilterDisplayName(timeFilter),
				},
				"tracks": []interface{}{},
			},
		}

		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			app.respondWithError(w, err)
		}
		return
	}

	// Use the first entry's playlist as our response data
	var responseData map[string]interface{}

	// Find the first entry with a playlist
	var playlist *models.Playlist
	var tracks []map[string]interface{}

	for _, entry := range entries {
		if entry.Playlist != nil {
			playlist = entry.Playlist

			// Convert tracks to the format expected by the frontend
			for _, track := range playlist.Tracks {
				tracks = append(tracks, map[string]interface{}{
					"artist":    track.Artist,
					"title":     track.Title,
					"mood_tags": []string{},
					"reason":    "",
				})
			}
			break
		}
	}

	// Create the response data
	if playlist != nil {
		responseData = map[string]interface{}{
			"playlist": map[string]string{
				"name": playlist.Name,
			},
			"tracks": tracks,
		}
	} else {
		// If no playlist found, create empty response
		responseData = map[string]interface{}{
			"playlist": map[string]string{
				"name": app.getTimeFilterDisplayName(timeFilter),
			},
			"tracks": []interface{}{},
		}
	}

	// Send the response
	response := FilterResponse{
		Success: true,
		Data:    responseData,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.respondWithError(w, err)
	}
}

// HandleToolsFilter handles the tools filter AJAX requests
func (app *Application) HandleToolsFilter(w http.ResponseWriter, r *http.Request) {
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Check if user is authenticated
	userID := middleware.GetUserID(r)
	if userID == "" {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get the time filter from the query string
	timeFilter := r.URL.Query().Get("timeFilter")
	if timeFilter == "" {
		timeFilter = "today" // Default to today
	}

	// Get mood entries based on time filter
	moodEntries, err := app.getMoodEntriesByTimeFilter(userID, timeFilter)
	if err != nil {
		app.respondWithError(w, err)
		return
	}

	// Format moods for the response
	var moodsData []map[string]interface{}
	for _, entry := range moodEntries {
		moodsData = append(moodsData, map[string]interface{}{
			"id":          entry.ID,
			"text":        entry.MoodText,
			"hasPlaylist": entry.Playlist != nil,
		})
	}

	// Create the response data
	responseData := map[string]interface{}{
		"title": app.getTimeFilterDisplayName(timeFilter),
		"moods": moodsData,
	}

	// Send the response
	response := struct {
		Success bool                   `json:"success"`
		Data    map[string]interface{} `json:"data"`
	}{
		Success: true,
		Data:    responseData,
	}

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		app.respondWithError(w, err)
	}
}

// respondWithError sends a JSON error response
func (app *Application) respondWithError(w http.ResponseWriter, err error) {
	response := APIResponse{
		Success: false,
		Error:   err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(w).Encode(response)
}

// getMoodEntriesByTimeFilter retrieves mood entries based on the selected time filter
func (app *Application) getMoodEntriesByTimeFilter(userID string, timeFilter string) ([]models.MoodEntry, error) {
	// Get today's date at the start of the day
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// Calculate the start date based on the time filter
	var startDate time.Time
	var endDate time.Time

	switch timeFilter {
	case "today":
		startDate = today
		endDate = today.Add(24 * time.Hour)
	case "yesterday":
		startDate = today.AddDate(0, 0, -1)
		endDate = today
	case "day-before":
		startDate = today.AddDate(0, 0, -2)
		endDate = today.AddDate(0, 0, -1)
	case "three-days-ago":
		startDate = today.AddDate(0, 0, -3)
		endDate = today.AddDate(0, 0, -2)
	case "four-days-ago":
		startDate = today.AddDate(0, 0, -4)
		endDate = today.AddDate(0, 0, -3)
	default:
		startDate = today
		endDate = today.Add(24 * time.Hour)
	}

	// Get all mood entries for the user
	allEntries, err := app.models.Moods.GetAllWithPlaylistForUser(userID)
	if err != nil {
		return nil, err
	}

	// Filter the entries by date
	var filteredEntries []models.MoodEntry
	for _, entry := range allEntries {
		if (entry.CreatedAt.After(startDate) || entry.CreatedAt.Equal(startDate)) && entry.CreatedAt.Before(endDate) {
			filteredEntries = append(filteredEntries, entry)
		}
	}

	return filteredEntries, nil
}

// getTimeFilterDisplayName returns a user-friendly name for the time filter
func (app *Application) getTimeFilterDisplayName(timeFilter string) string {
	switch timeFilter {
	case "today":
		return "Today's Moods"
	case "yesterday":
		return "Yesterday's Moods"
	case "day-before":
		return "The Day Before's Moods"
	case "three-days-ago":
		return "Three Days Ago Moods"
	case "four-days-ago":
		return "Four Days Ago Moods"
	default:
		return "Moods"
	}
}
