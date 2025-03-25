package models

import (
	"context"
	"database/sql"
	"time"
)

type Journal struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
}

type JournalModel struct {
	Database *sql.DB
}

func (m *JournalModel) Insert(journal *Journal) error {
	query := `
		INSERT INTO journal (title, content)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return m.Database.QueryRowContext(
		ctx,
		query,
		journal.Title,
		journal.Content,
	).Scan(&journal.ID, &journal.CreatedAt)
}

func (m *JournalModel) GetAll() ([]*Journal, error) {
	query := `
		SELECT id, created_at, title, content
		FROM journal
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	journals := []*Journal{}

	for rows.Next() {
		var journal Journal
		if err := rows.Scan(&journal.ID, &journal.CreatedAt, &journal.Title, &journal.Content); err != nil {
			return nil, err
		}
		journals = append(journals, &journal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return journals, nil
}
