package models

import "database/sql"

type Models struct {
	Moods     MoodModel
	Playlists PlaylistModel
	Tracks    TrackModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Moods:     MoodModel{Database: db},
		Playlists: PlaylistModel{Database: db},
		Tracks:    TrackModel{Database: db},
	}
}
