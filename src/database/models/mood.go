package models

import (
	"context"
	"database/sql"
	"time"
)

type MoodEntry struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	MoodText  string    `json:"mood_text"`
	Playlist  *Playlist `json:"playlist"`
}

type MoodModel struct {
	Database *sql.DB
}

func (m *MoodModel) Insert(entry *MoodEntry) error {
	query := `
		INSERT INTO mood_entries (mood_text)
		VALUES ($1)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.Database.QueryRowContext(
		ctx,
		query,
		entry.MoodText,
	).Scan(&entry.ID, &entry.CreatedAt)
}

func (m *MoodModel) GetToday() (*MoodEntry, error) {
	query := `
		SELECT id, created_at, mood_text
		FROM mood_entries 
		WHERE DATE(created_at) = CURRENT_DATE
		ORDER BY created_at DESC
		LIMIT 1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	entry := &MoodEntry{}
	err := m.Database.QueryRowContext(ctx, query).Scan(
		&entry.ID,
		&entry.CreatedAt,
		&entry.MoodText,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}

func (m *MoodModel) GetTodayWithPlaylist() (*MoodEntry, error) {
	query := `
		SELECT me.id, me.created_at, me.mood_text,
			   p.id, p.name, p.created_at,
			   t.id, t.artist, t.title, t.youtube_url
		FROM mood_entries me
		LEFT JOIN playlists p ON p.mood_id = me.id
		LEFT JOIN tracks t ON t.playlist_id = p.id
		WHERE DATE(me.created_at) = CURRENT_DATE
		ORDER BY me.created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entry *MoodEntry
	trackMap := make(map[string]Track)

	for rows.Next() {
		if entry == nil {
			entry = &MoodEntry{
				Playlist: &Playlist{},
			}
		}

		var trackID, trackArtist, trackTitle, trackURL sql.NullString

		err := rows.Scan(
			&entry.ID,
			&entry.CreatedAt,
			&entry.MoodText,
			&entry.Playlist.ID,
			&entry.Playlist.Name,
			&entry.Playlist.CreatedAt,
			&trackID,
			&trackArtist,
			&trackTitle,
			&trackURL,
		)
		if err != nil {
			return nil, err
		}

		if trackID.Valid {
			track := Track{
				ID:         trackID.String,
				Artist:     trackArtist.String,
				Title:      trackTitle.String,
				YouTubeURL: trackURL.String,
			}
			trackMap[trackID.String] = track
		}
	}

	if entry != nil && entry.Playlist != nil {
		for _, track := range trackMap {
			entry.Playlist.Tracks = append(entry.Playlist.Tracks, track)
		}
	}

	return entry, nil
}

func (m *MoodModel) GetAll() ([]MoodEntry, error) {
	query := `
		SELECT id, created_at, mood_text
		FROM mood_entries
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []MoodEntry
	for rows.Next() {
		var entry MoodEntry
		err := rows.Scan(
			&entry.ID,
			&entry.CreatedAt,
			&entry.MoodText,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (m *MoodModel) GetByID(id string) (*MoodEntry, error) {
	query := `
		SELECT id, created_at, mood_text
		FROM mood_entries
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	entry := &MoodEntry{}
	err := m.Database.QueryRowContext(ctx, query, id).Scan(
		&entry.ID,
		&entry.CreatedAt,
		&entry.MoodText,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}
