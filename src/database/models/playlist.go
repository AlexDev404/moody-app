package models

import (
	"context"
	"database/sql"
	"time"
)

type Playlist struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	Tracks    []Track   `json:"tracks"`
}

type PlaylistModel struct {
	Database *sql.DB
}

func (p *PlaylistModel) Insert(playlist *Playlist, moodID string) error {
	query := `
		INSERT INTO playlists (name, mood_id)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.Database.QueryRowContext(
		ctx,
		query,
		playlist.Name,
		moodID,
	).Scan(&playlist.ID, &playlist.CreatedAt)
}

func (p *PlaylistModel) Get(id string) (*Playlist, error) {
	query := `
		SELECT id, name, created_at
		FROM playlists
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	playlist := &Playlist{}
	err := p.Database.QueryRowContext(ctx, query, id).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return playlist, err
}

func (p *PlaylistModel) GetForMood(moodID string) (*Playlist, error) {
	query := `
		SELECT id, name, created_at
		FROM playlists
		WHERE mood_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	playlist := &Playlist{}
	err := p.Database.QueryRowContext(ctx, query, moodID).Scan(
		&playlist.ID,
		&playlist.Name,
		&playlist.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return playlist, err
}

func (p *PlaylistModel) GetAll() ([]Playlist, error) {
	query := `
		SELECT id, name, created_at
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
			&playlist.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		playlists = append(playlists, playlist)
	}
	return playlists, nil
}
