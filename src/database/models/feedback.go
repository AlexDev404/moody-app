// Filename: database/models/feedback.go
package models

import (
	"context"
	"database/sql"
	"time"
)

// Notice every column in our table has a matching field
// `json:"id"` etc are call struct tags. They help us if we
// want to represent our data as JSON (in the future)
type Feedback struct {
	ID        int64     `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Fullname  string    `json:"fullname"`
	Subject   string    `json:"subject"`
	Message   string    `json:"message"`
	Email     string    `json:"email"`
}

type FeedbackModel struct {
	Database *sql.DB
}

// Insert adds a new feedback record to the database
// Notice that we are sending a pointer. So we will not
// make a copy of the feedback data (more efficient)
func (m *FeedbackModel) Insert(feedback *Feedback) error {
	query := `
		INSERT INTO feedback (fullname, subject, message, email)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	// The Scan() method automatically assigns the returned id and created_at // values to the Feedback struct. This is the reason we needed the address
	// of the feedback data object
	return m.Database.QueryRowContext(
		ctx,
		query,
		feedback.Fullname,
		feedback.Subject,
		feedback.Message,
		feedback.Email,
	).Scan(&feedback.ID, &feedback.CreatedAt)
}

func (m *FeedbackModel) GetAll() ([]*Feedback, error) {
	query := `
		SELECT id, created_at, fullname, subject, message, email
		FROM feedback
		ORDER BY created_at DESC`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := m.Database.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	feedbacks := []*Feedback{}

	for rows.Next() {
		var feedback Feedback
		if err := rows.Scan(&feedback.ID, &feedback.CreatedAt, &feedback.Fullname, &feedback.Subject, &feedback.Message, &feedback.Email); err != nil {
			return nil, err
		}
		feedbacks = append(feedbacks, &feedback)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return feedbacks, nil
}
