package models

import (
	"context"
	"database/sql"
	"time"
)

type Track struct {
	ID         string `json:"id"`
	Artist     string `json:"artist"`
	Title      string `json:"title"`
	YouTubeURL string `json:"youtube_url"`
}

type TrackModel struct {
	Database *sql.DB
}

func (t *TrackModel) Insert(track *Track, playlistID string) error {
	query := `
		INSERT INTO tracks (artist, title, youtube_url, playlist_id)
		VALUES ($1, $2, $3, $4)
		RETURNING id`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return t.Database.QueryRowContext(
		ctx,
		query,
		track.Artist,
		track.Title,
		track.YouTubeURL,
		playlistID,
	).Scan(&track.ID)
}

func (t *TrackModel) Get(id string) (*Track, error) {
	query := `
		SELECT id, artist, title, youtube_url 
		FROM tracks
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	track := &Track{}
	err := t.Database.QueryRowContext(ctx, query, id).Scan(
		&track.ID,
		&track.Artist,
		&track.Title,
		&track.YouTubeURL,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return track, err
}

func (t *TrackModel) GetAllForPlaylist(playlistID string) ([]Track, error) {
	query := `
		SELECT id, artist, title, youtube_url
		FROM tracks 
		WHERE playlist_id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.Database.QueryContext(ctx, query, playlistID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var track Track
		err := rows.Scan(
			&track.ID,
			&track.Artist,
			&track.Title,
			&track.YouTubeURL,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}

func (t *TrackModel) GetAll() ([]Track, error) {
	query := `
		SELECT id, artist, title, youtube_url
		FROM tracks
		ORDER BY id DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := t.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tracks []Track
	for rows.Next() {
		var track Track
		err := rows.Scan(
			&track.ID,
			&track.Artist,
			&track.Title,
			&track.YouTubeURL,
		)
		if err != nil {
			return nil, err
		}
		tracks = append(tracks, track)
	}
	return tracks, nil
}
