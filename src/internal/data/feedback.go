// Filename: internal/data/feedback.go
package data

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
	DB *sql.DB
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
	return m.DB.QueryRowContext(
		ctx,
		query,
		feedback.Fullname,
		feedback.Subject,
		feedback.Message,
		feedback.Email,
	).Scan(&feedback.ID, &feedback.CreatedAt)
}
