package types

import (
	"baby-blog/database/models"
	"log/slog"
)

// Application type for the main application
type Application struct {
	Logger *slog.Logger
	// We'll initialize this in main.go after creating middleware
	Middleware interface{}
}

type TemplateData struct {
	Data            interface{}
	User            *AuthenticatedUser
	IsAuthenticated bool
	Flash           string
}

type AIPlaylist struct {
	Songs []AIPlaylistSong `json:"songs,omitempty"`
}

type AIPlaylistSong struct {
	Title    string   `json:"title"`
	Artist   string   `json:"artist"`
	MoodTags []string `json:"mood_tags"`
	Reason   string   `json:"reason"`
}

type AuthenticatedUser struct {
	IsAuthenticated bool
	ID              string
	Email           string
}

type Models struct {
	Moods     *models.MoodModel
	Playlists *models.PlaylistModel
	Tracks    *models.TrackModel
	Users     *models.UserModel
}
