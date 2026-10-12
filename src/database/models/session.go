package models

import (
	"context"
	"database/sql"
	"encoding/json"
	"sync/atomic"
	"time"
)

// SessionStore keeps folio sessions in the sessions table (migration 000003),
// so logging out ends a session on the server and sessions survive restarts.
// It satisfies folio.SessionStore.
type SessionStore struct {
	Database *sql.DB
	saves    atomic.Uint64
}

func (s *SessionStore) Load(ctx context.Context, id string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var data []byte
	err := s.Database.QueryRowContext(ctx, `
		SELECT data FROM sessions
		WHERE id = $1 AND expires_at > NOW()`, id).Scan(&data)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var values map[string]string
	return values, json.Unmarshal(data, &values)
}

func (s *SessionStore) Save(ctx context.Context, id string, values map[string]string, expires time.Time) error {
	data, err := json.Marshal(values)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	_, err = s.Database.ExecContext(ctx, `
		INSERT INTO sessions (id, data, user_id, expires_at)
		VALUES ($1, $2, NULLIF($3, '')::uuid, $4)
		ON CONFLICT (id) DO UPDATE
		SET data = EXCLUDED.data, user_id = EXCLUDED.user_id, expires_at = EXCLUDED.expires_at`,
		id, data, values["user_id"], expires)
	if err != nil {
		return err
	}

	// Clear out expired sessions every so often.
	if s.saves.Add(1)%500 == 0 {
		_, err = s.Database.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at <= NOW()`)
	}
	return err
}

func (s *SessionStore) Delete(ctx context.Context, id string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := s.Database.ExecContext(ctx, `DELETE FROM sessions WHERE id = $1`, id)
	return err
}

// Ready checks that the sessions table exists, so a missing migration is
// reported at startup instead of every login silently failing.
func (s *SessionStore) Ready(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := s.Database.ExecContext(ctx, `SELECT 1 FROM sessions LIMIT 1`)
	return err
}

// DeleteForUser ends every session a user has, e.g. after a password change.
func (s *SessionStore) DeleteForUser(ctx context.Context, userID string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	_, err := s.Database.ExecContext(ctx, `DELETE FROM sessions WHERE user_id = $1`, userID)
	return err
}
