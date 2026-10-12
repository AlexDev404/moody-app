package models

import "database/sql"

type Models struct {
	Moods     *MoodModel
	Playlists *PlaylistModel
	Tracks    *TrackModel
	Users     *UserModel
	Sessions  *SessionStore
}

func New(db *sql.DB) *Models {
	return &Models{
		Moods:     &MoodModel{Database: db},
		Playlists: &PlaylistModel{Database: db},
		Tracks:    &TrackModel{Database: db},
		Users:     &UserModel{Database: db},
		Sessions:  &SessionStore{Database: db},
	}
}
