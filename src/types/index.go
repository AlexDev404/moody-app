package types

import (
	"baby-blog/database/models"
	"baby-blog/middleware"
	"log/slog"
)

type Application struct {
	Logger     *slog.Logger
	Middleware *middleware.Application
}

type TemplateData struct {
	Data interface{}
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

type Models struct {
	Moods     *models.MoodModel
	Playlists *models.PlaylistModel
	Tracks    *models.TrackModel
}
