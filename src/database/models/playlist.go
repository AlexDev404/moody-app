package models

import (
	"context"
	"database/sql"
	"time"
)

type Playlist struct {
	ID        string    `json:"id"`
	MoodID    string    `json:"mood_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Tracks    []Track   `json:"tracks"`
	UserID    string    `json:"user_id"`
}

type PlaylistModel struct {
	Database *sql.DB
}

func (p *PlaylistModel) Insert(playlist *Playlist) error {
	query := `
		INSERT INTO playlists (name, mood_id, user_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.Database.QueryRowContext(
		ctx,
		query,
		playlist.Name,
		playlist.MoodID,
		playlist.UserID,
	).Scan(&playlist.ID, &playlist.CreatedAt)
}

func (p *PlaylistModel) GetByID(id string) (*Playlist, error) {
	query := `
		SELECT id, name, mood_id, created_at, user_id
		FROM playlists
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	playlist := &Playlist{}
	err := p.Database.QueryRowContext(ctx, query, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.MoodID,
		&playlist.CreatedAt,
		&playlist.UserID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return playlist, err
}

// Legacy method - calls the new method for compatibility
func (p *PlaylistModel) Get(id string) (*Playlist, error) {
	return p.GetByID(id)
}

// GetByIDForUser gets a playlist by ID and validates it belongs to the user
func (p *PlaylistModel) GetByIDForUser(id string, userID string) (*Playlist, error) {
	query := `
		SELECT id, name, mood_id, created_at, user_id
		FROM playlists
		WHERE id = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	playlist := &Playlist{}
	err := p.Database.QueryRowContext(ctx, query, id, userID).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.MoodID,
		&playlist.CreatedAt,
		&playlist.UserID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return playlist, err
}

func (p *PlaylistModel) GetForMood(moodID string) (*Playlist, error) {
	query := `
		SELECT id, name, mood_id, created_at, user_id
		FROM playlists
		WHERE mood_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	playlist := &Playlist{}
	err := p.Database.QueryRowContext(ctx, query, moodID).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.MoodID,
		&playlist.CreatedAt,
		&playlist.UserID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return playlist, err
}

// GetForMoodForUser gets a playlist for a specific mood and validates it belongs to the user
func (p *PlaylistModel) GetForMoodForUser(moodID string, userID string) (*Playlist, error) {
	query := `
		SELECT p.id, p.name, p.mood_id, p.created_at, p.user_id
		FROM playlists p
		JOIN mood_entries m ON p.mood_id = m.id
		WHERE p.mood_id = $1 AND m.user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	playlist := &Playlist{}
	err := p.Database.QueryRowContext(ctx, query, moodID, userID).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.MoodID,
		&playlist.CreatedAt,
		&playlist.UserID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return playlist, err
}

func (p *PlaylistModel) GetAllForUser(userID string) ([]Playlist, error) {
	query := `
		SELECT id, name, mood_id, created_at, user_id
		FROM playlists
		WHERE user_id = $1
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.Database.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var playlist Playlist
		err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.MoodID,
			&playlist.CreatedAt,
			&playlist.UserID,
		)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}

func (p *PlaylistModel) GetAll() ([]Playlist, error) {
	query := `
		SELECT id, name, mood_id, created_at, user_id
		FROM playlists
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := p.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var playlists []Playlist
	for rows.Next() {
		var playlist Playlist
		err := rows.Scan(
			&playlist.ID,
			&playlist.Name,
			&playlist.MoodID,
			&playlist.CreatedAt,
			&playlist.UserID,
		)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}
