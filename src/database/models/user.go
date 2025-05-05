package models

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("email already in use")
)

type User struct {
	ID             string    `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"-"`
	CreatedAt      time.Time `json:"created_at"`
}

type UserModel struct {
	Database *sql.DB
}

// Insert adds a new user to the database
func (m *UserModel) Insert(user *User, password string) error {
	// Hash the user's password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Store the user in the database
	query := `
		INSERT INTO users (email, hashed_password)
		VALUES ($1, $2)
		RETURNING id, created_at`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// Set user's fields based on sql query response
	err = m.Database.QueryRowContext(ctx, query, user.Email, string(hashedPassword)).Scan(&user.ID, &user.CreatedAt)
	if err != nil {
		// Check for unique constraint violation on email
		if err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"` {
			return ErrDuplicateEmail
		}
		return err
	}

	return nil
}

// Authenticate validates a user's email and password
func (m *UserModel) Authenticate(email, password string) (*User, error) {
	// Retrieve the user with the provided email
	query := `
		SELECT id, email, hashed_password, created_at
		FROM users
		WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.Database.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.HashedPassword,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	// Check if the provided password matches
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	return &user, nil
}

// GetByID retrieves a user by ID
func (m *UserModel) GetByID(id string) (*User, error) {
	query := `
		SELECT id, email, created_at
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.Database.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetByEmail retrieves a user by email
func (m *UserModel) GetByEmail(email string) (*User, error) {
	query := `
		SELECT id, email, created_at
		FROM users
		WHERE email = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var user User
	err := m.Database.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user's details
func (m *UserModel) Update(user *User) error {
	query := `
		UPDATE users
		SET email = $1
		WHERE id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := m.Database.ExecContext(ctx, query, user.Email, user.ID)
	return err
}

// ChangePassword updates a user's password
func (m *UserModel) ChangePassword(userID, currentPassword, newPassword string) error {
	// Get the current user to check the current password
	query := `
		SELECT hashed_password
		FROM users
		WHERE id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var hashedPassword string
	err := m.Database.QueryRowContext(ctx, query, userID).Scan(&hashedPassword)
	if err != nil {
		return err
	}

	// Check if the current password matches
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(currentPassword))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return ErrInvalidCredentials
		}
		return err
	}

	// Hash the new password
	newHashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// Update the password in the database
	updateQuery := `
		UPDATE users
		SET hashed_password = $1
		WHERE id = $2`

	_, err = m.Database.ExecContext(ctx, updateQuery, string(newHashedPassword), userID)
	return err
}