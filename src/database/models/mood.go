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
	UserID    string    `json:"user_id"`
}

type MoodModel struct {
	Database *sql.DB
}

func (m *MoodModel) Insert(entry *MoodEntry) error {
	query := `
		INSERT INTO mood_entries (mood_text, user_id)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.Database.QueryRowContext(
		ctx,
		query,
		entry.MoodText,
		entry.UserID,
	).Scan(&entry.ID, &entry.CreatedAt)
}

func (m *MoodModel) GetToday() (*MoodEntry, error) {
	query := `
		SELECT id, created_at, mood_text, user_id
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
		&entry.UserID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}

func (m *MoodModel) GetTodayWithPlaylist() (*MoodEntry, error) {
	query := `
		SELECT me.id, me.created_at, me.mood_text, me.user_id,
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
			&entry.UserID,
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
		SELECT id, created_at, mood_text, user_id
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
			&entry.UserID,
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
		SELECT id, created_at, mood_text, user_id
		FROM mood_entries
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	entry := &MoodEntry{}
	err := m.Database.QueryRowContext(ctx, query, id).Scan(
		&entry.ID,
		&entry.CreatedAt,
		&entry.MoodText,
		&entry.UserID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return entry, err
}

func (m *MoodModel) GetAllWithPlaylist() ([]MoodEntry, error) {
	query := `
		SELECT me.id, me.created_at, me.mood_text, me.user_id,
			   p.id, p.name, p.created_at,
			   t.id, t.artist, t.title, t.youtube_url
		FROM mood_entries me
		LEFT JOIN playlists p ON p.mood_id = me.id
		LEFT JOIN tracks t ON t.playlist_id = p.id
		ORDER BY me.created_at ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entriesMap := make(map[string]*MoodEntry)
	playlistsMap := make(map[string]*Playlist)

	for rows.Next() {
		var moodID string
		var playlistID, playlistName sql.NullString
		var playlistCreatedAt sql.NullTime
		var trackID, trackArtist, trackTitle, trackURL sql.NullString
		var moodText string
		var moodCreatedAt time.Time
		var userIDForMood string

		err := rows.Scan(
			&moodID,
			&moodCreatedAt,
			&moodText,
			&userIDForMood,
			&playlistID,
			&playlistName,
			&playlistCreatedAt,
			&trackID,
			&trackArtist,
			&trackTitle,
			&trackURL,
		)
		if err != nil {
			return nil, err
		}

		// Get or create mood entry
		entry, exists := entriesMap[moodID]
		if !exists {
			entry = &MoodEntry{
				ID:        moodID,
				CreatedAt: moodCreatedAt,
				MoodText:  moodText,
				UserID:    userIDForMood,
			}
			entriesMap[moodID] = entry
		}

		// If there's a playlist for this mood
		if playlistID.Valid && playlistName.Valid {
			// Get or create playlist
			playlist, exists := playlistsMap[playlistID.String]
			if !exists {
				playlist = &Playlist{
					ID:        playlistID.String,
					Name:      playlistName.String,
					CreatedAt: playlistCreatedAt.Time,
					Tracks:    []Track{},
				}
				playlistsMap[playlistID.String] = playlist
				entry.Playlist = playlist
			}

			// If there's a track for this playlist
			if trackID.Valid {
				track := Track{
					ID:         trackID.String,
					Artist:     trackArtist.String,
					Title:      trackTitle.String,
					YouTubeURL: trackURL.String,
				}
				// Avoid duplicate tracks
				trackExists := false
				for _, existingTrack := range playlist.Tracks {
					if existingTrack.ID == track.ID {
						trackExists = true
						break
					}
				}
				if !trackExists {
					playlist.Tracks = append(playlist.Tracks, track)
				}
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice
	var entries []MoodEntry
	for _, entry := range entriesMap {
		entries = append(entries, *entry)
	}

	return entries, nil
}

// GetAllForUser returns all mood entries for a specific user
func (m *MoodModel) GetAllForUser(userID string) ([]MoodEntry, error) {
	query := `
		SELECT id, created_at, mood_text, user_id
		FROM mood_entries
		WHERE user_id = $1
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Database.QueryContext(ctx, query, userID)
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
			&entry.UserID,
		)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return entries, nil
}

// GetAllWithPlaylistForUser returns all mood entries with their associated playlists for a specific user
func (m *MoodModel) GetAllWithPlaylistForUser(userID string) ([]MoodEntry, error) {
	query := `
		SELECT me.id, me.created_at, me.mood_text, me.user_id,
			   p.id, p.name, p.created_at,
			   t.id, t.artist, t.title, t.youtube_url
		FROM mood_entries me
		LEFT JOIN playlists p ON p.mood_id = me.id
		LEFT JOIN tracks t ON t.playlist_id = p.id
		WHERE me.user_id = $1
		ORDER BY me.created_at ASC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Database.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	entriesMap := make(map[string]*MoodEntry)
	playlistsMap := make(map[string]*Playlist)

	for rows.Next() {
		var moodID string
		var playlistID, playlistName sql.NullString
		var playlistCreatedAt sql.NullTime
		var trackID, trackArtist, trackTitle, trackURL sql.NullString
		var moodText string
		var moodCreatedAt time.Time
		var userIDForMood string

		err := rows.Scan(
			&moodID,
			&moodCreatedAt,
			&moodText,
			&userIDForMood,
			&playlistID,
			&playlistName,
			&playlistCreatedAt,
			&trackID,
			&trackArtist,
			&trackTitle,
			&trackURL,
		)
		if err != nil {
			return nil, err
		}

		// Get or create mood entry
		entry, exists := entriesMap[moodID]
		if !exists {
			entry = &MoodEntry{
				ID:        moodID,
				CreatedAt: moodCreatedAt,
				MoodText:  moodText,
				UserID:    userIDForMood,
			}
			entriesMap[moodID] = entry
		}

		// If there's a playlist for this mood
		if playlistID.Valid && playlistName.Valid {
			// Get or create playlist
			playlist, exists := playlistsMap[playlistID.String]
			if !exists {
				playlist = &Playlist{
					ID:        playlistID.String,
					Name:      playlistName.String,
					CreatedAt: playlistCreatedAt.Time,
					Tracks:    []Track{},
				}
				playlistsMap[playlistID.String] = playlist
				entry.Playlist = playlist
			}

			// If there's a track for this playlist
			if trackID.Valid {
				track := Track{
					ID:         trackID.String,
					Artist:     trackArtist.String,
					Title:      trackTitle.String,
					YouTubeURL: trackURL.String,
				}
				// Avoid duplicate tracks
				trackExists := false
				for _, existingTrack := range playlist.Tracks {
					if existingTrack.ID == track.ID {
						trackExists = true
						break
					}
				}
				if !trackExists {
					playlist.Tracks = append(playlist.Tracks, track)
				}
			}
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice
	var entries []MoodEntry
	for _, entry := range entriesMap {
		entries = append(entries, *entry)
	}

	return entries, nil
}
